package service

import (
	"github.com/forlifeproj/application/gfriends/order_manager/conf"
)

func IsProdEnv() bool {
	return conf.GConf.Global.Env == "prod"
}
