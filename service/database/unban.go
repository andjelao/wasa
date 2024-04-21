package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"fmt"
)

func (db *appdbimpl) Unban(banningUsername, bannedUsername string) error {
	// Execute the DELETE statement to remove the ban
	_, err := db.c.Exec("DELETE FROM Bans WHERE banningUser = ? AND bannedUser = ?", banningUsername, bannedUsername)
	if err != nil {
		return fmt.Errorf("error deleting ban from database: %w", err)
	}
	return nil
}
