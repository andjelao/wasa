package database

// "fantastic-coffee-decaffeinated/service/api/models"

func (db *appdbimpl) IsAuthenticatedUser(username string) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		// Handle error
		return false, err
	}
	return count > 0, nil
}
