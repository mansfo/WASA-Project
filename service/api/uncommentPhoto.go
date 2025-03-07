package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasaphoto-2009711/service/api/reqcontext"
)

// Function that removes a comment from a photo
func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	id := extractToken(r.Header.Get("Authorization"))

	// Returns the identifier of the comment
	cid, err := strconv.ParseInt(ps.ByName("comment"), 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("uncommentPhoto: error getting the comment's identifier")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Checks if the user is the one that wrote the comments that he wants to remove, if not returns HTTP status 403
	isAuthor, err := rt.db.IsCommentAuthor(Userid{ID: id}.toDatabase(), CommentId{IdComment: cid}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("uncommentPhoto: error checking comment's author")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isAuthor {
		ctx.Logger.WithError(err).Error("uncommentPhoto: can't remove a comment published by someone else")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Calling the db function that removes the comment from the photo
	err = rt.db.UncommentPhoto(CommentId{IdComment: cid}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("uncommentPhoto: error removing the comment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Returns the HTTP status 204
	w.WriteHeader(http.StatusNoContent)
}
