/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	// "fantastic-coffee-decaffeinated/service/api/models"
)

type Ban struct {
	BannedUsername string `json:"bannedUsername,omitempty"`
}

type Comment struct {
	CommentId int64  `json:"commentId"`
	PhotoId   int64  `json:"photoId"`
	Author    string `json:"author"`
	Text      string `json:"text"`
	Date      string `json:"date"`
}
type Followed struct {
	Following string `json:"followedUsername"`
}
type Follower struct {
	Follower string `json:"followerUsername"`
}
type Like struct {
	Username string `json:"username"`
	PhotoId  int64  `json:"photoId"`
}
type PhotoMultipart struct {
	// koji format da stavim da li jpg ili ovako bytes
	Photo []byte `json:"photo"`
	//
	PhotoId        int64     `json:"photoId"`
	Author         string    `json:"author"`
	UploadDateTime string    `json:"uploadDateTime"`
	Location       string    `json:"location,omitempty"`
	Caption        string    `json:"caption,omitempty"`
	LikesCount     int       `json:"likesCount"`
	Likes          []Like    `json:"likes"`
	CommentsCount  int       `json:"commentsCount"`
	Comments       []Comment `json:"comments"`
	PhotoEncoded   string    `json:"photoEncoded"`
}
type Profile struct {
	Username       string           `json:"username"`
	PhotosCount    int              `json:"photosCount"`
	FollowersCount int              `json:"followersCount"`
	FollowersList  []Follower       `json:"followersList"`
	FollowingCount int              `json:"followingCount"`
	FollowingList  []Followed       `json:"followingList"`
	UserPhotos     []PhotoMultipart `json:"userPhotos"`
}
type UpdateRequest struct {
	Caption  string `json:"caption,omitempty"`
	Location string `json:"location,omitempty"`
}
type User struct {
	Username    string           `json:"username"`
	Profile     Profile          `json:"profile"`
	Banned      []Ban            `json:"banned"`
	Photostream []PhotoMultipart `json:"photostream"`
}
type ConflictResponseBan struct {
	Message     string `json:"message"`
	ExistingBan Ban    `json:"existingBan"`
}

type ConflictResponseFollow struct {
	Message  string   `json:"message"`
	Followed Followed `json:"followed"`
}
type ConflictResponseLike struct {
	Message string `json:"message"`
	Like    Like   `json:"like"`
}

type CommentText struct {
	Text string `json:"text"`
}

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetName() (string, error)
	SetName(name string) error

	InsertPhoto(photo PhotoMultipart) (int64, error)
	IsAuthenticatedUser(username string) (bool, error)
	IsAuthorized(username string, photoID int64) (bool, error)
	UpdatePhoto(photoID int64, req UpdateRequest) error
	GetLikes(photoId int64) ([]Like, error)
	GetComments(photoID int64) ([]Comment, error)
	GetPhoto(photoID int64) (PhotoMultipart, error)
	DeletePhoto(photoId int64) error
	BanExists(banningUsername, bannedUsername string) (bool, error)
	AddBan(banningUsername, bannedUsername string) error
	GetBanned(username string) ([]Ban, error)
	Unban(banningUsername, bannedUsername string) error
	ExtractAuthor(photoID int64) (string, error)
	AddComment(comment Comment) (int64, error)
	LikeExists(like Like) (bool, error)
	AddLike(like Like) error
	IsAuthorizedToDeleteLike(like Like, username string) (bool, error)
	DeleteLike(like Like) error
	CommentExists(commentID int) (bool, error)
	IsAuthorizedToDeleteComment(username string, commentID int) (bool, error)
	DeleteComment(commentID int, photoID int64) error
	FollowExists(follower, following string) (bool, error)
	AddFollow(follower, following string) error
	Unfollow(follower, followed string) error
	GetFollowed(username string) ([]Followed, error)
	GetFollowers(username string) ([]Follower, error)
	Login(username string) (User, error)
	CreateUser(username string) (User, error)
	RetrievePhotos(users []string, myUsername string) ([]PhotoMultipart, error)
	PhotoStream(username string, sinceDateTime string) ([]PhotoMultipart, error)
	GetProfile(username string) (Profile, error)
	ChangeUsername(username, newUsername string) error
	GetTables() error

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Enable foreign keys
	_, err := db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, err
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `

			CREATE TABLE IF NOT EXISTS photos (
				photo_id INTEGER PRIMARY KEY AUTOINCREMENT,
				photo BLOB NOT NULL,
				author TEXT NOT NULL,
				upload_datetime TEXT NOT NULL,
				location TEXT,
				caption TEXT,
				FOREIGN KEY (author) REFERENCES users(username) ON DELETE CASCADE ON UPDATE CASCADE
			);
			CREATE TABLE IF NOT EXISTS likes (
				like_id INTEGER PRIMARY KEY,
				photo_id INTEGER,
				username TEXT,
				FOREIGN KEY (photo_id) REFERENCES photos(photo_id) ON DELETE CASCADE,
				FOREIGN KEY (username) REFERENCES users(username) ON DELETE CASCADE ON UPDATE CASCADE
			);
			
			CREATE TABLE IF NOT EXISTS comments (
				comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
				photo_id INTEGER,
				username TEXT,
				text TEXT,
				date TEXT,
				FOREIGN KEY (photo_id) REFERENCES photos(photo_id) ON DELETE CASCADE,
				FOREIGN KEY (username) REFERENCES users(username) ON DELETE CASCADE ON UPDATE CASCADE
			);
			CREATE TABLE IF NOT EXISTS users (
				username TEXT PRIMARY KEY
			);
			CREATE TABLE IF NOT EXISTS Bans (
				banId INTEGER PRIMARY KEY,
				banningUser TEXT NOT NULL,
				bannedUser TEXT NOT NULL,
				FOREIGN KEY (banningUser) REFERENCES users(username) ON DELETE CASCADE ON UPDATE CASCADE,
				FOREIGN KEY (bannedUser) REFERENCES users(username) ON DELETE CASCADE ON UPDATE CASCADE
			);

			
			CREATE TABLE IF NOT EXISTS follow (
				followId INTEGER PRIMARY KEY,
				follower TEXT,
				followed TEXT,
				FOREIGN KEY (follower) REFERENCES users(username) ON DELETE CASCADE ON UPDATE CASCADE,
				FOREIGN KEY (followed) REFERENCES users(username) ON DELETE CASCADE ON UPDATE CASCADE
			);

			`

		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

// getphoto by username, get all photos, by id
// delete photo
// save photo
// DeleteLikesForPhoto(photoID)
// rt.DB.DeleteCommentsForPhoto
// updatedata slike
// add like
// get likes
// delete like
// add comment
// get comments
// delete comment
// add follow
// add follower
// get followed
// unfollow
// get followers
// ban user
// get banned
// unban
// photostream
// profile
// getuser
// create user
// change username
