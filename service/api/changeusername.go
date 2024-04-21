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
)

type ChangeRequest struct {
	NewUsername string
}

// change username

func (rt *_router) changeusername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
	// NISAM PROVJERILA PATH USERNAME DA LI SE POKLAPA SA TOKENOM!!!!!
	pathusername := ps.ByName("username")
	if pathusername == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	if pathusername != username {
		http.Error(w, "Forbidden - username in bearer token does not match the provided path parameter username", http.StatusForbidden)
		return
	}

	// var newUsername string
	var newuser ChangeRequest
	if err := json.NewDecoder(r.Body).Decode(&newuser); err != nil {
		http.Error(w, "Failed to decode JSON body", http.StatusBadRequest)
		return
	}
	// fmt.Println(newuser)
	if newuser.NewUsername == "" {
		http.Error(w, "Invalid new username", http.StatusBadRequest)
		return
	}

	// check if the user exists
	authenticated, err = rt.db.IsAuthenticatedUser(newuser.NewUsername)
	if err != nil {
		// Handle error if check fails
		http.Error(w, "Authentication check fail", http.StatusInternalServerError)
		return
	}
	if authenticated {
		// Return error response if the new username already exists
		http.Error(w, "Username already in use", http.StatusConflict)
		return
	}

	var new database.User
	//
	err = rt.db.ChangeUsername(username, newuser.NewUsername)
	if err != nil {
		// Handle error if retrieving profile from the database fails.
		// fmt.Println(err)
		http.Error(w, "Internal Server Error- error changing the username", http.StatusInternalServerError)
		return
	}
	new.Username = newuser.NewUsername
	new.Profile, err = rt.db.GetProfile(new.Username)
	if err != nil {
		// Handle error if retrieving user from the database fails.
		http.Error(w, "Internal Server Error- error getting the profile", http.StatusInternalServerError)
		return
	}

	new.Banned, err = rt.db.GetBanned(new.Username)
	if err != nil {
		// Handle error if retrieving user from the database fails.
		http.Error(w, "Internal Server Error- error getting banned users", http.StatusInternalServerError)
		return
	}

	new.Photostream, err = rt.db.PhotoStream(new.Username, "")
	if err != nil {
		// Handle error if retrieving user from the database fails.
		http.Error(w, "Internal Server Error- error getting photostream", http.StatusInternalServerError)
		return
	}

	response := UserReply{
		Identifier:   new.Username,
		UserResource: new,
	}

	// Serialize the profile into JSON format.
	jsonResponse, err := json.Marshal(response)
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
		http.Error(w, "Internal Server Error- error writing response", http.StatusInternalServerError)
		return
	}
}
