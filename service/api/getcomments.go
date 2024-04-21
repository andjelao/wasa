package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	// "log"
	// "io/ioutil"
	// "mime"
	// "time"
	"database/sql"
	"errors"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	// "fantastic-coffee-decaffeinated/service/api/models"
)

func (rt *_router) getComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
	photoId, err := strconv.Atoi(ps.ByName("photoId"))
	if err != nil {
		http.Error(w, "Internal server error- error converting photoid to string", http.StatusInternalServerError)
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

	var comments []database.Comment
	// get comments from database
	comments, err = rt.db.GetComments(photoId)
	if len(comments) == 0 {
		http.Error(w, "Photo has no comments", http.StatusNotFound)
		return
	}
	if err != nil {
		// Handle error
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Serialize the list of likes into JSON format.
	jsonResponse, err := json.Marshal(comments)
	if err != nil {
		// Handle error if serialization fails.
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set content type header.
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response back to the client.
	_, err = w.Write(jsonResponse)
	if err != nil {
		// Handle error if writing response fails.
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
