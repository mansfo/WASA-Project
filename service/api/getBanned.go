package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasaphoto-2009711/service/api/reqcontext"
)

// Function that returns the user id of all the users banned by the one that calls it
func (rt *_router) getBanned(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	id := extractToken(r.Header.Get("Authorization"))
	valid := isValid(ps.ByName("uid"), id)
	// Check user identity for the operation
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Calling the db function that returns the list of the users banned
	banned, err := rt.db.GetBanned(Userid{ID: id}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("getBanned: error returning the users currently banned by this user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Returns HTTP status 200 and encodes the result of the db function
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(banned)
}
