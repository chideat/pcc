package models

import (
	"github.com/chideat/glog"
	. "github.com/chideat/pcc/action/modules/config"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var (
	db                *gorm.DB
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

func init() {
	var err error

	db, err = gorm.Open("postgres", Conf.Database)
	if err != nil {
		glog.Panic(err)
	}
	db.SingularTable(true)
	if Conf.Model == "debug" {
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

	db.Model(&LikeAction{}).AddUniqueIndex("idx_like_user_target", "user_id", "target")
}
