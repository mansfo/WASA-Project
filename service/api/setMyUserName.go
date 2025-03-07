package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"wasaphoto-2009711/service/api/reqcontext"
)

// Function that sets a new username for the user that calls it
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	uid := ps.ByName("uid")
	requestingUser := extractToken(r.Header.Get("Authorization"))
	valid := isValid(uid, requestingUser)
	// Checking user identity for the operation
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	var newname Username
	// Decodes the request body to get the new username
	jsondata, err := io.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("setMyUserName: error reading the request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(jsondata, &newname)
	if err != nil {
		ctx.Logger.WithError(err).Error("setMyUserName: error decoding json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// If the length of the new username is < 3 or > 16 it is returned the HTTP status 400
	if len(newname.NameUser) > 16 || len(newname.NameUser) < 3 {
		ctx.Logger.WithError(err).Error("setMyUserName: the length of the new username doesn't respect the constraints")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if this username already exists, regardless of letter case
	exists, err := rt.db.CheckIfUsernameIsTaken(Userid{ID: requestingUser}.toDatabase(), newname.NameUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("setMyUserName: error checking if the username already exists")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if exists {
		ctx.Logger.WithError(err).Error("setMyUserName: this username already exists")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Calling the db function that sets a new username
	err = rt.db.SetMyUserName(Userid{ID: requestingUser}.toDatabase(), newname.NameUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("setMyUserName: error setting the username")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Returns the HTTP status 200
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(newname)
}
