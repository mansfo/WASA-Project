package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"path/filepath"
	"wasaphoto-2009711/service/api/reqcontext"
)

// Function that returns a photo published by a user
func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	http.ServeFile(w, r, filepath.Join(photosFolder, ps.ByName("uid"), "photos", ps.ByName("photo")))
}
