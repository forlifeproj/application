package main

import (
	"flag"
	"fmt"

	"github.com/forlifeproj/application/gfriends/login/config"
	"github.com/forlifeproj/application/gfriends/login/dao"
	"github.com/forlifeproj/application/gfriends/login/service"
	fllog "github.com/forlifeproj/msf/log"
	flsvr "github.com/forlifeproj/msf/server"
)

var (
	cfg string
)

func init() {
	flag.StringVar(&cfg, "c", "../conf/login.toml", "config file path, default ../conf/login.toml")
}

func main() {
	flag.Parse()
	// log init
	if err := fllog.Init(cfg); err != nil {
		fmt.Printf("fllog init failed. err:%+v", err)
		return
	}
	// config init
	if err := config.Init(cfg); err != nil {
		fllog.Log().Error("config init failed. err=", err, " cfg=", cfg)
		return
	}
	// dao init
	if err := dao.Init(); err != nil {
		fllog.Log().Error("init database failed. err=", err)
		return
	}

	// server init
	svr := flsvr.NewFLServer(cfg)
	svr.RegisterFunc(service.LoginIn)
	svr.RegisterFunc(service.LoginAuth)
	svr.StartServer()
	fllog.Log().Debug("test fllog debug cfg:", cfg)

}
