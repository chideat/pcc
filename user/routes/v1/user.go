package v1

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/chideat/glog"
	"github.com/chideat/pcc/user/models"
	"github.com/chideat/pcc/user/modules/auth"
	. "github.com/chideat/pcc/user/modules/error"
	. "github.com/chideat/pcc/user/routes/utils"
	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {
	tokenStr := c.Request.Header.Get("Authorization")
	if tokenStr == "" {
		JsonWithError(c, NewUserError("202001", "未登录"))
		c.Abort()
		return
	}

	token, err := auth.Auth(tokenStr)
	if err != nil || !token.Valid {
		JsonWithError(c, NewUserError("202004", "授权失败，Token无效"))
		c.Abort()
		return
	}

	userInfoVal, ok := token.Claims["info"]
	if !ok {
		JsonWithError(c, NewUserError("202004", "授权失败，Token无效"))
		c.Abort()
		return
	}
	userInfoBytes, ok := userInfoVal.(string)
	if !ok {
		JsonWithError(c, NewUserError("202004", "授权失败，Token无效"))
		c.Abort()
		return
	}

	userInfo := auth.UserTokenInfo{}
	err = json.Unmarshal([]byte(userInfoBytes), &userInfo)
	if err != nil {
		JsonWithError(c, NewUserError("202004", "授权失败，Token无效"))
		c.Abort()
		return
	}

	if time.Now().After(userInfo.Expired) {
		JsonWithError(c, NewUserError("202005", "Token过期，请重新登录"))
		c.Abort()
		return
	}

	user, err := models.GetUserById(userInfo.ID)
	if err != nil {
		JsonWithError(c, NewUserError("100001", "系统错误"))
		c.Abort()
		return
	}
	if user == nil {
		JsonWithError(c, NewUserError("202004", "授权失败，Token无效"))
		c.Abort()
		return
	}
	if user.Ban == 1 {
		JsonWithError(c, NewUserError("202006", "您的账户已被禁言"))
		c.Abort()
		return
	}
}

// Route: /register
// Method: POST
func Register(c *gin.Context) {
	var req struct {
		Icon     string `form:"icon" json:"icon"`
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
	if req.Icon != "" {
		user.Icon = req.Icon
	}
	err = user.Save()
	if err != nil {
		glog.Error(err)
		JsonWithError(c, err)
		return
	}
	JsonWithData(c, "0", "OK", user.Info())
}

// Route: /login
// Method: POST
func Login(c *gin.Context) {
	var req struct {
		Name     string `form:"name" json:"name" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}

	if err := c.Bind(&req); err != nil {
		JsonWithError(c, NewUserError("200001", "缺少参数"))
		return
	}

	user, err := models.GetUserByName(req.Name)
	if err != nil {
		JsonWithError(c, NewUserError("100001", "系统错误"))
		return
	}

	if user == nil {
		JsonWithError(c, NewUserError("201002", "无效的用户名"))
		return
	}

	if !auth.PasswordValid(req.Password, user.Salt, user.Password) {
		JsonWithError(c, NewUserError("202003", "用户名或密码错误"))
		return
	}

	err = user.UpdateToken()
	if err != nil {
		glog.Error(err)
		JsonWithError(c, NewUserError("100001", "系统错误"))
		return
	}
	user.Login()
	err = user.Save()
	if err != nil {
		glog.Error(err)
		JsonWithError(c, NewUserError("100001", "系统错误"))
		return
	}

	info := user.Info()
	info["token"] = user.Token
	JsonWithData(c, "0", "OK", info)
}

// Route: /logout
// Method: POST
func Logout(c *gin.Context) {
	id, err := strconv.ParseInt(c.PostForm("user_id"), 10, 64)
	if err != nil {
		JsonWithError(c, NewUserError("100001", "非法请求"))
		return
	}

	user, err := models.GetUserById(id)
	if err != nil {
		JsonWithError(c, NewUserError("100001", "系统错误"))
		return
	}

	if user == nil {
		JsonWithError(c, NewUserError("201006", "无效的用户"))
		return
	}

	user.Logout()

	err = user.Save()
	if err != nil {
		glog.Error(err)
		JsonWithError(c, NewUserError("100001", "注销失败"))
		return
	}
	Json(c, "0", "OK")
}

// Route: /users/:user_id
// Method: GET
func GetUserInfo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Params.ByName("user_id"), 10, 64)
	if err != nil {
		Json(c, "200001", err.Error())
		return
	}

	user, err := models.GetUserById(id)
	if err != nil {
		JsonWithError(c, NewUserError("100001", "系统错误"))
		return
	}

	if user == nil {
		JsonWithError(c, NewUserError("201006", "无效的用户ID"))
		return
	}

	JsonWithData(c, "0", "OK", user.BaseInfo())
}
