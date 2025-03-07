package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasaphoto-2009711/service/api/reqcontext"
)

// Function that removes a user from the followers of another user
func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id := extractToken(r.Header.Get("Authorization"))
	valid := isValid(ps.ByName("uid"), id)
	// Checking the user identity for the operation
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Checking if the user that will be unfollowed exists, if not returns HTTP status 404
	followId := ps.ByName("followingId")
	exists, err := rt.db.CheckUser(Userid{ID: followId}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("unfollowUser: error checking user existance")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		ctx.Logger.WithError(err).Error("unfollowUser: this user doesn't exists")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Checking if the user is currently followed, if not returns HTTP status 400
	isFollowed, err := rt.db.IsFollowed(Userid{ID: id}.toDatabase(), Userid{ID: followId}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("unfollowUser: error checking if one user is followed by the other one")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isFollowed {
		ctx.Logger.WithError(err).Error("unfollowUser: the user is already not followed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Calling the db function that removes the follow
	err = rt.db.UnfollowUser(Userid{ID: id}.toDatabase(), Userid{ID: followId}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("unfollowUser: error unfollowing the user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Returns the HTTP status 204
	w.WriteHeader(http.StatusNoContent)
}
