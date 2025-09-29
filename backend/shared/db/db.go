package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	var dsn string

	// First priority: full DATABASE_URL
	if env := os.Getenv("DATABASE_URL"); env != "" {
		dsn = env
	} else {
		// Build DSN from individual env vars
		user := os.Getenv("POSTGRES_USER")
		pass := os.Getenv("POSTGRES_PASSWORD")
		host := os.Getenv("DB_HOST") // 👈 will be "auth-db" inside Docker
		port := os.Getenv("DB_PORT") // usually "5432"
		name := os.Getenv("POSTGRES_DB")

		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			user, pass, host, port, name,
		)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("DB ping error:", err)
	}
	log.Println("✅ Connected to Postgres")
	return db
}
