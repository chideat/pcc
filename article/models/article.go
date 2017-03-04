package models

import (
	"fmt"
	"time"

	"github.com/chideat/glog"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
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
		article.CreatedUtc = time.Now().Local().UnixNano() / int64(time.Millisecond)

		err = db.Create(article).Error
	} else {
		err = db.Save(article).Error
	}
	if err != nil {
		return err
	}

	go article.cache()

	return nil
}

func (article *Article) Delete() error {
	article.Deleted = true
	article.DeletedUtc = time.Now().Local().UnixNano() / int64(time.Millisecond)

	err := db.Save(article).Error
	if err != nil {
		return err
	}

	go article.cache()

	return nil
}

func (article *Article) cache() {
	key := fmt.Sprintf("index://articles/%d", article.Id)

	data, _ := proto.Marshal(article)
	_, err := cache.Do("SET", key, data)
	if err != nil {
		glog.Error(err)
	}
}

func (article *Article) Map() (map[string]interface{}, error) {
	var (
		err error
		ret = map[string]interface{}{}
	)

	ret["id"] = article.Id
	ret["user_id"] = article.UserId
	ret["data"] = article.Data
	ret["created_utc"] = article.CreatedUtc
	ret["modified_utc"] = article.ModifiedUtc

	ret["liked_count"], err = GetArticleLikeCount(article.Id)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	user, err := GetUserById(article.UserId)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	ret["user"] = user.Info()

	return ret, nil
}

func GetArticleById(id uint64) (*Article, error) {
	article := Article{}

	key := fmt.Sprintf("index://articles/%d", id)
	data, err := redis.Bytes(cache.Do("GET", key))
	if err == nil {
		err = proto.Unmarshal(data, &article)
		if err == nil {
			return &article, nil
		}
	}

	err = db.Where("deleted=false").First(&article, id).Error
	if err == ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		go article.cache()

		return &article, nil
	}
}

func GetArticles(count int, cursor uint64) ([]*Article, uint64, error) {
	articles := []*Article{}

	err := db.Where("deleted=false and created_utc>?", cursor).
		Order("created_utc asc").Limit(count).Find(&articles).Error
	if err != nil {
		return nil, 0, err
	}
	if len(articles) == count {
		cursor = uint64(articles[len(articles)-1].CreatedUtc)
	} else {
		cursor = 0
	}
	return articles, cursor, nil
}
