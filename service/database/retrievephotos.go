package database

// "database/sql"

func (db *appdbimpl) RetrievePhotos(users []string, myUsername string) ([]PhotoMultipart, error) {

	// RetrievePhotos retrieves photos from the database based on the given criteria.
	// It returns an array of PhotoMultipart objects.
	var photos []PhotoMultipart

	// Query to retrieve photos from the database
	var query string
	var args []interface{}

	if len(users) == 0 {
		// If users array is empty, retrieve all photos
		query = `SELECT * FROM photos WHERE author NOT IN (SELECT banningUser FROM Bans WHERE bannedUser = ?)`
		args = []interface{}{myUsername}
	} else {
		// If users array is not empty, retrieve photos from specific authors
		query = `SELECT * FROM photos WHERE author IN (?) AND author NOT IN (SELECT banningUser FROM Bans WHERE bannedUser = ?)`
		args = append(args, users, myUsername)
	}

	// Prepare the query
	rows, err := db.c.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the result rows and populate photos
	for rows.Next() {
		var photo PhotoMultipart
		if err := rows.Scan(&photo.PhotoId, &photo.Photo, &photo.Author, &photo.UploadDateTime, &photo.Location, &photo.Caption); err != nil {
			return nil, err
		}
		// populate likes i like count
		// makla unused
		// var likes []Like
		photo.Likes, err = db.GetLikes(photo.PhotoId)
		if err != nil {
			return photos, err
		}
		photo.LikesCount = len(photo.Likes)

		// populate comments
		// var comments []Comment
		photo.Comments, err = db.GetComments(photo.PhotoId)
		if err != nil {
			return photos, err
		}
		photo.CommentsCount = len(photo.Comments)
		photos = append(photos, photo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return photos, nil
}
