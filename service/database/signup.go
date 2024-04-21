package database

// "fantastic-coffee-decaffeinated/service/api/models"

func (db *appdbimpl) CreateUser(username string) (User, error) {
	// Execute the SQL INSERT statement to add a new user with the given username
	_, err := db.c.Exec("INSERT INTO users (username) VALUES (?)", username)
	if err != nil {
		// Handle the error
		return User{}, err
	}

	// Fetch the newly created user from the database
	var newUser User
	err = db.c.QueryRow("SELECT * FROM users WHERE username = ?", username).Scan(&newUser.Username)
	if err != nil {
		// Handle the error
		return User{}, err
	}

	// Return a copy of the newly created user object
	return newUser, nil
}
