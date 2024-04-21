package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"fmt"
)

// AddBan adds a new ban to the bans table.
func (db *appdbimpl) AddBan(banningUsername, bannedUsername string) error {
	// Prepare the SQL statement
	query := "INSERT INTO Bans (banningUser, bannedUser) VALUES (?, ?)"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return fmt.Errorf("error preparing SQL statement: %w", err)
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(banningUsername, bannedUsername)
	if err != nil {
		return fmt.Errorf("error executing SQL statement: %w", err)
	}

	return nil
}
