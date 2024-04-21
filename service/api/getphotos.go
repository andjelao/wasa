package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	// "fmt"
	// "log"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	// "io/ioutil"
	// "mime"
	// "time"
)

func (rt *_router) getPhotos(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Retrieve the value of two username query parameters from the request.
	username := r.URL.Query().Get("searchedUsername")
	myUsername := GetUsernameFromToken(r)

	// Check if the author exists in the database of authenticated users and is logged in
	authenticated, err := rt.db.IsAuthenticatedUser(myUsername)
	if err != nil {
		// Handle error if check fails
		http.Error(w, "Authentication check fail", http.StatusInternalServerError)
		return
	}
	if myUsername == "" || !authenticated {
		// Return error response if the author is not authenticated
		http.Error(w, "Unauthorized user", http.StatusUnauthorized)
		return
	}

	// da provjerim searched username da li je pravi user
	// nisam provjerila ako nije prazno
	if username != "" {
		authenticated, err = rt.db.IsAuthenticatedUser(username)
		if err != nil {
			// Handle error if check fails
			http.Error(w, "Authentication check fail", http.StatusInternalServerError)
			return
		}
		if !authenticated {
			// Return error response if the author is not authenticated
			http.Error(w, "Searched user not found", http.StatusNotFound)
			return
		}
	}
	// check if banned
	banned, err := rt.db.BanExists(username, myUsername)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if banned {
		http.Error(w, "Forbidden - not allowed access to other users' photos because of the ban", http.StatusForbidden)
		return
	}

	// If the username query parameter is not empty, use it to filter the photos
	// Otherwise, retrieve all photos.
	var photos []database.PhotoMultipart
	if username != "" {
		var users = []string{username}
		//
		photos, err = rt.db.RetrievePhotos(users, myUsername)
		// check if there are any photos present in the list
		if len(photos) == 0 {
			http.Error(w, "Photos not Found", http.StatusNotFound)
			return
		}
		if err != nil {
			// Handle error if retrieving photos from the database fails.
			http.Error(w, "Internal Server Error- error retrieving photos", http.StatusInternalServerError)
			return
		}
	} else {
		var users []string
		//
		photos, err = rt.db.RetrievePhotos(users, myUsername)
		if len(photos) == 0 {
			http.Error(w, "Photos not Found", http.StatusNotFound)
			return
		}
		if err != nil {
			// Handle error if retrieving photos from the database fails.
			http.Error(w, "Internal Server Error- error retrieving photos", http.StatusInternalServerError)
			return
		}
	}

	// Serialize the list of photos into JSON format.
	jsonResponse, err := json.Marshal(photos)
	if err != nil {
		// Handle error if serialization fails.
		http.Error(w, "Internal Server Error- error marshaling responce", http.StatusInternalServerError)
		return
	}

	// Set content type header.
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response back to the client.
	_, err = w.Write(jsonResponse)
	if err != nil {
		// Handle error if writing response fails.
		http.Error(w, "Internal Server Error- error writing responce", http.StatusInternalServerError)
		return
	}
}

//
// mozda da napravim samo jednu funkciju retrieve photos da uzme za parametar array usernameova koji isto moze biti prazan da uzmemo sve slike iz database
// jedan usere koje zelimo da vidimo a to su ili jedan user kojeg eksplicitno trazim jer gledam njegov profil ili moj following array
// i uvijek mi treba da pogledam dje se moj user nalazi u listama banned drugih usera da sakrijem njihove slike od njega
// retrievePhotosByAuthor is a placeholder function to simulate retrieving photos by au
