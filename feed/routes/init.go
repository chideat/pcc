package routes

import (
	"github.com/chideat/pcc/feed/routes/v1"
	"github.com/gin-gonic/gin"
)

var Handler *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)

	Handler = gin.New()

	group := Handler.Group("/api/v1")
	{
		group.GET("/feeds", v1.GetFeeds)
		group.GET("/feeds/:id", v1.GetFeed)
		group.POST("/feeds", v1.CreateFeed)
		group.PUT("/feeds/:id", v1.UpdateFeed)
		group.DELETE("/feeds/:id", v1.DeleteFeed)
	}
}
