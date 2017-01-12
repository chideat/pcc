package models

import (
	"fmt"
	"time"

	. "github.com/chideat/pcc/pig/models"
	"github.com/chideat/pcc/sdk/pig"
	"github.com/chideat/pcc/sdk/user"
	"github.com/jinzhu/gorm"
)

func (feed *Feed) _BeforeSave() error {
	if feed.UserId == 0 {
		return fmt.Errorf("invalid user_id")
	}
	if feed.Data == "" {
		return fmt.Errorf("data is needed")
	}

	return nil
}

func (feed *Feed) Save() error {
	err := feed._BeforeSave()
	if err != nil {
		return err
	}

	feed.ModifiedUtc = time.Now().Local().UnixNano() / int64(time.Millisecond)
	if feed.Id == 0 {
		feed.Id, err = pig.Int64(TYPE_FEED)
		if err != nil {
			return err
		}
		feed.CreatedUtc = time.Now().Local().UnixNano() / int64(time.Millisecond)

		err = db.Create(feed).Error
	} else {
		err = db.Save(feed).Error
	}

	return err
}

func (feed *Feed) Delete() error {
	feed.Deleted = true
	feed.DeletedUtc = time.Now().Local().UnixNano() / int64(time.Millisecond)

	return feed.Save()
}

func (feed *Feed) Like() error {
	return db.Model(feed).UpdateColumn("like_count", gorm.Expr("like_count+?", 1)).Error
}

func (feed *Feed) CancelLike() error {
	return db.Model(feed).UpdateColumn("like_count", gorm.Expr("like_count-?", 1)).Error
}

func (feed *Feed) Map() (map[string]interface{}, error) {
	var (
		err error
		ret = map[string]interface{}{}
	)

	ret["id"] = feed.Id
	ret["user_id"] = feed.UserId
	ret["data"] = feed.Data
	ret["like_count"] = feed.LikeCount
	ret["created_utc"] = feed.CreatedUtc
	ret["modified_utc"] = feed.ModifiedUtc
	ret["user"], err = user.UserBaseInfo(feed.UserId)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func GetFeedById(id int64) (*Feed, error) {
	feed := Feed{}
	err := db.Where("deleted=false").First(&feed, id).Error
	if err == ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &feed, nil
	}
}

func GetFeeds(cursor int64, limit int) ([]*Feed, error) {
	feeds := []*Feed{}

	err := db.Where("deleted=false and created_utc>?", cursor).Order("created_utc asc").Limit(limit).Find(&feeds).Error
	if err != nil {
		return nil, err
	} else {
		return feeds, nil
	}
}
