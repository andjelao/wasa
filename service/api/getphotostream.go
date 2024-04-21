package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	// "fmt"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	// "log"
	// "io/ioutil"
	// "mime"
	// "time"
)

// get users photostream
func (rt *_router) getPhotoStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	// Retrieve the value of two username path parameter from the request and authentication token
	pathusername := ps.ByName("username")
	if pathusername == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	if pathusername != username {
		http.Error(w, "Forbidden - user is not allowed access other users' photo streams", http.StatusForbidden)
		return
	}

	sinceDateTime := r.URL.Query().Get("sinceDateTime")

	// get photostream from database
	var photos []database.PhotoMultipart
	photos, err = rt.db.PhotoStream(username, sinceDateTime)

	if err != nil {
		// Handle error if retrieving photostream from the database fails.
		http.Error(w, "Internal Server Error- error retrieving photos", http.StatusInternalServerError)
		return
	}

	// check if there are any photos present in the list
	if len(photos) == 0 {
		http.Error(w, "User not following anyone or followers have no photos", http.StatusNotFound)
		return
	}

	// Serialize the list of photos into JSON format.
	jsonResponse, err := json.Marshal(photos)
	if err != nil {
		// Handle error if serialization fails.
		http.Error(w, "Internal Server Error- error marshalling response", http.StatusInternalServerError)
		return
	}

	// Set content type header.
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response back to the client.
	_, err = w.Write(jsonResponse)
	if err != nil {
		// Handle error if writing response fails.
		http.Error(w, "Internal Server Error-error writing responce", http.StatusInternalServerError)
		return
	}
}
