package main

import (
	"fmt"

	"github.com/forlifeproj/application/gfriends/order_manager/conf"
	"github.com/forlifeproj/application/gfriends/order_manager/dao"
	"github.com/forlifeproj/application/gfriends/order_manager/service"
	msfcfg "github.com/forlifeproj/msf/config"
	fllog "github.com/forlifeproj/msf/log"
	flsvr "github.com/forlifeproj/msf/server"
)

const (
	cfgPath = "../conf/config.toml"
)

func main() {
	if err := msfcfg.ParseConfigWithPath(&conf.GConf, cfgPath); err != nil {
		fmt.Printf("load config failed. err:%+v cfg:%s", err, cfgPath)
		return
	}
	// log init
	if err := fllog.Init(cfgPath); err != nil {
		fmt.Printf("fllog init failed. err:%+v", err)
		return
	}
	// db init
	if err := dao.Init(); err != nil {
		fllog.Log().Errorf("init db failed, err:%v", err)
		return
	}

	// wx pay init
	service.InitWxPayClient()

	// server init
	svr := flsvr.NewFLServer(cfgPath)
	svr.RegisterFunc(service.GenerateOrderId)
	svr.RegisterFunc(service.GetOrderList)
	svr.RegisterFunc(service.QueryOrderStatus)
	svr.RegisterFunc(service.WxPrePay)
	svr.RegisterFunc(service.WxPayCallback)
	svr.StartServer()
	fllog.Log().Debug("test fllog debug cfg:", cfgPath)

}
