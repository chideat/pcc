package models

import (
	"github.com/chideat/glog"
	. "github.com/chideat/pcc/feed/modules/config"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var (
	db                *gorm.DB
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

func init() {
	var err error

	if db, err = gorm.Open("postgres", Conf.Database); err == nil {
		db.SingularTable(true)
		db.LogMode(false)
	} else {
		glog.Panic(err)
	}

	// TODO
	// single table may have problems in data increase in w.
	if !db.HasTable(&Feed{}) {
		if err = db.CreateTable(&Feed{}).Error; err != nil {
			glog.Panic(err)
		}
	}
	db.AutoMigrate(&Feed{})
}
