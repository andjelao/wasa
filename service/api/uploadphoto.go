package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	// "fmt"
	// "log"
	"io/ioutil"
	// "mime"
	"encoding/base64"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"

	// "fantastic-coffee-decaffeinated/service/api/models"
	"math/rand"
)

// //////////////////////////////
// type PhotoMultipart struct {
// /////////////// koji format da stavim da li jpg ili ovako bytes
// photo 		 []byte `json:"photo"`
// ///////////////////////////////
// photoId           int    `json:"photoId"`
// author       string `json:"author"`
// uploadDateTime   string `json:"uploadDateTime"`
// location     string `json:"location,omitempty"`
// caption      string `json:"caption,omitempty"`
// likesCount        int    `json:"likesCount,omitempty"`
// likes       []Like `json:"likes,omitempty"`
// commentsCount     int    `json:"commentsCount,omitempty"`
// comments     []Comment `json:"comments,omitempty"`
// }
// //////////////////////////////
type UpdateResponse struct {
	Message string                  `json:"message"`
	Photo   database.PhotoMultipart `json:"photo"`
}

func isValidFileType(contentType string) bool {
	return contentType == "image/jpeg" || contentType == "image/png"
}

// ostalo da validiram length stringova location and photo
// da implementujem is authenticated user
// da implementujem upload photo in the database

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parse multipart form data
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
		return
	}
	// Get photo file from form data
	file, _, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Error retrieving photo from form data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read photo file into byte slice
	photoData, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading photo data", http.StatusInternalServerError)
		return
	}

	// Check if the photo data is empty
	if len(photoData) == 0 {
		http.Error(w, "Empty photo data", http.StatusBadRequest)
		return
	}

	// Determine file type
	fileType := http.DetectContentType(photoData)
	if !isValidFileType(fileType) {
		http.Error(w, "Unsupported file type", http.StatusBadRequest)
		return
	}

	// Encode photo data using base64
	encodedPhoto := base64.StdEncoding.EncodeToString(photoData)

	// Get other form fields
	// treba li da provjeravam validnost stringova paterne duzinu i ostalo
	caption := r.FormValue("caption")
	location := r.FormValue("location")

	author := GetUsernameFromToken(r)

	// Check if the author exists in the database of authenticated users and is logged in
	authenticated, err := rt.db.IsAuthenticatedUser(author)
	if err != nil {
		// Handle error if check fails
		http.Error(w, "Authentication check fail", http.StatusInternalServerError)
		return
	}
	if author == "" || !authenticated {
		// Return error response if the author is not authenticated
		http.Error(w, "Unauthorized user", http.StatusUnauthorized)
		return
	}

	// Calculate upload date and time
	uploadDateTime := time.Now().Format(time.RFC3339)
	photoId := rand.Int()

	// Create PhotoMultipart object
	photo := database.PhotoMultipart{
		Photo:          encodedPhoto,
		PhotoId:        photoId,
		Caption:        caption,
		Location:       location,
		Author:         author,
		UploadDateTime: uploadDateTime,
		LikesCount:     0, // Default value for likes
		Likes:          make([]database.Like, 0),
		CommentsCount:  0, // Default value for comments
		Comments:       make([]database.Comment, 0),
	}

	// Do something with the photo data, such as saving it to a database or storage
	// kako da stavim u bazu
	// Insert photo into the database
	if err := rt.db.InsertPhoto(photo); err != nil {
		// Handle error if insertion fails
		http.Error(w, "Error inserting photo into the database", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	response := UpdateResponse{
		Message: "Photo updated successfully",
		Photo:   photo, // Assuming 'photo' is a PhotoMultipart type
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
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
