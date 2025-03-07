package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasaphoto-2009711/service/api/reqcontext"
)

// Function that returns a list of users whose username or userid matches the one specified
func (rt *_router) searchUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	requestedUser := ps.ByName("searchedname")
	requestingUser := extractToken(r.Header.Get("Authorization"))

	// Calling the db function that returns the profile of the users
	users, err := rt.db.SearchUser(Userid{ID: requestingUser}.toDatabase(), Userid{ID: requestedUser}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("searchUser: error executing the query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Returns the HTTP status 200 and encodes the result
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(users)
}
