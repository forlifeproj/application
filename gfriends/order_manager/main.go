package main

import (
	"flag"
	"fmt"

	"github.com/forlifeproj/application/gfriends/order_manager/dao"
	"github.com/forlifeproj/application/gfriends/order_manager/service"
	fllog "github.com/forlifeproj/msf/log"
	flsvr "github.com/forlifeproj/msf/server"
)

var (
	cfg string
)

func init() {
	flag.StringVar(&cfg, "c", "../conf/config.toml", "config file path, default ../conf/config.toml")
}

func main() {
	flag.Parse()
	// log init
	if err := fllog.Init(cfg); err != nil {
		fmt.Printf("fllog init failed. err:%+v", err)
		return
	}
	// db init
	if err := dao.Init(); err != nil {
		fllog.Log().Debug("init db failed, err:%v", err)
		return
	}
	// server init
	svr := flsvr.NewFLServer(cfg)
	svr.RegisterFunc(service.GenerateOrderId)
	svr.RegisterFunc(service.GetOrderList)
	svr.RegisterFunc(service.QueryOrderStatus)
	svr.StartServer()
	fllog.Log().Debug("test fllog debug cfg:", cfg)

}
