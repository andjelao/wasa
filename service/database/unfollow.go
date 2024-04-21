package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) Unfollow(follower, followed string) error {
	_, err := db.c.Exec("DELETE FROM follow WHERE follower = ? AND followed = ?", follower, followed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Handle case where the follow relationship doesn't exist
			return err
		}
		// Handle other errors
		return fmt.Errorf("failed to unfollow: %w", err)
	}
	return nil
}
