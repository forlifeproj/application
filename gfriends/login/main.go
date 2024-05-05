package main

import (
	"flag"
	"fmt"

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
	// server init
	svr := flsvr.NewFLServer(cfg)
	svr.RegisterFunc(service.LoginIn)
	// svr.RegisterFunc(Add)
	svr.StartServer()
	fllog.Log().Debug("test fllog debug cfg:", cfg)

}
