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

	db, err = gorm.Open("postgres", Conf.Database)
	if err == nil {
		db.SingularTable(true)
	} else {
		glog.Panic(err)
	}

	if Conf.Model == "debug" {
		db.LogMode(true)
	} else {
		db.LogMode(false)
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
