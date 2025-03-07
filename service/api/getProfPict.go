package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"path/filepath"
	"wasaphoto-2009711/service/api/reqcontext"
)

// Function that returns the profile picture published by a user
func (rt *_router) getProfPict(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	http.ServeFile(w, r, filepath.Join(photosFolder, ps.ByName("uid"), "profile_pictures", ps.ByName("photo")))
}
