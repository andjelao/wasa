package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"database/sql"
	"errors"
)

func (db *appdbimpl) Login(username string) (User, error) {
	// Query the database for the user record with the given username
	row := db.c.QueryRow("SELECT username FROM users WHERE username = ?", username)

	// Initialize a variable to store the user record
	var user User

	// Scan the row into the user struct
	if err := row.Scan(&user.Username); err != nil {
		// Check if the error is due to no rows found
		if errors.Is(err, sql.ErrNoRows) {
			// Return an empty user and the error
			return User{}, err
		}
		// Handle other errors
		return User{}, err
	}

	// Return the user record
	return user, nil
}
