package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"time"
	"wasaphoto-2009711/service/api/reqcontext"
)

// Function that adds a comment on a photo
func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	id := extractToken(r.Header.Get("Authorization"))
	isBanned, err := rt.db.CheckBan(Userid{ID: id}.toDatabase(), Userid{ID: ps.ByName("uid")}.toDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("commentPhoto: error checking if the user is banned")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if isBanned {
		ctx.Logger.WithError(err).Error("commentPhoto: the user is banned")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var comment CommentBody
	// Reading and decoding the request body
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		ctx.Logger.WithError(err).Error("commentPhoto: error decoding the json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the comment length respects the limit, if not returns HTTP status 400
	if len(comment.Comment) > 300 || len(comment.Comment) < 1 {
		ctx.Logger.WithError(err).Error("commentPhoto: the comment doesn't match the length constraints")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Returns the identifier of the photo
	photo, err := strconv.Atoi(ps.ByName("photo"))
	if err != nil {
		ctx.Logger.WithError(err).Error("commentPhoto: error getting photo's identifier")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Calls the db function that adds the comment under the photo and returns the new comment
	comm, err := rt.db.CommentPhoto(Userid{ID: id}.toDatabase(), PhotoId{IdPhoto: photo}.toDatabase(), comment.Comment, time.Now().Format(time.RFC3339))
	if err != nil {
		ctx.Logger.WithError(err).Error("commentPhoto: error adding the comment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Returns HTTP status 201
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(comm)
}
