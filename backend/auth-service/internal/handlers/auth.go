package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/xmashedxtomatox/auth-service/internal/services"
	"github.com/xmashedxtomatox/auth-service/internal/utils"
	sharedredis "github.com/xmashedxtomatox/shared/redis"
)

type SignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignUp(db *sql.DB, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SignUpRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Create the user in DB
		user, err := services.CreateUser(db, req.Username, req.Email, req.Password)
		if err != nil {
			http.Error(w, "Could not create user", http.StatusInternalServerError)
			return
		}

		// Fetch JWT secret from env
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			http.Error(w, "JWT secret not configured", http.StatusInternalServerError)
			return
		}

		// Cache user in Redis
		redisClient.HSet(sharedredis.Ctx, fmt.Sprintf("user:%d", user.ID), map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		})

		// Return token
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
		})
	}
}

func Login(db *sql.DB, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		user, err := services.AuthenticateUser(db, req.Email, req.Password)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Fetch secret from environment (set in docker-compose.yml or .env)
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			http.Error(w, "JWT secret not configured", http.StatusInternalServerError)
			return
		}

		// Generate JWT
		token, err := utils.GenerateJWT(user.ID, jwtSecret)
		if err != nil {
			http.Error(w, "Could not generate token", http.StatusInternalServerError)
			return
		}

		// Cache user in Redis
		redisClient.HSet(sharedredis.Ctx, fmt.Sprintf("user:%d", user.ID), map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		})

		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})
	}
}
