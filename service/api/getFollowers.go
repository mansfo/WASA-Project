package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasaphoto-2009711/service/api/reqcontext"
)

// Function that returns the user id of the users that are following a specified user
func (rt *_router) getFollowers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	// The id of the requesting user
	id := extractToken(r.Header.Get("Authorization"))
	// The id of the user whose followers will be returned
	followersOf := ps.ByName("uid")

	// Calling the db funcion to get the followers
	count, followers, err := rt.db.GetFollowers(Userid{ID: id}.toDatabase(), Userid{ID: followersOf}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("getFollowers: error returning the followers")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Returns the HTTP status 200 and encodes the results of the db function
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(followers)
	_ = json.NewEncoder(w).Encode(count)
}
