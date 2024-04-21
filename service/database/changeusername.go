package database

// "database/sql"

func (db *appdbimpl) ChangeUsername(username, newUsername string) error {

	// ChangeUsername updates the username in the database
	// Prepare the SQL statement
	query := "UPDATE users SET username = ? WHERE username = ?"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(newUsername, username)
	if err != nil {
		return err
	}

	return nil

}
