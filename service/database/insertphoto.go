package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"fmt"
)

// InsertPhoto inserts a new photo record into the database.
func (db *appdbimpl) InsertPhoto(photo PhotoMultipart) (int64, error) {
	result, err := db.c.Exec(`
		INSERT INTO photos ( photo, author,upload_datetime,location,caption)
		VALUES ( ?, ?, ?, COALESCE(?, ''), ?)
	`, photo.Photo, photo.Author, photo.UploadDateTime, photo.Location, photo.Caption)

	if err != nil {
		return 0, fmt.Errorf("error inserting photo record: %w", err)
	}
	// Retrieve the last inserted ID
	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error inserting photo record: %w", err)
	}
	return lastInsertedID, nil
}
