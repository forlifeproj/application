package dao

import (
	"fmt"

	"github.com/forlifeproj/application/gfriends/order_manager/conf"
	fllog "github.com/forlifeproj/msf/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var gormDb *gorm.DB

func Init() error {
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.GConf.DB.User, conf.GConf.DB.PassWord, conf.GConf.DB.StrIp, conf.GConf.DB.Port, conf.GConf.DB.Database)
	fllog.Log().Debugf("mysql url: %s", url)

	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		fllog.Log().Errorf("init db error, err: %v", err)
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		fllog.Log().Errorf("init db error, err: %v", err)
		return err
	}

	sqlDB.SetMaxOpenConns(conf.GConf.DB.MaxOpen)
	sqlDB.SetMaxIdleConns(conf.GConf.DB.MaxIdol)

	gormDb = db
	fllog.Log().Debugf("init db success, url: %s, maxopen=%d, maxidle=%d", url, conf.GConf.DB.MaxOpen, conf.GConf.DB.MaxIdol)
	return nil
}
