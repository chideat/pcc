package v1

import (
	"github.com/chideat/pcc/user/models"
	. "github.com/chideat/pcc/user/routes/utils"
	"github.com/gin-gonic/gin"
)

// Route: /register
// Method: POST
func Register(c *gin.Context) {
	var req struct {
		Name     string `form:"name" json:"name" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}

	err := c.Bind(&req)
	if err != nil {
		Json(c, "200001", "缺少参数")
		return
	}

	user, err := models.NewUser(req.Name, req.Password)
	if err != nil {
		JsonWithError(c, err)
		return
	}
	JsonWithData(c, "0", "OK", user.Info())
}
