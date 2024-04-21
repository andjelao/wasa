package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	rt.router.POST("/photos/", rt.wrap(rt.uploadPhoto))
	rt.router.PATCH("/photos/:photoId", rt.wrap(rt.updatePhoto))
	rt.router.DELETE("/photos/:photoId", rt.wrap(rt.deletePhoto))
	rt.router.POST("/users/:username/banned/", rt.wrap(rt.BanUser))
	rt.router.GET("/users/:username/banned/", rt.wrap(rt.getBanned))
	rt.router.DELETE("/users/:username/banned/:bannedUsername", rt.wrap(rt.unban))
	rt.router.GET("/photos/:photoId/likes/", rt.wrap(rt.getLikes))
	rt.router.GET("/photos/:photoId/comments/", rt.wrap(rt.getComments))
	rt.router.POST("/photos/:photoId/comments/", rt.wrap(rt.commentPhoto))
	rt.router.POST("/photos/:photoId/likes/", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/photos/:photoId/likes/:username", rt.wrap(rt.dislike))
	rt.router.DELETE("/photos/:photoId/comments/:commentId", rt.wrap(rt.uncomment))
	rt.router.POST("/users/:username/followed/", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:username/followed/:followedUsername", rt.wrap(rt.unfollow))
	rt.router.GET("/users/:username/followed/", rt.wrap(rt.getFollowed))
	rt.router.GET("/users/:username/followers/", rt.wrap(rt.getFollowers))
	rt.router.POST("/session", rt.wrap(rt.login))
	rt.router.GET("/photos/", rt.wrap(rt.getPhotos))
	rt.router.GET("/users/:username/photo-stream", rt.wrap(rt.getPhotoStream))
	rt.router.GET("/users/:username/profile", rt.wrap(rt.getProfile))
	rt.router.PUT("/users/:username", rt.wrap(rt.changeusername))
	rt.router.GET("/tables", rt.wrap(rt.getTables))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}

// ovdje dodajes entities sa njihovim metodama
