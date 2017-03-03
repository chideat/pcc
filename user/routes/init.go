package routes

import (
	"github.com/chideat/pcc/user/routes/v1"
	"github.com/gin-gonic/gin"
)

var Handler *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)

	Handler = gin.New()

	group := Handler.Group("/api/v1")
	{
		group.POST("/register", v1.Register)
	}
}
