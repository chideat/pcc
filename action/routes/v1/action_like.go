package v1

import (
	"strconv"

	"github.com/chideat/glog"
	"github.com/chideat/pcc/action/models"
	. "github.com/chideat/pcc/action/modules/config"
	. "github.com/chideat/pcc/action/routes/utils"
	. "github.com/chideat/pcc/pig/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/nsqio/go-nsq"
)

var (
	producer *nsq.Producer
)

func init() {
	var err error

	config := nsq.NewConfig()
	producer, err = nsq.NewProducer(Config.Queue.NsqdAddress, config)
	if err != nil {
		glog.Panic(err)
	}
	producer.SetLogger(nil, nsq.LogLevelError)
}

// Route: /feeds/:id/like
// Method: POST
func FeedLike(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil || TYPE_USER&uint8(255&userId) == 0 {
		Json(c, "200001", "无效的用户ID")
		return
	}
	target, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
	if err != nil {
		Json(c, "200001", "无效的参数")
		return
	}
	mood, ok := models.LikeMood_value[c.Params.ByName("mood")]
	if !ok {
		Json(c, "200001", "invalid mood")
		return
	}

	action := models.NewLikeAction(userId, target, models.LikeMood(mood))
	err = action.Save()
	if err != nil {
		Json(c, "100001", "点赞失败")
		return
	}

	Json(c, "0", "OK")

	c.Next()

	data, _ := proto.Marshal(&models.Request{Method: models.RequestMethod_Add, Data: action.Bytes()})
	err = producer.Publish("pcc.action", data)
	if err != nil {
		glog.Error(err)
	}
}

// Route: /feeds/:id/like
// Method: POST
func FeedUnlike(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil || TYPE_USER&uint8(255&userId) == 0 {
		Json(c, "200001", "无效的用户ID")
		return
	}
	target, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
	if err != nil {
		Json(c, "200001", "无效的参数")
		return
	}

	action, err := models.GetLikeActionByUserAndTarget(userId, target)
	if err != nil {
		Json(c, "100001", "点赞失败")
		return
	}
	if action == nil {
		Json(c, "100001", "用户没有关注过")
		return
	}
	err = action.Delete()
	if err != nil {
		Json(c, "100001", "点赞失败")
		return
	}
	Json(c, "0", "OK")

	c.Next()

	data, _ := proto.Marshal(&models.Request{Method: models.RequestMethod_Delete, Data: action.Bytes()})
	err = producer.Publish("pcc.action", data)
	if err != nil {
		glog.Error(err)
	}
}
