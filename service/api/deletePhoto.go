package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"wasaphoto-2009711/service/api/reqcontext"
)

// Function that removes a photo from the profile of the user that published it
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	id := extractToken(r.Header.Get("Authorization"))
	valid := isValid(ps.ByName("uid"), id)
	// Check the user identity for the operation
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Returns the identifier of the photo
	photoId := ps.ByName("photo")
	pId, err := strconv.Atoi(photoId)
	if err != nil {
		ctx.Logger.WithError(err).Error("deletePhoto: error getting the photo's identifier")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Calling of the db function that removes the photo
	err = rt.db.RemovePhoto(PhotoId{IdPhoto: pId}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("deletePhoto: error removing the photo")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Returns the path of the local folder of the user to remove the photo from there
	pathPhoto, err := getUserPhotoPath(id)
	if err != nil {
		ctx.Logger.WithError(err).Error("deletePhoto: error getting the directory's path")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Removes the photo from the local directory
	err = os.Remove(filepath.Join(pathPhoto, photoId))
	if err != nil {
		ctx.Logger.WithError(err).Error("deletePhoto: error deleting the file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Returns the HTTP status 204
	w.WriteHeader(http.StatusNoContent)
}
