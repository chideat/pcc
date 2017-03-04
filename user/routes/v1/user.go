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
		Id   uint64 `form:"id" json:"id" binding:"required"`
		Name string `form:"name" json:"name" binding:"required"`
	}

	err := c.Bind(&req)
	if err != nil {
		Json(c, "200001", "缺少参数")
		return
	}

	user, err := models.NewUser(req.Id, req.Name, "pcc")
	if err != nil {
		JsonWithError(c, err)
		return
	}
	JsonWithData(c, "0", "OK", user.Info())
}
