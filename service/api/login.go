package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strings"
	"wasaphoto-2009711/service/api/reqcontext"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	var username Username

	// Reading and decoding the request body
	jsondata, err := io.ReadAll(r.Body)
	if err != nil {
		// Error reading the request body
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(jsondata, &username)
	if err != nil {
		// The body was not a valid json, so it is discarded
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !validId(username.NameUser) {
		// The userid is not valid
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Insert the user into the database
	user, isNew, err := rt.db.NewUser(username.NameUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("session: error checking if the user is a new user or wants to login")
		return
	}
	if !isNew {
		// The user is already in the database
		w.WriteHeader(http.StatusOK)
		// Encodes the user
		_ = json.NewEncoder(w).Encode(user)
		return
	}

	// Create user's folders locally
	err = createLocalFolder(strings.ToLower(username.NameUser), ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("session: can't create local folder")
		return
	}

	// Return the output
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("session: can't create response json")
		return
	}
}
