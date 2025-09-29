package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/xmashedxtomatox/feedback-service/internal/handlers"
	"github.com/xmashedxtomatox/feedback-service/internal/middleware"
	"github.com/xmashedxtomatox/feedback-service/internal/services"

	"github.com/xmashedxtomatox/shared/db"
)

func setupTestRouter(database *sql.DB, jwtSecret string) *mux.Router {
	service := services.NewFeedbackService(database)
	handler := handlers.NewFeedbackHandler(service)

	r := mux.NewRouter()
	r.Handle("/feedback", middleware.JWTAuth(jwtSecret, handler)).Methods("POST", "GET")
	return r
}

func generateJWT(userID int, secret string) string {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": "user",
		"exp":  time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}

func TestFeedbackFlow(t *testing.T) {
	// Use shared DB connection
	database := db.ConnectDB()
	defer database.Close()

	jwtSecret := os.Getenv("JWT_SECRET")
	assert.NotEmpty(t, jwtSecret, "JWT_SECRET must be set for integration tests")

	router := setupTestRouter(database, jwtSecret)

	// --- 1. Create feedback ---
	body := []byte(`{"message":"Integration test feedback"}`)
	req := httptest.NewRequest("POST", "/feedback", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+generateJWT(123, jwtSecret))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var created map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &created)
	assert.NoError(t, err)
	assert.Equal(t, "Integration test feedback", created["message"])

	// --- 2. Get feedbacks ---
	req2 := httptest.NewRequest("GET", "/feedback", nil)
	req2.Header.Set("Authorization", "Bearer "+generateJWT(123, jwtSecret))
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	var feedbacks []map[string]interface{}
	err = json.Unmarshal(w2.Body.Bytes(), &feedbacks)
	assert.NoError(t, err)
	assert.NotEmpty(t, feedbacks)
	assert.Equal(t, "Integration test feedback", feedbacks[0]["message"])
}
