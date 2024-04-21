package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"fmt"
)

// GetLikes returns all likes associated with the given photo ID.
func (db *appdbimpl) GetLikes(photoId int) ([]Like, error) {
	// SQL query to select likes for the given photoId
	query := "SELECT username FROM likes WHERE photo_id = ?"

	// Execute the query
	rows, err := db.c.Query(query, photoId)
	if err != nil {
		return nil, fmt.Errorf("error querying likes: %w", err)
	}
	defer rows.Close()

	// Initialize a slice to hold the likes
	likes := make([]Like, 0)

	// Iterate over the result set
	for rows.Next() {
		var username string
		// Scan the username from the row
		if err := rows.Scan(&username); err != nil {
			return nil, fmt.Errorf("error scanning like: %w", err)
		}
		// Create a Like object and append it to the slice
		like := Like{
			Username: username,
			PhotoId:  photoId,
		}
		likes = append(likes, like)
	}
	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over likes: %w", err)
	}
	if len(likes) == 0 {
		return make([]Like, 0), nil
	}

	return likes, nil
}
