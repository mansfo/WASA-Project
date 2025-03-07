package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasaphoto-2009711/service/api/reqcontext"
)

// Function that adds the username of the one that calls it in the list of the users that like a certain photo
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	id := extractToken(r.Header.Get("Authorization"))
	valid := isValid(ps.ByName("like"), id)
	// Checking user identity for the operation
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Returns the photo's identifier
	photo, err := strconv.Atoi(ps.ByName("photo"))
	if err != nil {
		ctx.Logger.WithError(err).Error("LikePhoto: error getting the photo's identifier")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Calling the db function to like the photo
	err = rt.db.LikePhoto(Userid{ID: id}.toDatabase(), PhotoId{IdPhoto: photo}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("LikePhoto: error putting the like")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Returns the HTTP status 204
	w.WriteHeader(http.StatusNoContent)
}
