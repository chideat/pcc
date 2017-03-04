package v1

import (
	"strconv"

	"github.com/chideat/glog"
	"github.com/chideat/pcc/action/models"
	. "github.com/chideat/pcc/action/routes/utils"
	"github.com/gin-gonic/gin"
)

// Route: /articles/:id/liked_users
func GetLikedUsers(c *gin.Context) {
	count, _ := strconv.Atoi(c.Query("count"))
	if count == 0 {
		count = 20
	}
	target, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
	if err != nil {
		Json(c, "200001", "无效的参数")
		return
	}
	cursor, _ := strconv.ParseUint(c.Query("cursor"), 10, 64)

	actions := []*models.LikeAction{}
	actions, cursor, err = models.GetLikeActions(target, count, cursor)
	if err != nil {
		Json(c, "100001", err.Error())
		return
	}

	rets := []map[string]interface{}{}
	for _, action := range actions {
		ret, err := action.Map()
		if err != nil {
			glog.Error(err)
			continue
		}
		rets = append(rets, ret)
	}

	JsonWithDataInfo(c, "0", "OK", rets, map[string]interface{}{
		"cursor": cursor,
	})
}

// Route: /articles/:id/is_liked
// Method: POST
func IsLiked(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Query("user_id"), 10, 64)
	if err != nil {
		Json(c, "200001", "无效的用户ID")
		return
	}

	target, err := strconv.ParseUint(c.Params.ByName("id"), 10, 64)
	if err != nil {
		Json(c, "200001", "无效的参数")
		return
	}

	action, err := models.GetLikeAction(userId, target)
	if err != nil {
		Json(c, "100001", err.Error())
		return
	}

	if action == nil || action.Deleted {
		JsonWithData(c, "0", "OK", map[string]interface{}{"is_liked": false})
		return
	}
	JsonWithData(c, "0", "OK", map[string]interface{}{"is_liked": true})
}

// Route: /articles/:id/like
// Method: POST
func DoLike(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Query("user_id"), 10, 64)
	if err != nil {
		Json(c, "200001", "无效的用户ID")
		return
	}

	target, err := strconv.ParseUint(c.Params.ByName("id"), 10, 64)
	if err != nil {
		Json(c, "200001", "无效的参数")
		return
	}

	action, err := models.GetLikeAction(userId, target)
	if err != nil {
		Json(c, "100001", err.Error())
		return
	}

	if action != nil && !action.Deleted {
		Json(c, "200001", "已经点过赞了")
		return
	}

	if action != nil {
		action.Deleted = false
	} else {
		action, _ = models.NewLikeAction(userId, target)
	}

	err = action.Broadcast(models.RequestMethod_Add)
	if err != nil {
		glog.Error(err)
		Json(c, "10001", err.Error())
		return
	}
	ret, err := action.Map()
	if err != nil {
		glog.Error(err)
		Json(c, "10001", err.Error())
		return
	}
	JsonWithData(c, "0", "OK", ret)
}

// Route: /articles/:id/like
// Method: DELETE
func UndoLike(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Query("user_id"), 10, 64)
	if err != nil {
		Json(c, "200001", "无效的用户ID")
		return
	}
	target, err := strconv.ParseUint(c.Params.ByName("id"), 10, 64)
	if err != nil {
		Json(c, "200001", "无效的参数")
		return
	}

	action, err := models.GetLikeAction(userId, target)
	if err != nil {
		Json(c, "100001", "取消点赞失败")
		return
	}
	if action == nil || action.Deleted {
		Json(c, "200001", "未点赞")
		return
	}

	err = action.Broadcast(models.RequestMethod_Delete)
	if err != nil {
		glog.Error(err)
		Json(c, "10001", "取消点赞失败")
		return
	}
	Json(c, "0", "OK")
}
