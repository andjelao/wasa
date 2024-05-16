package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	// "log"
	// "io/ioutil"
	// "mime"
	"database/sql"
	"errors"

	// "math/rand"
	"strconv"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	// "git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	// "fantastic-coffee-decaffeinated/service/api/models"
)

// ///////////////////////
// type Comment struct {
// commentId  int `json:"commentId"`
// photoId int `json:"commentId"`
// author  string `json:"author"`
// text string `json:"text"`
// date string `json:"date"`
// }
// /////////////////////////
// comment a photo
func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := GetUsernameFromToken(r)
	// Check if the author exists in the database of authenticated users and is logged in
	authenticated, err := rt.db.IsAuthenticatedUser(username)
	if err != nil {
		// Handle error if check fails
		http.Error(w, "Authentication check fail", http.StatusInternalServerError)
		return
	}
	if username == "" || !authenticated {
		// Return error response if the author is not authenticated
		http.Error(w, "Unauthorized - user is not authenticated", http.StatusUnauthorized)
		return
	}
	// Extract the photoId from the path parameters
	photoIdd := ps.ByName("photoId")
	if photoIdd == "" {
		http.Error(w, "Photo ID is required", http.StatusBadRequest)
		return
	}

	photoId, err := strconv.ParseInt(photoIdd, 10, 64)
	if err != nil {
		http.Error(w, "Internal server error-error converting photoid to string", http.StatusInternalServerError)
		return
	}

	// moram da znam autora fotke da bi znala je li banned
	author, err := rt.db.ExtractAuthor(photoId)
	if err != nil {
		// Handle error
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Photo not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	banned, err := rt.db.BanExists(author, username)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if banned {
		http.Error(w, "Forbidden - the requester is banned by the author of the photo", http.StatusForbidden)
		return
	}

	// Decode the JSON request body
	var newComment database.Comment

	////// mozda je unused
	var commentText database.CommentText

	// var comm string
	if err := json.NewDecoder(r.Body).Decode(&commentText); err != nil {
		http.Error(w, "Failed to decode JSON body", http.StatusBadRequest)
		return
	}
	newComment.Text = commentText.Text
	newComment.PhotoId = photoId
	newComment.Author = username
	// Calculate upload date and time
	uploadDateTime := time.Now().Format(time.RFC3339)
	newComment.Date = uploadDateTime
	// commentId := rand.Int()
	// newComment.CommentId = commentId

	//
	// Save the updated photo back to the database
	var idd int64
	if idd, err = rt.db.AddComment(newComment); err != nil {
		http.Error(w, "Failed to add comment in the database", http.StatusInternalServerError)
		return
	}
	newComment.CommentId = idd

	// respond with a comment
	jsonResponse, err := json.Marshal(newComment)
	if err != nil {
		// Handle error
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(jsonResponse)
	if err != nil {
		// Handle error if writing response fails.
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
