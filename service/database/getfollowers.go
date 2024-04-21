package database

// "fantastic-coffee-decaffeinated/service/api/models"

func (db *appdbimpl) GetFollowers(username string) ([]Follower, error) {
	// Query the database for followers of the given username
	rows, err := db.c.Query("SELECT follower FROM follow WHERE followed = ?", username)
	if err != nil {
		// Handle the error
		return nil, err
	}
	defer rows.Close()

	// Initialize an empty slice to store followers
	var followers []Follower

	// Iterate over the rows and extract follower usernames
	for rows.Next() {
		var follower string
		if err := rows.Scan(&follower); err != nil {
			// Handle the error
			return nil, err
		}
		// Append the follower username to the slice
		followers = append(followers, Follower{Follower: follower})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// Check if there are no followers found
	if len(followers) == 0 {
		// Return an empty array without error
		return make([]Follower, 0), nil
	}

	return followers, nil
}
