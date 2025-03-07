package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasaphoto-2009711/service/api/reqcontext"
)

// Function that adds a user in the followed list of the user that calls it
func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id := extractToken(r.Header.Get("Authorization"))
	valid := isValid(ps.ByName("uid"), id)
	// Checks user identity for the operation
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	followId := ps.ByName("followingId")
	// Check if the id of the user that will be followed exists, if not return HTTP status 404
	exists, err := rt.db.CheckUser(Userid{ID: followId}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("followUser: error checking user existance")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		ctx.Logger.WithError(err).Error("followUser: this user doesn't exists")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Check if one of the users is banned by the other one, if so returns HTTP status 403
	ban, err := rt.db.CheckBan(Userid{ID: followId}.toDatabase(), Userid{ID: id}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("followUser: error checking if one of the user is banned by the other one")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if ban {
		ctx.Logger.WithError(err).Error("followUser: one of the user is banned by the other one")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	ban1, err1 := rt.db.CheckBan(Userid{ID: id}.toDatabase(), Userid{ID: followId}.toDatabase())
	if err1 != nil {
		ctx.Logger.WithError(err).Error("followUser: error checking if one of the user is banned by the other one")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if ban1 {
		ctx.Logger.WithError(err).Error("followUser: one of the user is banned by the other one")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Calling the db function to follow the user
	err = rt.db.FollowUser(Userid{ID: id}.toDatabase(), Userid{ID: followId}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("followUser: error following the user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Returns HTTP status 204
	w.WriteHeader(http.StatusNoContent)
}
