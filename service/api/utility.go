package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"wasaphoto-2009711/service/api/reqcontext"
)

// Function that checks if the identifier respects the length constraints
func validId(identifier string) bool {
	return (len(identifier) >= 3 && len(identifier) <= 16)
}

// Function that extracts the berear token from the Authorization header
func extractToken(token string) string {
	var bearer = strings.Split(token, " ")
	if len(bearer) == 2 {
		return strings.Trim(bearer[1], " ")
	}
	return ""
}

// Function that checks if the user has a valid token for a specific endpoint; if so returns 0, the error as a int (the HTTP status) otherwise
func isValid(id string, token string) int {
	if !logged(token) {
		return http.StatusForbidden
	}

	if id != token {
		return http.StatusUnauthorized
	}
	return 0
}

// Function that checks if a user is logged
func logged(token string) bool {
	return token != ""
}

// Function that returns the path of the local directory where the profile pictures are stored
func getUserProfPictPath(u string) (string, error) {
	path := filepath.Join(photosFolder, u, "profile_pictures")
	return path, nil
}

// Function that returns the path of the local directory where the photos published by the user are stored
func getUserPhotoPath(u string) (string, error) {
	path := filepath.Join(photosFolder, u, "photos")
	return path, nil
}

func createLocalFolder(id string, ctx reqcontext.RequestContext) error {

	// Create the path media/userid/ in the wasaproject directory
	path := filepath.Join(photosFolder, id)
	// Create the directory and insert in it a subdir called "photos" for the posts and a subdir called "profile_pictures" to store the profile pictures of the user
	err := os.MkdirAll(filepath.Join(path, "photos"), os.ModePerm)
	if err != nil {
		ctx.Logger.WithError(err).Error("session/createLocalFolder: can't create local folders")
	}
	err1 := os.MkdirAll(filepath.Join(path, "profile_pictures"), os.ModePerm)
	if err1 != nil {
		ctx.Logger.WithError(err).Error("session/createLocalFolder: can't create local folders")
	}
	return nil
}
