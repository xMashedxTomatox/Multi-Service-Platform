package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/xmashedxtomatox/auth-service/internal/handlers"
	"github.com/xmashedxtomatox/auth-service/internal/utils"
	"github.com/xmashedxtomatox/shared/db"
	"github.com/xmashedxtomatox/shared/redis"
)

func TestSignUpIntegration(t *testing.T) {
	// Arrange: connect to Dockerized DB
	db := db.ConnectDB()
	defer db.Close()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		t.Fatal("JWT_SECRET environment variable is not set")
	}

	redisClient := redis.NewRedisClient()
	defer redisClient.Close()

	router := mux.NewRouter()
	router.HandleFunc("/signup", handlers.SignUp(db, redisClient)).Methods("POST")

	// Generate unique username/email per test run
	ts := time.Now().UnixNano()
	username := fmt.Sprintf("integration_user_%d", ts)
	email := fmt.Sprintf("integration_%d@example.com", ts)
	t.Logf("🔑 Using test credentials: username=%s email=%s", username, email)

	payload := fmt.Sprintf(`{"username":"%s","email":"%s","password":"password123"}`, username, email)

	req := httptest.NewRequest("POST", "/signup", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}

	// Extra: validate JWT works with the secret
	token, err := utils.GenerateJWT(1, jwtSecret)
	if err != nil {
		t.Fatalf("failed to generate token in integration test: %v", err)
	}
	if token == "" {
		t.Error("expected a non-empty JWT")
	}
}

func TestLoginIntegration(t *testing.T) {
	db := db.ConnectDB()
	defer db.Close()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		t.Fatal("JWT_SECRET environment variable is not set")
	}

	router := mux.NewRouter()
	redisClient := redis.NewRedisClient()
	defer redisClient.Close()

	router.HandleFunc("/signup", handlers.SignUp(db, redisClient)).Methods("POST")
	router.HandleFunc("/login", handlers.Login(db, redisClient)).Methods("POST")

	username := "login_test_user"
	email := "login_test@example.com"
	password := "mypassword"

	// --- SignUp first ---
	signupPayload := fmt.Sprintf(`{"username":"%s","email":"%s","password":"%s"}`, username, email, password)
	reqSignup := httptest.NewRequest("POST", "/signup", strings.NewReader(signupPayload))
	reqSignup.Header.Set("Content-Type", "application/json")
	wSignup := httptest.NewRecorder()
	router.ServeHTTP(wSignup, reqSignup)

	if wSignup.Code != http.StatusOK {
		t.Fatalf("signup failed: got %d: %s", wSignup.Code, wSignup.Body.String())
	}

	// --- Login with same credentials ---
	loginPayload := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)
	reqLogin := httptest.NewRequest("POST", "/login", strings.NewReader(loginPayload))
	reqLogin.Header.Set("Content-Type", "application/json")
	wLogin := httptest.NewRecorder()
	router.ServeHTTP(wLogin, reqLogin)

	if wLogin.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for login, got %d: %s", wLogin.Code, wLogin.Body.String())
	}

	t.Logf("✅ Login response: %s", wLogin.Body.String())
}
