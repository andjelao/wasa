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
	// "fantastic-coffee-decaffeinated/service/database"
	// "fantastic-coffee-decaffeinated/service/api/models"
)

// check da li imamo required parametre

// deletePhoto deletes the photo with the given photo ID.
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Extract photo ID from the request parameters
	photoIDD := ps.ByName("photoId")
	if photoIDD == "" {
		http.Error(w, "PhotoID is required", http.StatusBadRequest)
		return
	}
	photoID, err := strconv.Atoi(ps.ByName("photoId"))
	if err != nil {
		http.Error(w, "Internal server error-error converting to string", http.StatusInternalServerError)
		return
	}
	// Extract username from the query parameters
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

	// Check if the user is authorized to delete the photo
	//
	authorized, err := rt.db.IsAuthorized(username, photoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Photo not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	if !authorized {
		// If user is not authorized to delete the photo, send forbidden response
		http.Error(w, "Forbidden - the requester is not the author of the photo", http.StatusForbidden)
		return
	}

	// Delete the photo from the database
	//
	err = rt.db.DeletePhoto(photoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Photo not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	// Delete associated likes and comments
	//
	// gets deleted by the design of the database

	// Send success response
	w.WriteHeader(http.StatusNoContent)
}
