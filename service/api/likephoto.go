package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"

	"github.com/julienschmidt/httprouter"

	// "log"
	// "io/ioutil"
	// "mime"
	// "time"
	"database/sql"
	"errors"
	"strconv"
	// "fantastic-coffee-decaffeinated/service/api/models"
)

//
// type Like struct {
// username  string `json:"username"`
// photoId  int `json:"photoId"`
// }
//
// ConflictResponse represents the structure of the response for a 409 Conflict for Liking action
// type ConflictResponse struct {
// string `json:"message"`
// Like     Like   `json:"like"`
// }

// like a photo
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	// Extract the photoId from the path parameters
	photoIdd := ps.ByName("photoId")
	if photoIdd == "" {
		http.Error(w, "Photo ID is required", http.StatusBadRequest)
		return
	}
	photoId, err := strconv.Atoi(ps.ByName("photoId"))
	if err != nil {
		http.Error(w, "Internal server error- error converting photoid to string", http.StatusInternalServerError)
		return
	}

	// create Like object with obtained properties
	var like database.Like
	like.Username = username
	like.PhotoId = photoId

	// check if the like already exists for this photo
	//
	likeExists, err := rt.db.LikeExists(like)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if likeExists {
		// ConflictResponse(w, "Conflict - Like already exists", like)
		response := database.ConflictResponseLike{
			Message: "Conflict - Like already exists",
			Like:    like,
		}

		jsonResponseLike, err := json.Marshal(response)
		if err != nil {
			// If there's an error in marshalling JSON, return an internal server error
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)

		_, err = w.Write(jsonResponseLike)
		if err != nil {
			// Handle error if writing response fails.
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	// check if the user is allowed to like a photo - cant like your own and must not be banned by the author of the photo
	isauthor, err := rt.db.IsAuthorized(username, photoId)
	if isauthor {
		http.Error(w, "Forbidden - the requester is the author of the photo", http.StatusForbidden)
		return
	}
	if err != nil {
		// Handle error
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Photo not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	// moram da znam autora fotke da bi znala je li banned
	//
	author, err := rt.db.ExtractAuthor(photoId)
	if err != nil {
		// Handle error
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Photo not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	//
	banned, err := rt.db.BanExists(author, username)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if banned {
		http.Error(w, "Forbidden - the requester is banned by the author of the photo", http.StatusForbidden)
		return
	}

	// database ubaci lajk u strukturu za ovu fotku
	// Assuming you have a function to add a like to a photo
	err = rt.db.AddLike(like)
	if err != nil {
		// Handle error
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	jsonResponse, err := json.Marshal(like)
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

// func ConflictResponse(w http.ResponseWriter, message string, like  Like) {
// response := ConflictResponse{
// Message: message,
// Like:    like,
//   }

// jsonResponse, err := json.Marshal(response)
// if err != nil {
// If there's an error in marshalling JSON, return an internal server error
//    http.Error(w, "Internal server error", http.StatusInternalServerError)
// return
// }

//   w.Header().Set("Content-Type", "application/json")
// w.WriteHeader(http.StatusConflict)
// w.Write(jsonResponse)
// }
