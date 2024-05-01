package config

import (
	"fmt"

	msfcfg "github.com/forlifeproj/msf/config"
	fllog "github.com/forlifeproj/msf/log"
)

var Conf GateWayConfig

type RouterCfg struct {
	Url         string
	ServiceName string
	Method      string
	PassRatio   int `default:"100"`
}

type GateWayConfig struct {
	RouterTableMap map[string]RouterCfg
}

func Init(cfg string) error {
	if err := msfcfg.ParseConfigWithPath(&Conf, cfg); err != nil {
		fllog.Log().Debug(fmt.Sprintf("load config failed. err:%+v cfg:%s", err, cfg))
		return err
	}

	tmpRouterTableMap := make(map[string]RouterCfg)
	for _, v := range Conf.RouterTableMap {
		tmpRouterTableMap[v.Url] = v
	}
	Conf.RouterTableMap = tmpRouterTableMap
	fllog.Log().Debug(fmt.Sprintf("load config succ cfg:%s Conf:%+v", cfg, Conf))
	return nil
}
