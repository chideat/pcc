package routes

import (
	"github.com/chideat/pcc/action/routes/v1"
	"github.com/gin-gonic/gin"
)

var Handler *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)

	Handler = gin.New()

	group := Handler.Group("/api/v1")
	{
		group.GET("/articles/:id/liked_users", v1.GetLikedUsers)
		group.GET("/articles/:id/is_liked", v1.IsLiked)
		group.POST("/articles/:id/like", v1.DoLike)
		group.DELETE("/articles/:id/like", v1.UndoLike)

		group.GET("/users/:id/followers", v1.GetFollowers)
		group.GET("/users/:id/is_followed", v1.IsFollowed)
		group.POST("/users/:id/follow", v1.Follow)
		group.DELETE("/users/:id/follow", v1.Unfollow)
	}
}
