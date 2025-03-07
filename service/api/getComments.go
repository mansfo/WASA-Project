package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasaphoto-2009711/service/api/reqcontext"
)

// Function that returns the comments published on a photo
func (rt *_router) getComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	id := extractToken(r.Header.Get("Authorization"))

	// Returns the photo identifier
	photo, err := strconv.Atoi(ps.ByName("photo"))
	if err != nil {
		ctx.Logger.WithError(err).Error("GetComments: error getting the photo's identifier")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Calling the db function to get the comments
	count, comments, err := rt.db.GetComments(Userid{ID: id}.toDatabase(), PhotoId{IdPhoto: photo}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("GetComments: error returning the comments")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Returns HTTP status 200 and encodes the results of the db function
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(comments)
	_ = json.NewEncoder(w).Encode(count)
}
