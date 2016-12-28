package v1

import (
	"strconv"

	"github.com/chideat/pcc/feed/models"
	. "github.com/chideat/pcc/feed/routes/utils"
	"github.com/gin-gonic/gin"
)

// Route: /feeds
// Method: GET
func GetFeeds(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 0 {
		page = 1
	}
	count, _ := strconv.Atoi(c.Query("count"))
	if count < 0 || count > 100 {
		count = 10
	}
	// cursor, _ := strconv.ParseInt(c.Query("cursor"), 10, 64)

}

// Route: /feeds/:id
// Method: GET
func GetFeed(c *gin.Context) {
	id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
	if err != nil {
		Json(c, "1", "invalid feed id")
		return
	}

	feed, err := models.GetFeedById(id)
	if err != nil {
		JsonWithError(c, "1", err)
		return
	}
	if feed == nil {
		Json(c, "1", "invalid id")
		return
	}
	ret, err := feed.Map()
	if err != nil {
		JsonWithError(c, "1", err)
		return
	}
	JsonWithData(c, "0", "OK", ret)
}

// Route: /feeds
// Method: POST
func CreateFeed(c *gin.Context) {
	feed := models.Feed{}
	feed.UserId, _ = strconv.ParseInt(c.PostForm("user_id"), 10, 64)
	feed.Data = c.PostForm("data")

	err := feed.Save()
	if err != nil {
		JsonWithError(c, "1", err)
		return
	}

	ret, err := feed.Map()
	if err != nil {
		JsonWithError(c, "1", err)
		return
	}
	JsonWithData(c, "0", "OK", ret)
}

// Route: /feeds/:id
// Method: PUT
func UpdateFeed(c *gin.Context) {
	id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
	if err != nil {
		Json(c, "1", "invalid feed id")
		return
	}

	feed, err := models.GetFeedById(id)
	if err != nil {
		JsonWithError(c, "1", err)
		return
	}
	if feed == nil {
		Json(c, "1", "invalid id")
		return
	}

	feed.Data = c.DefaultPostForm("data", feed.Data)

	ret, err := feed.Map()
	if err != nil {
		JsonWithError(c, "1", err)
		return
	}
	JsonWithData(c, "0", "OK", ret)
}

// Route: /feeds/:id
// Method: DELETE
func DeleteFeed(c *gin.Context) {
	id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
	if err != nil {
		Json(c, "1", "invalid feed id")
		return
	}

	feed, err := models.GetFeedById(id)
	if err != nil {
		JsonWithError(c, "1", err)
		return
	}
	if feed == nil {
		Json(c, "1", "invalid id")
		return
	}

	err = feed.Delete()
	if err != nil {
		Json(c, "1", "delete feed failed")
		return
	}
	Json(c, "0", "OK")
}
