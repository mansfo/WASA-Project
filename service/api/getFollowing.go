package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasaphoto-2009711/service/api/reqcontext"
)

// Function that returns the user id of the users followed by a specified user
func (rt *_router) getFollowing(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	// The userid of the requesting user
	id := extractToken(r.Header.Get("Authorization"))
	// The userid of the user whose following will be returned
	followersOf := ps.ByName("uid")

	// Calling the db function to get the user followed by the user
	count, following, err := rt.db.GetFollowing(Userid{ID: id}.toDatabase(), Userid{ID: followersOf}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("getFollowing: error returning the users followed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Returns the HTTP status 200 and encodes the db function results
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(following)
	_ = json.NewEncoder(w).Encode(count)
}
