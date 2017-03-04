package models

import (
	"github.com/chideat/glog"
	. "github.com/chideat/pcc/action/modules/cache"
	. "github.com/chideat/pcc/action/modules/config"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	nsq "github.com/nsqio/go-nsq"
)

var (
	db                *gorm.DB
	cache             *Cache
	producer          *nsq.Producer
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

const (
	TYPE_USER    uint8 = 1
	TYPE_ACTION        = 2
	TYPE_ARTICLE       = 3
)

func init() {
	var err error

	{
		db, err = gorm.Open("postgres", Conf.Database)
		if err != nil {
			glog.Panic(err)
		}
		db.SingularTable(true)
		if Conf.IsDebug() {
			db.LogMode(true)
		} else {
			db.LogMode(false)
		}

		// TODO
		// single table may have problems in data increase in w.
		for _, table := range []interface{}{&LikeAction{}, &FollowAction{}} {
			if !db.HasTable(table) {
				if err = db.CreateTable(table).Error; err != nil {
					glog.Panic(err)
				}
			}
			db.AutoMigrate(table)
		}

		db.Model(&LikeAction{}).AddIndex("idx_like_target", "target")
		db.Model(&LikeAction{}).AddIndex("idx_like_user", "user_id")
		db.Model(&LikeAction{}).AddUniqueIndex("idx_like_user_target", "user_id", "target")

		db.Model(&FollowAction{}).AddIndex("idx_follow_target", "target")
		db.Model(&FollowAction{}).AddIndex("idx_follow_user", "user_id")
		db.Model(&FollowAction{}).AddUniqueIndex("idx_follow_user_target", "user_id", "target")
	}

	{
		config := nsq.NewConfig()
		producer, err = nsq.NewProducer(Conf.MQ.ProducerTCPAddress, config)
		if err != nil {
			panic(err)
		}
	}

	cache = NewCache(Conf.Caches)
}
