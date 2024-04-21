package database

import (
	// "fantastic-coffee-decaffeinated/service/api/models"
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetProfile(username string) (Profile, error) {

	var profile Profile
	profile.Username = username
	// Get username, photos count, and user photos
	userQuery := `SELECT COUNT(*) AS photosCount FROM photos WHERE author = ? GROUP BY author`
	row := db.c.QueryRow(userQuery, username)
	if err := row.Scan(&profile.PhotosCount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Set PhotosCount to 0 if no rows are found
			profile.PhotosCount = 0
		} else {
			return profile, err
		}
	}
	// Get followers count
	followersQuery := `SELECT COUNT(*) AS followersCount FROM follow WHERE followed = ?`
	row = db.c.QueryRow(followersQuery, username)
	if err := row.Scan(&profile.FollowersCount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Set FollowersCount to 0 if no rows are found
			profile.FollowersCount = 0
		} else {
			return profile, err
		}
	}
	// Get followers list
	followersListQuery := `SELECT follower FROM follow WHERE followed = ?`
	rows, err := db.c.Query(followersListQuery, username)
	if err != nil {
		return profile, err
	}
	defer rows.Close()

	for rows.Next() {
		var follower Follower
		if err := rows.Scan(&follower.Follower); err != nil {
			return profile, err
		}
		profile.FollowersList = append(profile.FollowersList, follower)
	}
	if err := rows.Err(); err != nil {
		return profile, err
	}

	if len(profile.FollowersList) == 0 {
		profile.FollowersList = []Follower{}
	}

	// Get following count
	followingQuery := `SELECT COUNT(*) AS followingCount FROM follow WHERE follower = ?`
	row = db.c.QueryRow(followingQuery, username)
	if err := row.Scan(&profile.FollowingCount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Set FollowersCount to 0 if no rows are found
			profile.FollowingCount = 0
		} else {
			return profile, err
		}
	}

	// Get following list
	followingListQuery := `SELECT followed FROM follow WHERE follower = ?`
	rows, err = db.c.Query(followingListQuery, username)
	if err != nil {
		return profile, err
	}
	defer rows.Close()

	for rows.Next() {
		var followed Followed
		if err := rows.Scan(&followed.Following); err != nil {
			return profile, err
		}
		profile.FollowingList = append(profile.FollowingList, followed)
	}
	if err := rows.Err(); err != nil {
		return profile, err
	}

	if len(profile.FollowingList) == 0 {
		profile.FollowingList = []Followed{}
	}

	query := "SELECT * FROM photos WHERE author = ? ORDER BY upload_datetime DESC"
	rows, err = db.c.Query(query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Set FollowersCount to 0 if no rows are found
			profile.UserPhotos = []PhotoMultipart{}
			return profile, nil
		} else {
			return profile, err
		}
	}
	defer rows.Close()

	var photos []PhotoMultipart
	for rows.Next() {
		var photo PhotoMultipart
		if err := rows.Scan(&photo.PhotoId, &photo.Photo, &photo.Author, &photo.UploadDateTime, &photo.Location, &photo.Caption); err != nil {
			return profile, err
		}
		// populate likes i like count
		var likes []Like
		photo.Likes, err = db.GetLikes(photo.PhotoId)
		if err != nil {
			return profile, err
		}
		photo.LikesCount = len(likes)

		// populate comments
		var comments []Comment
		photo.Comments, err = db.GetComments(photo.PhotoId)
		if err != nil {
			return profile, err
		}
		photo.CommentsCount = len(comments)
		photos = append(photos, photo)
	}

	if err := rows.Err(); err != nil {
		return profile, err
	}

	if len(photos) == 0 {
		profile.UserPhotos = []PhotoMultipart{}
		return profile, nil
	}
	profile.UserPhotos = photos

	return profile, nil
}
