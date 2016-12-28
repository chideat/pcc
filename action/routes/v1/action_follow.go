package v1

/*

// Route: /users/:user_id/following
// Method: GET
func UserFollowing(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Params.ByName("user_id"), 10, 64)
	if err != nil {
		Json(c, "200001", "无效的用户ID")
		return
	}
	page, _ := strconv.ParseInt(c.DefaultPostForm("page", "1"), 10, 64)
	if page < 1 {
		page = 1
	}
	count, _ := strconv.ParseInt(c.DefaultPostForm("count", "100"), 10, 64)
	if count < 1 || count > 200 {
		count = 100
	}
	cursor, _ := strconv.ParseInt(c.Query("cursor"), 10, 64)

	var (
		ids   []int64
		users []map[string]interface{} = []map[string]interface{}{}
	)

	key := models.Index(userId, true, models.ActionType_Follow)
	for {
		ids, cursor, err = models.Range(key, (page-1)*count, count, cursor, true)
		if err != nil {
			glog.Error(err)
			Json(c, "100001", "系统错误")
			return
		}
		for _, id := range ids {
			if int64(len(users)) == count {
				cursor = id
				break
			}

			user, err := user.UserBaseInfo(id)
			if err != nil {
				glog.Error(err)
				continue
			}

			// get user extra info
			info := models.GetUserExtraInfo(id, 0)
			for k, v := range info {
				user[k] = v
			}
			user["followed"] = true

			users = append(users, user)
		}
		if cursor == 0 || int64(len(users)) == count {
			break
		}
	}

	if cursor == 0 {
		page = -1
	}

	JsonWithDataInfo(c, "0", "OK", users, map[string]interface{}{
		"cursor":    fmt.Sprintf("%d", cursor),
		"count":     count,
		"next_page": page + 1,
	})
}

// Route: /users/:user_id/follow
// Method: POST
func FollowUser(c *gin.Context) {
	access, err := ParseAccess(c, false)
	if err != nil {
		Json(c, "200001", "缺少参数")
		return
	}
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		Json(c, "200001", "请先登录")
		return
	}
	target, err := strconv.ParseInt(c.Params.ByName("user_id"), 10, 64)
	if err != nil {
		Json(c, "200001", "无效的用户ID")
		return
	}

	action, _ := models.GetActionByTargetAndUser(target, userId, models.ActionType_Follow)
	if action != nil {
		Json(c, "200001", "您已经关注了此用户")
		return
	}

	action = models.NewAction(models.ActionType_Follow, access.App, access.DeviceId, access.UserId, target, "")
	action.Ip = access.Ip
	action.Net = access.Net

	dataRaw, err := proto.Marshal(&models.Request{Method: "ADD", Action: action})
	if err != nil {
		Json(c, "100001", "关注用户失败")
		return
	}

	err = ActionProducer.Publish(ACTION_TOPIC, dataRaw)
	if err == nil {
		Json(c, "0", "OK")
	} else {
		glog.Error(err)
		Json(c, "100001", "关注用户失败")
	}
}

// Route: /users/:user_id/sections/:section_id
// Method: DELETE
func UnfollowUser(c *gin.Context) {
	_, err := ParseAccess(c, false)
	if err != nil {
		Json(c, "200001", "缺少参数")
		return
	}
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		Json(c, "200001", "请先登录")
		return
	}
	target, err := strconv.ParseInt(c.Params.ByName("user_id"), 10, 64)
	if err != nil {
		Json(c, "200001", "无效的用户ID")
		return
	}

	action, err := models.GetActionByTargetAndUser(target, userId, models.ActionType_Follow)
	if err != nil {
		glog.Error(err)
		Json(c, "100001", "系统错误")
		return
	}

	if action == nil {
		action = &models.Action{Id: target}
	}
	dataRaw, err := proto.Marshal(&models.Request{Method: "DEL", Action: action})
	if err != nil {
		glog.Error(err)
		Json(c, "100001", "取消关注失败")
		return
	}

	err = ActionProducer.Publish(ACTION_TOPIC, dataRaw)
	if err == nil {
		Json(c, "0", "OK")
	} else {
		glog.Error(err)
		Json(c, "100001", "取消关注失败")
	}
}

*/
