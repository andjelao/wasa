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
	// "fantastic-coffee-decaffeinated/service/api/models"
	"database/sql"
	"errors"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	// "strconv"
	// "fantastic-coffee-decaffeinated/service/database"
)

// uunfollow user
func (rt *_router) unfollow(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	followedUsername := ps.ByName("followedUsername")
	if followedUsername == "" {
		http.Error(w, "followed username is required", http.StatusBadRequest)
		return
	}
	// check if the user exists
	authenticated, err = rt.db.IsAuthenticatedUser(followedUsername)
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

	//
	banned, err := rt.db.BanExists(followedUsername, username)
	if banned {
		http.Error(w, "Forbidden - not allowed access because of the ban", http.StatusForbidden)
		return
	}
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// check if follow object exists
	followExists, err := rt.db.FollowExists(username, followedUsername)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !followExists {
		http.Error(w, "Already not following this user", http.StatusNotFound)
		return
	}

	// delete follow from the database
	err = rt.db.Unfollow(username, followedUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	// Send success response
	w.WriteHeader(http.StatusNoContent)

}
