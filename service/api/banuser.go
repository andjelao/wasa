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
	// "fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	// "git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	// "fantastic-coffee-decaffeinated/service/database"
	// "fantastic-coffee-decaffeinated/service/api/models"
)

// type Ban struct {
// bannedUsername string `json:"bannedUsername"`
// }

// type ConflictResponseBan struct {
// Message string `json:"message"`
// ExistingBan rt.db.Ban  `json:"existingBan"`
// }

// ban a user

func (rt *_router) BanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
	var ban database.Ban
	if err := json.NewDecoder(r.Body).Decode(&ban); err != nil {
		http.Error(w, "Failed to decode JSON body", http.StatusBadRequest)
		return
	}

	// check if the banned user is an authenticated user
	// da provjerim searched username da li je pravi user
	authenticated, err = rt.db.IsAuthenticatedUser(ban.BannedUsername)
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

	if ban.BannedUsername == username {
		http.Error(w, "Forbidden - user can not ban themselves", http.StatusForbidden)
		return
	}

	// check existing ban
	//
	banExists, err := rt.db.BanExists(username, ban.BannedUsername)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if banExists {
		// conflictResponseBan(w, "Conflict - Already banned this user", ban)
		response := database.ConflictResponseBan{
			Message:     "Conflict - Already banned this user",
			ExistingBan: ban,
		}

		jsonResponseBan, err := json.Marshal(response)
		if err != nil {
			// If there's an error in marshalling JSON, return an internal server error
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, err = w.Write(jsonResponseBan)
		if err != nil {
			// Handle error if writing response fails.
			http.Error(w, "Internal Server Error- error writing response", http.StatusInternalServerError)
			return
		}
		return
	}

	// sacuvaj ban
	//
	// Save the ban back to the database
	if err := rt.db.AddBan(username, ban.BannedUsername); err != nil {
		http.Error(w, "Failed to add ban in the database", http.StatusInternalServerError)
		return
	}
	// respond with a ban
	jsonResponse, err := json.Marshal(ban)
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

//
// func conflictResponseBan(w http.ResponseWriter, message string, ban  rt.db.Ban) {
// response := ConflictResponseBan{
// Message: message,
// ExistingBan: ban,
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
