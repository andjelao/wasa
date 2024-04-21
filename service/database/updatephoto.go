package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"database/sql"
	"errors"
	"fmt"
)

// UpdatePhoto updates the photo's caption and location based on the provided photoID.
func (db *appdbimpl) UpdatePhoto(photoID int, req UpdateRequest) error {
	// Prepare SQL statement to update the photo
	sqlStmt := `UPDATE photos SET caption = ?, location = ? WHERE photo_id = ?`

	// Execute the SQL statement
	_, err := db.c.Exec(sqlStmt, req.Caption, req.Location, photoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No photo found with the given photoID
			return err
		}
		return fmt.Errorf("error updating photo: %w", err)
	}

	return nil
}
