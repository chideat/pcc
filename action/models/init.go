package models

import (
	"github.com/chideat/glog"
	. "github.com/chideat/pcc/action/modules/config"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	nsq "github.com/nsqio/go-nsq"
)

var (
	db                *gorm.DB
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
		if !db.HasTable(&LikeAction{}) {
			if err = db.CreateTable(&LikeAction{}).Error; err != nil {
				glog.Panic(err)
			}
		}
		db.AutoMigrate(&LikeAction{})

		db.Model(&LikeAction{}).AddIndex("idx_like_target", "target")
		db.Model(&LikeAction{}).AddIndex("idx_like_user", "user_id")
		db.Model(&LikeAction{}).AddUniqueIndex("idx_like_user_target", "user_id", "target")
	}

	{
		config := nsq.NewConfig()
		producer, err = nsq.NewProducer(Conf.MQ.ProducerTCPAddress, config)
		if err != nil {
			panic(err)
		}
		producer.SetLogger(nil, nsq.LogLevelError)
	}
}
