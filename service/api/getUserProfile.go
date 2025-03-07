package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasaphoto-2009711/service/api/reqcontext"
)

// Function that returns the profile of a user
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	// The userid of the requesting user and of the user whose profile will be returned
	requestingUser := extractToken(r.Header.Get("Authorization"))
	id := ps.ByName("uid")

	// Calling the db function to return the profile
	profile, err := rt.db.GetProfile(Userid{ID: requestingUser}.toDatabase(), Userid{ID: id}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile: error executing the query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check if the user is banned
	isbanned, err := rt.db.CheckBan(Userid{ID: id}.toDatabase(), Userid{ID: requestingUser}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile/CheckBan: error executing the query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if isbanned {
		w.WriteHeader(http.StatusPartialContent)
		// Returns partial content (HTTP status 206) as the profile is banned by the user who called the function
		_ = json.NewEncoder(w).Encode(profile)
	} else {
		// Returns the HTTP status 200 and encodes the profile
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(profile)
	}
}
