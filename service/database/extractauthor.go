package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) ExtractAuthor(photoID int64) (string, error) {
	// Execute the SELECT statement to retrieve the author of the photo
	var author string
	err := db.c.QueryRow("SELECT author FROM photos WHERE photo_id = ?", photoID).Scan(&author)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return an error indicating that the photo was not found
			return "", err
		}
		// Return the error if it's not related to a missing photo
		return "", fmt.Errorf("error retrieving author from database: %w", err)
	}
	return author, nil
}
