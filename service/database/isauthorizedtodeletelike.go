package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) IsAuthorizedToDeleteLike(like Like, username string) (bool, error) {
	// Query the database to check if the given username is the same as the username in the like
	var storedUsername string
	err := db.c.QueryRow("SELECT username FROM likes WHERE photo_id = ? AND username = ?", like.PhotoId, username).Scan(&storedUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return false if no matching record is found
			return false, nil
		}
		// Return the error if the query fails for other reasons
		return false, fmt.Errorf("error querying database: %w", err)
	}

	// Compare the stored username with the given username
	return storedUsername == username, nil
}
