package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"fmt"
)

// GetComments retrieves all comments associated with the given photo ID.
func (db *appdbimpl) GetComments(photoID int64) ([]Comment, error) {
	// Prepare SQL statement to select comments for the given photo ID
	query := "SELECT comment_id, photo_id, username, text, date FROM comments WHERE photo_id = ?"
	rows, err := db.c.Query(query, photoID)
	if err != nil {
		return nil, fmt.Errorf("error querying comments: %w", err)
	}
	defer rows.Close()

	// Iterate over rows and scan comments into a slice
	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.CommentId, &comment.PhotoId, &comment.Author, &comment.Text, &comment.Date); err != nil {
			return nil, fmt.Errorf("error scanning comment row: %w", err)
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over comments rows: %w", err)
	}

	if len(comments) == 0 {
		return make([]Comment, 0), nil
	}

	return comments, nil
}
