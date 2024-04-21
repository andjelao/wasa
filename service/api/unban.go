package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	// "encoding/json"
	// "fmt"
	// "log"
	// "io/ioutil"
	// "mime"
	// "time"
	// "errors"
	// "database/sql"
	// "strconv"
	// "fantastic-coffee-decaffeinated/service/database"
	// "fantastic-coffee-decaffeinated/service/api/models"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
)

// unban user
func (rt *_router) unban(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
		http.Error(w, "Forbidden - username in bearer token does not match the provided path parameter username", http.StatusForbidden)
		return
	}

	bannedUsername := ps.ByName("bannedUsername")
	if bannedUsername == "" {
		http.Error(w, "banned username is required", http.StatusBadRequest)
		return
	}
	// check if the user exists
	authenticated, err = rt.db.IsAuthenticatedUser(bannedUsername)
	if err != nil {
		// Handle error if check fails
		http.Error(w, "Authentication check fail", http.StatusInternalServerError)
		return
	}
	if !authenticated {
		// Return error response if the author is not authenticated
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	// check if user wants to unban themselves
	if bannedUsername == username {
		http.Error(w, "Forbidden - user can not unban themselves", http.StatusForbidden)
		return
	}
	// check if ban exists
	exists, err := rt.db.BanExists(username, bannedUsername)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "User not banned", http.StatusNotFound)
		return
	}

	// delete ban from the database
	err = rt.db.Unban(username, bannedUsername)
	//
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// Send success response
	w.WriteHeader(http.StatusNoContent)
}
