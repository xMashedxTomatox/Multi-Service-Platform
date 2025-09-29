package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func ListUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, username, email FROM users")
		if err != nil {
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []map[string]interface{}
		for rows.Next() {
			var id int
			var username, email string
			if err := rows.Scan(&id, &username, &email); err != nil {
				http.Error(w, "Scan error", http.StatusInternalServerError)
				return
			}
			users = append(users, map[string]interface{}{
				"id":       id,
				"username": username,
				"email":    email,
			})
		}
		json.NewEncoder(w).Encode(users)
	}
}
