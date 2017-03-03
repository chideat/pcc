package models

import (
	"github.com/chideat/glog"
	. "github.com/chideat/pcc/article/modules/cache"
	. "github.com/chideat/pcc/article/modules/config"
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

	{
		db, err = gorm.Open("postgres", Conf.Database)
		if err == nil {
			db.SingularTable(true)
		} else {
			glog.Panic(err)
		}
		if Conf.IsDebug() {
			db.LogMode(true)
		} else {
			db.LogMode(false)
		}

		// TODO
		// single table may have problems in data increase in w.
		if !db.HasTable(&Article{}) {
			if err = db.CreateTable(&Article{}).Error; err != nil {
				glog.Panic(err)
			}
		}
		db.AutoMigrate(&Article{})
	}

	{
		cache = NewCache(Conf.Caches)
	}
}
