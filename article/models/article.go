package models

import (
	"fmt"
	"time"

	"github.com/chideat/pcc/article/modules/pig"
	"github.com/jinzhu/gorm"
)

func (article *Article) _BeforeSave() error {
	if article.UserId == 0 {
		return fmt.Errorf("invalid user_id")
	}
	if article.Data == "" {
		return fmt.Errorf("data is needed")
	}

	return nil
}

func (article *Article) Save() error {
	err := article._BeforeSave()
	if err != nil {
		return err
	}

	article.ModifiedUtc = time.Now().Local().UnixNano() / int64(time.Millisecond)
	if article.Id == 0 {
		article.Id = pig.Next(1, TYPE_ARTICLE)
		if err != nil {
			return err
		}
		article.CreatedUtc = time.Now().Local().UnixNano() / int64(time.Millisecond)

		err = db.Create(article).Error
	} else {
		err = db.Save(article).Error
	}

	return err
}

func (article *Article) Delete() error {
	article.Deleted = true
	article.DeletedUtc = time.Now().Local().UnixNano() / int64(time.Millisecond)

	return article.Save()
}

func (article *Article) Like() error {
	return db.Model(article).UpdateColumn("like_count", gorm.Expr("like_count+?", 1)).Error
}

func (article *Article) CancelLike() error {
	return db.Model(article).UpdateColumn("like_count", gorm.Expr("like_count-?", 1)).Error
}

func (article *Article) Map() (map[string]interface{}, error) {
	var (
		err error
		ret = map[string]interface{}{}
	)

	ret["id"] = article.Id
	ret["user_id"] = article.UserId
	ret["data"] = article.Data
	ret["like_count"] = article.LikeCount
	ret["created_utc"] = article.CreatedUtc
	ret["modified_utc"] = article.ModifiedUtc
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func GetArticleById(id int64) (*Article, error) {
	article := Article{}
	err := db.Where("deleted=false").First(&article, id).Error
	if err == ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &article, nil
	}
}

func GetArticles(cursor int64, limit int) ([]*Article, error) {
	articles := []*Article{}

	err := db.Where("deleted=false and created_utc>?", cursor).Order("created_utc asc").Limit(limit).Find(&articles).Error
	if err != nil {
		return nil, err
	} else {
		return articles, nil
	}
}
