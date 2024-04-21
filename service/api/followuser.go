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

// type Followed struct {
// following string `json:"following"`
// }
// type Follower struct {
// follower string `json:"follower"`
// }
// type ConflictResponseFollow struct {
// Message string `json:"message"`
// Followed  Followed  `json:"followed"`
// }

// follow a user
func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	// da procitam request body
	// Decode the JSON request body
	var follow database.Followed
	if err := json.NewDecoder(r.Body).Decode(&follow); err != nil {
		http.Error(w, "Failed to decode JSON body", http.StatusBadRequest)
		return
	}

	// check if the followed user is an authenticated user
	// da provjerim searched username da li je pravi user
	//
	authenticated, err = rt.db.IsAuthenticatedUser(follow.Following)
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

	// check if banned
	banned, err := rt.db.BanExists(follow.Following, username)
	if err != nil {
		// Handle error
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if banned {
		http.Error(w, "Forbidden - user can not send a follow request because of the ban", http.StatusForbidden)
		return
	}

	// check if user wants to follow themselves
	if follow.Following == username {
		http.Error(w, "Forbidden - user can not follow themselves", http.StatusForbidden)
		return
	}

	// check existing following
	//
	followExists, err := rt.db.FollowExists(username, follow.Following)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if followExists {
		// ConflictResponseFollow(w, , follow)
		response := database.ConflictResponseFollow{
			Message:  "Conflict - Already following this user",
			Followed: follow,
		}

		jsonResponseFollow, err := json.Marshal(response)
		if err != nil {
			// If there's an error in marshalling JSON, return an internal server error
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, err = w.Write(jsonResponseFollow)
		if err != nil {
			// Handle error if writing response fails.
			http.Error(w, "Internal Server Error- error writing response", http.StatusInternalServerError)
			return
		}
		return
	}

	// sacuvaj following
	// sacuvaj follower
	//
	// Save the follow back to the database
	if err := rt.db.AddFollow(username, follow.Following); err != nil {
		http.Error(w, "Failed to add follow in the database", http.StatusInternalServerError)
		return
	}
	// var follower Follower
	// follower.follower = username
	// if err := db.AddFollower(follower, follow.following); err != nil {
	// http.Error(w, "Failed to add follower in the database", http.StatusInternalServerError)
	// return
	// }
	// respond with a follow
	jsonResponse, err := json.Marshal(follow)
	if err != nil {
		// Handle error
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(jsonResponse)
	if err != nil {
		// Handle error if writing response fails.
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}

// da li vec postoji ovaj objekat u bazi podataka

//
// func ConflictResponseFollow(w http.ResponseWriter, message string, follow  Followed) {
// response := ConflictResponseFollow{
// Message: message,
// Followed: follow,
// }

// jsonResponse, err := json.Marshal(response)
// if err != nil {
// If there's an error in marshalling JSON, return an internal server error
// http.Error(w, "Internal server error", http.StatusInternalServerError)
// return
// }

// w.Header().Set("Content-Type", "application/json")
// w.WriteHeader(http.StatusConflict)
// w.Write(jsonResponse)
// }
