package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"fmt"
)

func (db *appdbimpl) GetFollowed(username string) ([]Followed, error) {
	rows, err := db.c.Query("SELECT followed FROM follow WHERE follower = ?", username)
	if err != nil {
		return nil, fmt.Errorf("failed to get followed users: %w", err)
	}
	defer rows.Close()

	var followedUsers []Followed
	for rows.Next() {
		var followedUser string
		if err := rows.Scan(&followedUser); err != nil {
			return nil, fmt.Errorf("failed to scan followed user: %w", err)
		}
		followedUsers = append(followedUsers, Followed{Following: followedUser})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over followed users: %w", err)
	}
	if len(followedUsers) == 0 {
		return make([]Followed, 0), nil
	}
	return followedUsers, nil
}
