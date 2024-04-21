package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"fmt"
)

func (db *appdbimpl) DeleteComment(commentID int, photoID int) error {

	// Delete the comment from the database where both comment ID and photo ID match
	_, err := db.c.Exec("DELETE FROM comments WHERE comment_id = ? AND photo_id = ?", commentID, photoID)
	if err != nil {
		return fmt.Errorf("error deleting comment from database: %w", err)
	}

	return nil // Return nil if deletion is successful
}
