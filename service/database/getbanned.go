package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"fmt"
)

func (db *appdbimpl) GetBanned(username string) ([]Ban, error) {
	var bans []Ban

	// Query the Bans table for rows where the given username is the banningUsername
	rows, err := db.c.Query("SELECT bannedUser FROM Bans WHERE banningUser = ?", username)
	if err != nil {
		return nil, fmt.Errorf("error querying bans table: %w", err)
	}
	defer rows.Close()

	// Iterate over the rows and populate the bans slice
	for rows.Next() {
		var bannedUsername string
		if err := rows.Scan(&bannedUsername); err != nil {
			return nil, fmt.Errorf("error scanning ban row: %w", err)
		}
		bans = append(bans, Ban{BannedUsername: bannedUsername})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over ban rows: %w", err)
	}
	if len(bans) == 0 {
		return make([]Ban, 0), nil
	}

	return bans, nil
}
