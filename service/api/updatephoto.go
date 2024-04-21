package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"

	// "log"

	// "io/ioutil"
	// "mime"
	// "time"
	"database/sql"
	"errors"

	"strconv"
	// "fantastic-coffee-decaffeinated/service/api/models"
)

// type updateRequest struct {
// Caption  string `json:"caption"`
// Location string `json:"location"`
// }

// update photo attributes
func (rt *_router) updatePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	// Extract photo ID from the request parameters
	photoIDD := ps.ByName("photoId")
	if photoIDD == "" {
		http.Error(w, "Photo ID is required", http.StatusBadRequest)
		return
	}
	photoID, err := strconv.Atoi(ps.ByName("photoId"))
	if err != nil {
		// fmt.Println("Error converting string to int:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// Check if the user is authorized to update the photo
	authorized, err := rt.db.IsAuthorized(username, photoID)
	if !authorized {
		http.Error(w, "Forbidden - the requester is not the author of the photo", http.StatusForbidden)
		return
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Photo not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	// Decode the JSON request body
	var updateData database.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Failed to decode JSON body", http.StatusBadRequest)
		return
	}

	// update the photo from the database by id
	if err := rt.db.UpdatePhoto(photoID, updateData); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Photo not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update photo in the database", http.StatusInternalServerError)
		return
	}
	// fetch updated photo from the database
	var photo database.PhotoMultipart
	photo, err = rt.db.GetPhoto(photoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Photo not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// Return the updated photo in the response
	if err := json.NewEncoder(w).Encode(photo); err != nil {
		// Handle error if encoding fails
		http.Error(w, "Failed to encode photo object", http.StatusInternalServerError)
		return
	}

}
