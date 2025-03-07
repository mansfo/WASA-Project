package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasaphoto-2009711/service/api/reqcontext"
)

// Function that removes a ban that was put on a user
func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	id := extractToken(r.Header.Get("Authorization"))
	valid := isValid(ps.ByName("uid"), id)
	// Checking the user identity for the operation
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Checking if the user that will be unbanned exists, if not returns the HTTP status 404
	banId := ps.ByName("bannedId")
	exists, err := rt.db.CheckUser(Userid{ID: banId}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("unbanUser: error checking user existance")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		ctx.Logger.WithError(err).Error("unbanUser: this user doesn't exists")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Checking if the user is currently banned, if not returns HTTP status 400
	isBanned, err := rt.db.CheckBan(Userid{ID: banId}.toDatabase(), Userid{ID: id}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("unbanUser: error checking if the user is banned currently by the other user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isBanned {
		ctx.Logger.WithError(err).Error("unbanUser: the user is already not banned by the other user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Calling the db function that removes the ban
	err = rt.db.UnbanUser(Userid{ID: id}.toDatabase(), Userid{ID: banId}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("unbanUser: error unbanning the user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Returns the HTTP status 204
	w.WriteHeader(http.StatusNoContent)
}
