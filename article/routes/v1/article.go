package v1

import (
	"strconv"

	"github.com/chideat/glog"
	"github.com/chideat/pcc/article/models"
	. "github.com/chideat/pcc/article/routes/utils"
	"github.com/gin-gonic/gin"
)

// Route: /articles
// Method: GET
func GetArticles(c *gin.Context) {
	count, _ := strconv.Atoi(c.Query("count"))
	if count < 0 || count > 100 {
		count = 10
	}
	cursor, _ := strconv.ParseUint(c.Query("cursor"), 10, 64)

	var (
		err      error
		articles = []*models.Article{}
	)
	articles, cursor, err = models.GetArticles(count, cursor)
	if err != nil {
		JsonWithError(c, "1", err)
		return
	}
	articlesMap := []map[string]interface{}{}
	for i := len(articles) - 1; i >= 0; i-- {
		feedMap, err := articles[i].Map()
		if err != nil {
			glog.Error(err)
			continue
		}
		articlesMap = append(articlesMap, feedMap)
	}
	JsonWithData(c, "0", "OK", articlesMap)
}

// Route: /articles/:id
// Method: GET
func GetArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Params.ByName("id"), 10, 64)
	if err != nil {
		Json(c, "1", "invalid article id")
		return
	}

	article, err := models.GetArticleById(id)
	if err != nil {
		JsonWithError(c, "1", err)
		return
	}
	if article == nil {
		Json(c, "1", "invalid id")
		return
	}
	ret, err := article.Map()
	if err != nil {
		JsonWithError(c, "1", err)
		return
	}
	JsonWithData(c, "0", "OK", ret)
}

func GetArticleLikeCount(c *gin.Context) {
	id, err := strconv.ParseUint(c.Params.ByName("id"), 10, 64)
	if err != nil {
		Json(c, "1", "invalid article id")
		return
	}
	count, err := models.GetArticleLikeCount(id)
	if err != nil {
		Json(c, "1", err.Error())
		return
	}
	JsonWithData(c, "0", "OK", map[string]int{"liked_count": count})
}

// Route: /articles
// Method: POST
func CreateArticle(c *gin.Context) {
	article := models.Article{}
	article.UserId, _ = strconv.ParseUint(c.PostForm("user_id"), 10, 64)
	article.Data = c.PostForm("data")

	err := article.Save()
	if err != nil {
		JsonWithError(c, "1", err)
		return
	}

	ret, err := article.Map()
	if err != nil {
		JsonWithError(c, "1", err)
		return
	}
	JsonWithData(c, "0", "OK", ret)
}

// Route: /articles/:id
// Method: PUT
func UpdateArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Params.ByName("id"), 10, 64)
	if err != nil {
		Json(c, "1", "invalid article id")
		return
	}

	article, err := models.GetArticleById(id)
	if err != nil {
		JsonWithError(c, "1", err)
		return
	}
	if article == nil {
		Json(c, "1", "invalid id")
		return
	}

	article.Data = c.DefaultPostForm("data", article.Data)

	ret, err := article.Map()
	if err != nil {
		JsonWithError(c, "1", err)
		return
	}
	JsonWithData(c, "0", "OK", ret)
}

// Route: /articles/:id
// Method: DELETE
func DeleteArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Params.ByName("id"), 10, 64)
	if err != nil {
		Json(c, "1", "invalid article id")
		return
	}

	article, err := models.GetArticleById(id)
	if err != nil {
		JsonWithError(c, "1", err)
		return
	}
	if article == nil {
		Json(c, "1", "invalid id")
		return
	}

	err = article.Delete()
	if err != nil {
		Json(c, "1", "delete article failed")
		return
	}
	Json(c, "0", "OK")
}
