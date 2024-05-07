package database

import (
	// "database/sql"

	"database/sql"
	"encoding/base64"
	"errors"
	"strings"
	"time"
)

func (db *appdbimpl) PhotoStream(username string, sinceDateTime string) ([]PhotoMultipart, error) {
	var photos []PhotoMultipart
	var err error
	var sinceTime time.Time
	if sinceDateTime != "" {
		sinceTime, err = time.Parse(time.RFC3339, sinceDateTime)
		if err != nil {
			return nil, err
		}
	}

	// Query to retrieve users followed by the given username
	followQuery := `SELECT followed FROM follow WHERE follower = ?`

	// Execute the query
	rows, err := db.c.Query(followQuery, username)
	if err != nil {
		if strings.Contains(err.Error(), "no such column") {
			// Return an empty user and the error
			return []PhotoMultipart{}, nil
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	var followedUsers []string
	// Iterate over the result rows and populate followedUsers slice
	for rows.Next() {
		var followedUser string
		if err := rows.Scan(&followedUser); err != nil {
			return nil, err
		}
		followedUsers = append(followedUsers, followedUser)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	placeholders := make([]string, len(followedUsers))
	for i := range followedUsers {
		placeholders[i] = "?"
	}

	// Base query to retrieve photos for the followed users
	query := "SELECT * FROM photos WHERE author IN (" + strings.Join(placeholders, ",") + ")"

	// If sinceDateTime is provided, add condition to fetch photos uploaded after sinceDateTime

	if !sinceTime.IsZero() {
		query += ` AND upload_datetime >= ?`

	}

	// Add order by clause to sort photos in reverse chronological order
	query += ` ORDER BY upload_datetime DESC`
	args := make([]interface{}, len(followedUsers))
	for i, user := range followedUsers {
		args[i] = user
	}

	if !sinceTime.IsZero() {
		args = append(args, sinceTime)
	}
	// Execute the query
	rows, err = db.c.Query(query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Set PhotosCount to 0 if no rows are found
			return []PhotoMultipart{}, nil
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	// Iterate over the result rows and populate photos
	for rows.Next() {
		var photo PhotoMultipart
		if err := rows.Scan(&photo.PhotoId, &photo.Photo, &photo.Author, &photo.UploadDateTime, &photo.Location, &photo.Caption); err != nil {
			return nil, err
		}

		if len(photo.Photo) > 0 {
			photo.PhotoEncoded = base64.StdEncoding.EncodeToString(photo.Photo)
			// fmt.Println("encoded" + photo.PhotoEncoded)
			photo.Photo = nil

		}

		// populate likes i like count

		photo.Likes, err = db.GetLikes(photo.PhotoId)
		if err != nil {
			return photos, err
		}
		photo.LikesCount = len(photo.Likes)

		// populate comments

		photo.Comments, err = db.GetComments(photo.PhotoId)
		if err != nil {
			return photos, err
		}
		photo.CommentsCount = len(photo.Comments)
		photos = append(photos, photo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(photos) == 0 {
		return []PhotoMultipart{}, nil
	}

	return photos, nil
}
