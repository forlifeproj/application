package main

import (
	"context"
	"flag"
	"fmt"

	fllog "github.com/forlifeproj/msf/log"
	flsvr "github.com/forlifeproj/msf/server"

	"github.com/forlifeproj/protocol/gfriends/json/demo"
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
	reply.C = args.A + args.B
	fllog.Log().Debug("req=", args, "reply=", reply)
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
