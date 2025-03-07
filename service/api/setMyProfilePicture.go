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

// Function that sets a new profile picture for the user that calls it
func (rt *_router) setMyProfilePicture(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
		ctx.Logger.WithError(err).Error("setMyProfilePicture: error reading the request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// The request body is refreshed with a buffer
	r.Body = io.NopCloser(bytes.NewBuffer(data))
	// A new buffer is used to read the request body and store the content as a string in the db
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("setMyProfilePicture: error reading the bytes buffer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var datetime = time.Now().Format("2006-01-02 15:04:05")
	// Calling the db function to set a new profile picture for the user
	imageId, err := rt.db.NewImage(datetime)
	if err != nil || imageId < 0 {
		ctx.Logger.WithError(err).Error("setMyProfilePicture: error creating a new instance in images table")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Returns the path of the local directory where the profile pictures of the user are stored
	Path, err := getUserProfPictPath(id)
	if err != nil {
		ctx.Logger.WithError(err).Error("setMyProfilePicture: error returning the user's directory path")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var pictPath string
	pictPath = filepath.Join(Path, strconv.FormatInt(imageId, 10))
	// Creates a new file in the local directory to store the new image
	file, err := os.Create(pictPath)
	if err != nil {
		ctx.Logger.WithError(err).Error("setMyProfilePicture: error creating user's local photo file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// The request body is refreshed with a buffer
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	// The request body is copied in the new file
	_, err = io.Copy(file, r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("setMyProfilePicture: error copying the body into the file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// The new file is closed
	file.Close()

	// Add the path of the image into the db table
	err = rt.db.SetImagePath(ImageIdentifier{IdImage: imageId}.toDatabase(), pictPath)
	if err != nil {
		ctx.Logger.WithError(err).Error("setMyProfilePicture: error setting the new image's path in the db")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Calling the db function to update the profile picture
	err = rt.db.SetMyProfilePict(Userid{ID: id}.toDatabase(), ImageIdentifier{IdImage: imageId}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("setMyProfilePicture: error uploading this photo")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Returns HTTP status 204
	w.WriteHeader(http.StatusNoContent)

}
