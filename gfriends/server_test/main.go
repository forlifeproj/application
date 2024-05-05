package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	fllog "github.com/forlifeproj/msf/log"
	flsvr "github.com/forlifeproj/msf/server"
	"github.com/forlifeproj/protocol/gfriends/json/demo"
	"github.com/forlifeproj/protocol/gfriends/json/login"
	"github.com/forlifeproj/protocol/gfriends/json/order_manager"
)

var (
	cfg string
)

func init() {
	flag.StringVar(&cfg, "c", "../conf/server_test.toml", "config file path, default ../conf/server_test.toml")
}

func Mul(ctx context.Context, args *demo.Args, reply *demo.Reply) error {
	reply.C = args.A * args.B
	fllog.Log().Debug("req=", args, "reply=", reply)
	return nil
}

func Add(ctx context.Context, args *demo.Args, reply *demo.Reply) error {
	tt := login.LoginReq{}
	fllog.Log().Debug("token=", tt.Token)
	reply.C = args.A + args.B
	fllog.Log().Debug("req=", args, "reply=", reply)
	return nil
}

func GenerateOrderId(ctx context.Context, req *order_manager.GenerateOrderIdReq, rsp *order_manager.GenerateOrderIdRsp) error {
	rsp.OrderId = fmt.Sprintf("%d", time.Now().UnixNano())
	fllog.Log().Debug("req=", req, "rsp=", rsp)
	return nil
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
	svr.RegisterFunc(Mul)
	svr.RegisterFunc(Add)
	svr.StartServer()
	fllog.Log().Debug("test fllog debug cfg:", cfg)

}
