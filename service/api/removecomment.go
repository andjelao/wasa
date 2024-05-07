package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	// "encoding/json"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	// "fantastic-coffee-decaffeinated/service/database"

	// "log"
	// "io/ioutil"
	// "mime"
	// "time"
	"database/sql"
	"errors"
	"strconv"
	// "fantastic-coffee-decaffeinated/service/api/models"
)

// user removes their own comment
func (rt *_router) uncomment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Extract the photoId from the path parameters
	photoIdd := ps.ByName("photoId")
	if photoIdd == "" {
		http.Error(w, "Photo ID is required", http.StatusBadRequest)
		return
	}
	photoId, err := strconv.ParseInt(photoIdd, 10, 64)
	if err != nil {
		http.Error(w, "Internal server error- error converting photoid to string", http.StatusInternalServerError)
		return
	}

	commentIdd := ps.ByName("commentId")
	if commentIdd == "" {
		http.Error(w, "Comment ID is required", http.StatusBadRequest)
		return
	}
	commentId, err := strconv.Atoi(ps.ByName("commentId"))
	if err != nil {
		http.Error(w, "Internal server error- error converting photoid to string", http.StatusInternalServerError)
		return
	}
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
		http.Error(w, "Unauthorized user", http.StatusUnauthorized)
		return
	}

	// check if the comment already exists for this photo
	commentExists, err := rt.db.CommentExists(commentId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Photo not found", http.StatusNotFound)
			return
		}
		// Handle error
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !commentExists {
		http.Error(w, "Comment not found", http.StatusNotFound)
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

	// is requestor the author of the comment
	authorized, err := rt.db.IsAuthorizedToDeleteComment(username, commentId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Comment not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	if !authorized {
		// If user is not authorized to delete the comment, send forbidden response
		http.Error(w, "Forbidden - the requester is not the author of the comment", http.StatusForbidden)
		return
	}

	// delete comment from the database
	err = rt.db.DeleteComment(commentId, photoId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusNoContent)
}
