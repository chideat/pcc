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
		group.GET("/feeds/:id/like/users", v1.GetFeedLikeUsers)
		group.POST("/feeds/:id/like", v1.FeedLike)
		group.DELETE("/feeds/:id/like", v1.FeedUnlike)
		// group.GET("/users/:user_id/following", v1.UserFollowing)
		// group.POST("/users/:user_id/follow", v1.Auth, v1.FollowUser)
		// group.DELETE("/users/:user_id/follow", v1.Auth, v1.UnfollowUser)
	}
}
