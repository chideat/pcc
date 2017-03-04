package routes

import (
	"github.com/chideat/pcc/article/routes/v1"
	"github.com/gin-gonic/gin"
)

var Handler *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)

	Handler = gin.New()

	group := Handler.Group("/api/v1")
	{
		group.GET("/articles", v1.GetArticles)
		group.GET("/articles/:id", v1.GetArticle)
		group.GET("/articles/:id/liked_count", v1.GetArticleLikeCount)
		group.POST("/articles", v1.CreateArticle)
		group.PUT("/articles/:id", v1.UpdateArticle)
		group.DELETE("/articles/:id", v1.DeleteArticle)
	}
}
