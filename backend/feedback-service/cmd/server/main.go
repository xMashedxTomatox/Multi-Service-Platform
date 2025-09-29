package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/xmashedxtomatox/feedback-service/internal/handlers"
	"github.com/xmashedxtomatox/feedback-service/internal/middleware"
	"github.com/xmashedxtomatox/feedback-service/internal/services"

	"github.com/xmashedxtomatox/shared/db"
)

func main() {
	// DB connection
	var database = db.ConnectDB() // ✅ picks up DATABASE_URL or POSTGRES_* env vars
	defer database.Close()

	// JWT secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not set")
	}

	// Services & handlers
	feedbackService := services.NewFeedbackService(database)
	feedbackHandler := handlers.NewFeedbackHandler(feedbackService)

	// Router
	r := mux.NewRouter()
	r.Handle("/feedback", middleware.JWTAuth(jwtSecret, feedbackHandler)).Methods("POST", "GET")

	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)

	log.Printf("Feedback service running on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
