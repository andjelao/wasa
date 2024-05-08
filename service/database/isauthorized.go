package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) IsAuthorized(username string, photoID int64) (bool, error) {
	// Query the database to check if the provided username is the author of the photo identified by photoID

	var authorr string
	err := db.c.QueryRow("SELECT author FROM photos WHERE photo_id = ?", photoID).Scan(&authorr)
	// err := db.c.QueryRow("SELECT author FROM photos WHERE photo_id = 7371312785867452000").Scan(&authorr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No photo found with the given photoID
			return false, err
		}
		// Some other error occurred
		return false, fmt.Errorf("error querying database: %w", err)
	}

	// Check if the provided username matches the author of the photo
	return username == authorr, nil
}
