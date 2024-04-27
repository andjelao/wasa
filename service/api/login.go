package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"

	// "fmt"
	// "log"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	// "io/ioutil"
	// "mime"
	// "time"
	// "fantastic-coffee-decaffeinated/service/api/models"
)

// type User struct{
// username string `json:"username"`
// profile Profile `json:"profile"`
// banned []Ban `json:"profile"`
// photostream []PhotoMultipart `json:"photostream"`
// }
type UserReply struct {
	Identifier   string        `json:"Identifier"`
	UserResource database.User `json:"UserResource"`
}
type LoginRequest struct {
	Username string `json:"username"`
}

// do login
func (rt *_router) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// da procitam request body
	// Decode the JSON request body
	var user database.User
	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Failed to decode JSON body", http.StatusBadRequest)
		return
	}
	// username := loginReq.Username
	if loginReq.Username == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// check if the user is an already authenticated user -- if yes he just logs in
	authenticated, err := rt.db.IsAuthenticatedUser(loginReq.Username)
	if err != nil {
		// Handle error if check fails
		if !errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Authentication check fail", http.StatusInternalServerError)
			return
		}
	}
	var exist bool = false
	if authenticated {
		//
		user, err = rt.db.Login(loginReq.Username)
		exist = true
	} else {
		//
		user, err = rt.db.CreateUser(loginReq.Username)
	}

	if err != nil {
		// Handle error if retrieving user from the database fails.
		http.Error(w, "Internal Server Error - error retrieving user from database", http.StatusInternalServerError)
		return
	}
	user.Profile, err = rt.db.GetProfile(user.Username)
	if err != nil {
		// Handle error if retrieving user from the database fails.
		http.Error(w, "Internal Server Error- error retrieving profile", http.StatusInternalServerError)
		return
	}

	user.Banned, err = rt.db.GetBanned(user.Username)
	if err != nil {
		// Handle error if retrieving user from the database fails.
		http.Error(w, "Internal Server Error- error retrieving banned users", http.StatusInternalServerError)
		return
	}

	user.Photostream, err = rt.db.PhotoStream(user.Username, "")
	if err != nil {
		// Handle error if retrieving user from the database fails.
		http.Error(w, "Internal Server Error- error retrieving photostream", http.StatusInternalServerError)
		return
	}
	response := UserReply{
		Identifier:   loginReq.Username,
		UserResource: user,
	}

	// Serialize the profile into JSON format.
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		// Handle error if serialization fails.
		http.Error(w, "Internal Server Error- error encoding response", http.StatusInternalServerError)
		return
	}
	// Set content type header.
	w.Header().Set("Content-Type", "application/json")
	if exist {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	// Write the JSON response back to the client.
	_, err = w.Write(jsonResponse)
	if err != nil {
		// Handle error if writing response fails.
		http.Error(w, "Internal Server Error- error writing response", http.StatusInternalServerError)
		return
	}

}
