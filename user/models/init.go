package models

import (
	"github.com/chideat/glog"
	. "github.com/chideat/pcc/user/modules/cache"
	. "github.com/chideat/pcc/user/modules/config"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var (
	db    *gorm.DB
	cache *Cache

	ErrRecordNotFound = gorm.ErrRecordNotFound
)

func init() {
	var err error

	db, err = gorm.Open("postgres", Conf.Database)
	if err != nil {
		glog.Panic(err)
	}
	if Conf.IsDebug() {
		db.LogMode(true)
	} else {
		db.LogMode(false)
	}
	db.SingularTable(true)

	// sql table
	{
		if !db.HasTable(&User{}) {
			if err = db.CreateTable(&User{}).Error; err != nil {
				glog.Panic(err)
			}
		}
		db.AutoMigrate(&User{})

		db.AddIndex("idx_user_name", "name")
		db.AddIndex("idx_user_name_password", "name", "password")
	}

	// cache
	{
		cache = NewCache(Conf.Caches)
	}
}
