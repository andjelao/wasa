package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	// "fmt"
	// "log"
	// "io/ioutil"
	// "mime"
	// "time"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	// "fantastic-coffee-decaffeinated/service/api/models"
)

// get a list of banned users
func (rt *_router) getBanned(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := GetUsernameFromToken(r)
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
		http.Error(w, "Forbidden - user can not access other users bannned list", http.StatusForbidden)
		return
	}

	var banned []database.Ban
	// get list of followers of this user from database
	banned, err = rt.db.GetBanned(pathusername)
	if len(banned) == 0 {
		http.Error(w, "User has not banned anyone", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Serialize the list of likes into JSON format.
	jsonResponse, err := json.Marshal(banned)
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
