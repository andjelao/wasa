package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"database/sql"
	"errors"
	"fmt"
)

// deletePhoto deletes a photo with the given photoId from the database.
func (db *appdbimpl) DeletePhoto(photoId int) error {
	// Prepare the SQL statement to delete the photo from the database.
	sqlStmt := "DELETE FROM photos WHERE photo_id= ?"

	// Execute the SQL statement with the provided photoId parameter.
	_, err := db.c.Exec(sqlStmt, photoId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return a specific error indicating that the photo was not found.
			return err
		}
		// Return the error if the SQL statement execution fails.
		return fmt.Errorf("failed to delete photo: %w", err)
	}

	// Return nil to indicate success.
	return nil
}
