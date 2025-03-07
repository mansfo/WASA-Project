package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasaphoto-2009711/service/api/reqcontext"
)

// Function that removes a like from a photo
func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	id := extractToken(r.Header.Get("Authorization"))
	valid := isValid(ps.ByName("like"), id)
	// Checking user identity for the operation
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Returns the identifier of the photo
	p, err := strconv.Atoi(ps.ByName("photo"))
	if err != nil {
		ctx.Logger.WithError(err).Error("unlikePhoto: error getting the photo's identifier")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Calling the db function that checks if the user likes the photo, if not returns HTTP status 400
	liked, err := rt.db.CheckLike(Userid{ID: id}.toDatabase(), PhotoId{IdPhoto: p}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("unlikePhoto: error checking if the users likes the photo")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if liked != 1 {
		ctx.Logger.WithError(err).Error("unlikePhoto: user can't unlike a photo that he doesn't like")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Calling the db function that removes the like from the photo
	err = rt.db.UnlikePhoto(Userid{ID: id}.toDatabase(), PhotoId{IdPhoto: p}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("unlikePhoto: error removing the like")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Returns HTTP status 204
	w.WriteHeader(http.StatusNoContent)
}
