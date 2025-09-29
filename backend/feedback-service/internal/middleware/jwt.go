package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	ContextUserID contextKey = "user_id"
	ContextRole   contextKey = "role"
)

func JWTAuth(secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			http.Error(w, "invalid Authorization format", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid claims", http.StatusUnauthorized)
			return
		}

		// Expiry
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				http.Error(w, "token expired", http.StatusUnauthorized)
				return
			}
		} else {
			http.Error(w, "missing exp claim", http.StatusUnauthorized)
			return
		}

		// Extract user_id and role
		userID, ok := claims["sub"].(float64)
		if !ok {
			http.Error(w, "missing sub claim", http.StatusUnauthorized)
			return
		}
		role, _ := claims["role"].(string)

		// Inject into context
		ctx := context.WithValue(r.Context(), ContextUserID, int(userID))
		if role != "" {
			ctx = context.WithValue(ctx, ContextRole, role)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
