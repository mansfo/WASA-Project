package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasaphoto-2009711/service/api/reqcontext"
)

// Function that returns the stream of the user that calls it
func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	id := extractToken(r.Header.Get("Authorization"))
	valid := isValid(ps.ByName("uid"), id)
	// Checking the user identity for the operation
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Calling the db function to get the stream
	stream, err := rt.db.GetMyStream(Userid{ID: id}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("getMyStream: error returning the stream from the database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Returns the HTTP status 200 and encodes the stream
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(stream)
}
