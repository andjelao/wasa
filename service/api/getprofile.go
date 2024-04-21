package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	// "fmt"
	// "log"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
)

// type Profile struct{
// username string `json:"username"`
// photosCount int `json:"photosCount"`
// followersCount int `json:"followersCount"`
// followersList []Follower `json:"followersList"`
// followingCount int `json:"followingCount"`
// followingList []Followed `json:"followingList"`
// userPhotos []Photo `json:"userPhotos"`
// }

// get user profile
func (rt *_router) getProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := GetUsernameFromToken(r)
	authenticated, err := rt.db.IsAuthenticatedUser(username)
	if err != nil {
		// Handle error if check fails
		http.Error(w, "Authentication check fail", http.StatusInternalServerError)
		return
	}
	// Check if the author exists in the database of authenticated users and is logged in
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

	// check if the user exists
	authenticated, err = rt.db.IsAuthenticatedUser(pathusername)
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
	// if asking for other users
	if username != pathusername {
		// check if banned
		banned, err := rt.db.BanExists(pathusername, username)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if banned {
			http.Error(w, "Forbidden - not allowed access because of the ban", http.StatusForbidden)
			return
		}
	}

	// get photostream from database
	var profile database.Profile
	profile, err = rt.db.GetProfile(pathusername)

	if err != nil {
		// Handle error if retrieving profile from the database fails.
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Serialize the profile into JSON format.
	jsonResponse, err := json.Marshal(profile)
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
