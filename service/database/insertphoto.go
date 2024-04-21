package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"fmt"
)

// InsertPhoto inserts a new photo record into the database.
func (db *appdbimpl) InsertPhoto(photo PhotoMultipart) error {
	_, err := db.c.Exec(`
		INSERT INTO photos (photo_id, photo, author,upload_datetime,location,caption)
		VALUES (?, ?, ?, ?, ?)
	`, photo.PhotoId, photo.Photo, photo.Author, photo.UploadDateTime, photo.Location, photo.Caption)
	if err != nil {
		return fmt.Errorf("error inserting photo record: %w", err)
	}
	return nil
}
