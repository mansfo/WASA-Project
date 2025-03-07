package api

import (
	"bytes"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"wasaphoto-2009711/service/api/reqcontext"
)

// Function that adds a photo in the profile of the user that calls it
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	id := extractToken(r.Header.Get("Authorization"))
	valid := isValid(ps.ByName("uid"), id)
	// Checking user identity for the operation
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// The request body can be read only one time, so it is stored in the variable data and it could be used again
	data, err := io.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("uploadPhoto: error reading the request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// The request body is refreshed with a buffer
	r.Body = io.NopCloser(bytes.NewBuffer(data))
	// A new buffer is used to read the request body
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("uploadPhoto: error reading the bytes buffer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var datetime = time.Now().Format("2006-01-02 15:04:05")
	// Calling the db function that adds a new image in the db
	imageId, err := rt.db.NewImage(datetime)
	if err != nil || imageId < 0 {
		ctx.Logger.WithError(err).Error("uploadPhoto: error creating a new instance in images table")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Returns the path of the local directory where the photos of the user are stored
	PhotoPath, err := getUserPhotoPath(id)
	if err != nil {
		ctx.Logger.WithError(err).Error("uploadPhoto: error returning the user's directory path")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Calling the db function that uploads a new photo in the profile of the user
	phId, err := rt.db.UploadPhoto(Userid{ID: id}.toDatabase(), ImageIdentifier{IdImage: imageId}.toDatabase(), datetime)
	if err != nil {
		ctx.Logger.WithError(err).Error("uploadPhoto: error uploading this photo")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var imagePath string
	imagePath = filepath.Join(PhotoPath, strconv.FormatInt(phId, 10))
	// Creates a new file in the local file to store the new image
	file, err := os.Create(imagePath)
	if err != nil {
		ctx.Logger.WithError(err).Error("uploadPhoto: error creating user's local photo file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Add the path of the image in the db table
	err = rt.db.SetImagePath(ImageIdentifier{IdImage: imageId}.toDatabase(), imagePath)
	if err != nil {
		ctx.Logger.WithError(err).Error("uploadPhoto: error setting the new image's path in the db")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// The request body is refreshed with a buffer
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	// The request body is copied in the new file
	_, err = io.Copy(file, r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("uploadPhoto: error copying the body into the file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// The new file is closed
	file.Close()

	// Returns HTTP status 204
	w.WriteHeader(http.StatusNoContent)
}
