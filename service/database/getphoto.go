package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) GetPhoto(photoID int64) (PhotoMultipart, error) {
	// Prepare the SQL query to retrieve the photo data
	query := "SELECT photo, photo_id, author, upload_datetime, location, caption FROM photos WHERE photo_id = ?"

	// Execute the query
	row := db.c.QueryRow(query, photoID)

	// Initialize variables to store the retrieved photo data
	var photoData []byte
	var photo PhotoMultipart

	// Scan the row to extract photo data
	err := row.Scan(&photoData, &photo.PhotoId, &photo.Author, &photo.UploadDateTime, &photo.Location, &photo.Caption)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return photo, err
		}
		return photo, fmt.Errorf("error retrieving photo data: %w", err)
	}

	// Decode photo data from byte slice
	photo.Photo = photoData

	// populate likes i like count
	var likes []Like
	photo.Likes, err = db.GetLikes(photoID)
	if err != nil {
		return photo, err
	}
	photo.LikesCount = len(likes)

	// populate comments
	var comments []Comment
	photo.Comments, err = db.GetComments(photoID)
	if err != nil {
		return photo, err
	}
	photo.CommentsCount = len(comments)

	return photo, nil
}
