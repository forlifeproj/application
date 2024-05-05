package config

import (
	"fmt"

	msfcfg "github.com/forlifeproj/msf/config"
	fllog "github.com/forlifeproj/msf/log"
)

var GConf LoginConfig

type LoginConfig struct {
	DB struct {
		User     string
		PassWord string
		StrIp    string
		Port     int
		Database string
		MaxOpen  int
		MaxIdol  int
	}
}

func Init(cfg string) error {
	if err := msfcfg.ParseConfigWithPath(&GConf, cfg); err != nil {
		fllog.Log().Debug(fmt.Sprintf("load config failed. err:%+v cfg:%s", err, cfg))
		return err
	}

	fllog.Log().Debug(fmt.Sprintf("load config succ cfg:%s Conf:%+v", cfg, GConf))
	return nil
}
