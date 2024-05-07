package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	// "encoding/json"

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

func (rt *_router) dislike(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
	// Retrieve the value of two username path parameter from the request and authentication token
	pathusername := ps.ByName("username")
	if pathusername == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	username := GetUsernameFromToken(r)

	if pathusername != username {
		http.Error(w, "Forbidden - username in bearer token does not match the provided path parameter username", http.StatusForbidden)
		return
	}
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
	// create Like object with obtained properties
	var like database.Like
	like.Username = username
	like.PhotoId = photoId

	// check if the like already exists for this photo
	//
	likeExists, err := rt.db.LikeExists(like)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !likeExists {
		http.Error(w, "Like not found", http.StatusNotFound)
		return
	}

	// moram da znam autora fotke da bi znala je li banned
	//
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
		// Handle error
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if banned {
		http.Error(w, "Forbidden - the requester is banned by the author of the photo", http.StatusForbidden)
		return
	}
	// is requestor the author of the like
	//
	authorized, err := rt.db.IsAuthorizedToDeleteLike(like, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Like not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	if !authorized {
		// If user is not authorized to delete the like, send forbidden response
		http.Error(w, "Forbidden - the requester is not the author of the like", http.StatusForbidden)
		return
	}

	// delete like from the database
	err = rt.db.DeleteLike(like)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// Send success response
	w.WriteHeader(http.StatusNoContent)

}
