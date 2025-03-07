package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasaphoto-2009711/service/api/reqcontext"
)

// A function that adds a user in the banned list of the user that calls it
func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	id := extractToken(r.Header.Get("Authorization"))

	// Check user identity for the operation
	valid := isValid(ps.ByName("uid"), id)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Id of the user that is going to be banned
	banId := ps.ByName("bannedId")

	// Check if this user exists, if not return HTTP status 404
	exists, err := rt.db.CheckUser(Userid{ID: banId}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("banUser: error checking user existance")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		ctx.Logger.WithError(err).Error("banUser: this user doesn't exists")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Check if the user is already banned by banId, if so return HTTP status 403
	isBanned, err := rt.db.CheckBan(Userid{ID: id}.toDatabase(), Userid{ID: banId}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("banUser: error checking if one user is banned by the other user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if isBanned {
		ctx.Logger.WithError(err).Error("banUser: the user is banned by the other user")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Calling the db function to ban the user
	err = rt.db.BanUser(Userid{ID: id}.toDatabase(), Userid{ID: banId}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("banUser: error banning the user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return HTTP status 204
	w.WriteHeader(http.StatusNoContent)
}
