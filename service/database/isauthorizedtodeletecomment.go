package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) IsAuthorizedToDeleteComment(username string, commentID int) (bool, error) {
	// Query the database to get the username associated with the comment
	var commentAuthor string
	err := db.c.QueryRow("SELECT username FROM comments WHERE comment_id = ?", commentID).Scan(&commentAuthor)
	if err != nil {
		// Check if the error is due to no rows being found
		if errors.Is(err, sql.ErrNoRows) {
			// Return false and an error indicating that the comment doesn't exist
			return false, fmt.Errorf("comment with ID %d does not exist", commentID)
		}
		// Return false and the error for any other issues
		return false, fmt.Errorf("error retrieving comment author: %w", err)
	}

	// Compare the retrieved username with the given username
	return username == commentAuthor, nil
}
