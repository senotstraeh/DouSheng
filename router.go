package main

import (
	"DouSheng/controller"
	"github.com/gin-gonic/gin"
)

func initDouShengRouter() *gin.Engine {
	r := gin.Default()
	apiRouter := r.Group("/douyin")
	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.POST("/publish/action/", controller.Publish)
	apiRouter.GET("/publish/list/", controller.PublishList)
	apiRouter.GET("/user/", controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	// extra apis - I
	apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.GetFavouriteList)
	apiRouter.POST("/comment/action/", controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)
	// extra apis - II
	//apiRouter.POST("/relation/action/", jwt.Auth(), controller.RelationAction)
	//apiRouter.GET("/relation/follow/list/", jwt.Auth(), controller.GetFollowing)
	//apiRouter.GET("/relation/follower/list", jwt.Auth(), controller.GetFollowers)

	return r
}
