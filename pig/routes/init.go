package routes

import (
	"github.com/chideat/pcc/pig/routes/v1"
	"github.com/gin-gonic/gin"
)

var Handler *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)

	Handler = gin.New()

	Handler.GET("/api/v1/id/:type_id", v1.Auth, v1.Next)
}
