package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	// Login endpoint
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	// Profile endpoint
	rt.router.GET("/users/:uid/profile", rt.wrap(rt.getUserProfile))
	rt.router.PUT("/users/:uid/profile/username", rt.wrap(rt.setMyUserName))
	rt.router.PUT("/users/:uid/profile/profile_picture", rt.wrap(rt.setMyProfilePicture))
	rt.router.GET("/users/:uid/profile/profile_picture/:photo", rt.wrap(rt.getProfPict))

	// Search endpoint
	rt.router.GET("/users/:uid/search/:searchedname", rt.wrap(rt.searchUser))

	// Stream endpoint
	rt.router.GET("/users/:uid/stream", rt.wrap(rt.getMyStream))

	// Photos endpoint
	rt.router.GET("/users/:uid/photos/:photo", rt.wrap(rt.getPhoto))
	rt.router.POST("/users/:uid/photos", rt.wrap(rt.uploadPhoto))
	rt.router.DELETE("/users/:uid/photos/:photo", rt.wrap(rt.deletePhoto))

	// Likes endpoint
	rt.router.GET("/users/:uid/photos/:photo/likes", rt.wrap(rt.getLikeOnPhoto))
	rt.router.PUT("/users/:uid/photos/:photo/likes/:like", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/users/:uid/photos/:photo/likes/:like", rt.wrap(rt.unlikePhoto))

	// Comments endpoint
	rt.router.GET("/users/:uid/photos/:photo/comments", rt.wrap(rt.getComments))
	rt.router.POST("/users/:uid/photos/:photo/comments", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/users/:uid/photos/:photo/comments/:comment", rt.wrap(rt.uncommentPhoto))

	// Followers endpoint
	rt.router.GET("/users/:uid/followers", rt.wrap(rt.getFollowers))
	rt.router.GET("/users/:uid/following", rt.wrap(rt.getFollowing))
	rt.router.PUT("/users/:uid/following/:followingId", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:uid/following/:followingId", rt.wrap(rt.unfollowUser))

	// Bans endpoint
	rt.router.GET("/users/:uid/banned", rt.wrap(rt.getBanned))
	rt.router.PUT("/users/:uid/banned/:bannedId", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:uid/banned/:bannedId", rt.wrap(rt.unbanUser))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
