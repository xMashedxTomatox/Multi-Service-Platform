package services

import (
	"database/sql"
	"errors"
	"log"

	"github.com/xmashedxtomatox/auth-service/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(db *sql.DB, username, email, password string) (*models.User, error) {
	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Insert user and return the newly created row
	var id int
	err = db.QueryRow(
		"INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id",
		username, email, string(hashed),
	).Scan(&id)
	if err != nil {
		// ✅ Log DB error before returning
		log.Printf("CreateUser DB error: %v", err)
		return nil, err
	}

	// Construct user model
	user := &models.User{
		ID:       id,
		Username: username,
		Email:    email,
	}
	return user, nil
}

func AuthenticateUser(db *sql.DB, email, password string) (*models.User, error) {
	var u models.User
	var hashedPassword string

	// Note: use password_hash (consistent with CreateUser)
	err := db.QueryRow(
		"SELECT id, username, email, password_hash FROM users WHERE email=$1",
		email,
	).Scan(&u.ID, &u.Username, &u.Email, &hashedPassword)

	if err != nil {
		return nil, errors.New("user not found")
	}

	// Compare stored hash with provided password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &u, nil
}
