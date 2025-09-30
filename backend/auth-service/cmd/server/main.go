package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/xmashedxtomatox/auth-service/internal/handlers"
	"github.com/xmashedxtomatox/shared/db"
	"github.com/xmashedxtomatox/shared/middleware"
	"github.com/xmashedxtomatox/shared/redis"
)

func main() {
	// Connect DB
	db := db.ConnectDB()
	defer db.Close()

	// Connect Redis
	redisClient := redis.NewRedisClient()
	defer redisClient.Close()

	// Router
	r := mux.NewRouter()
	r.HandleFunc("/signup", handlers.SignUp(db, redisClient)).Methods("POST")
	r.HandleFunc("/login", handlers.Login(db, redisClient)).Methods("POST")

	// Debug
	r.HandleFunc("/debug/users", handlers.ListUsers(db)).Methods("GET")

	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)

	log.Printf("Auth service running on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, middleware.WithCORS(r)))
}
