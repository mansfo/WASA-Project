package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasaphoto-2009711/service/api/reqcontext"
)

// Function that returns the username of the users that like a specified photo
func (rt *_router) getLikeOnPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	id := extractToken(r.Header.Get("Authorization"))
	valid := isValid(ps.ByName("uid"), id)
	// Checking the user identity for the operation
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Returns the photo's identifier
	p, err := strconv.Atoi(ps.ByName("photo"))
	if err != nil {
		ctx.Logger.WithError(err).Error("GetLikeOnPhoto: error getting photo's identifier")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Calling the db function to get the likes
	count, likes, err := rt.db.GetLikeOnPhoto(Userid{ID: id}.toDatabase(), PhotoId{IdPhoto: p}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("GetLikeOnPhoto: error returning the likes on this photo")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Returns the HTTP status 200 and encodes the results of the db function
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(likes)
	_ = json.NewEncoder(w).Encode(count)
}
