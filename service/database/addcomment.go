package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"fmt"
)

func (db *appdbimpl) AddComment(comment Comment) error {
	// Execute the INSERT statement to add the new comment to the comments table
	_, err := db.c.Exec("INSERT INTO comments (comment_id , photo_id, username, text, date) VALUES (?, ?, ?, ?, ?)",
		comment.CommentId, comment.PhotoId, comment.Author, comment.Text, comment.Date)
	if err != nil {
		// Return the error if the INSERT operation fails
		return fmt.Errorf("error adding comment to database: %w", err)
	}
	return nil
}
