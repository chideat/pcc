package models

import (
	"github.com/chideat/glog"
	. "github.com/chideat/pcc/user/modules/config"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
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

	// create table
	if !db.HasTable(&User{}) {
		if err = db.CreateTable(&User{}).Error; err != nil {
			glog.Panic(err)
		}
	}
	db.AutoMigrate(&User{})
}
