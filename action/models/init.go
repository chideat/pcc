package models

import (
	"github.com/chideat/glog"
	. "github.com/chideat/pcc/action/modules/config"
	"github.com/jinzhu/gorm"
)

var (
	db                *gorm.DB
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

func init() {
	var err error

	if db, err = gorm.Open("postgres", Config.Database); err == nil {
		db.SingularTable(true)
		db.LogMode(false)
	} else {
		glog.Panic(err)
	}

	// TODO
	// single table may have problems in data increase in w.
	if !db.HasTable(&LikeAction{}) {
		if err = db.CreateTable(&LikeAction{}).Error; err != nil {
			glog.Panic(err)
		}
	}
	db.AutoMigrate(&LikeAction{})
}
