package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/forlifeproj/msf/consul"
	"github.com/gin-gonic/gin"

	conf "github.com/forlifeproj/application/gfriends/http_gateway/config"
	"github.com/forlifeproj/application/gfriends/http_gateway/errno"
	"github.com/forlifeproj/application/gfriends/http_gateway/middle"
	"github.com/forlifeproj/application/gfriends/http_gateway/trans"

	fllog "github.com/forlifeproj/msf/log"
	"github.com/gin-gonic/gin/binding"
)

func RegisterRouters() (*gin.Engine, error) {
	svr := gin.New()
	svr.Use(gin.Recovery())
	svr.Use(middle.HttpsHandler())
	// gin.SetMode(gin.ReleaseMode)

	groupUrl := fmt.Sprintf("/svr/%s", consul.GetConsulEnvironment())
	routerGroup := svr.Group(groupUrl)
	for _, routerCfg := range conf.Conf.RouterTableMap {
		routerGroup.POST(routerCfg.Url, RouteHandler)
	}

	return svr, nil
}

type HttpRsp struct {
	Code    int
	Message string
	RspBody string
}

func RouteHandler(ctx *gin.Context) {

	var err error
	fullPath := ctx.FullPath()

	// tArgs := &demo.Args{
	// 	A: 100,
	// 	B: 300,
	// }
	// argBytes, _ := json.Marshal(tArgs)

	// var ttJ interface{}
	// terr := json.Unmarshal(argBytes, &ttJ)
	// if terr != nil {
	// 	fllog.Log().Debug(fmt.Sprintf("json unmarshal failed err:%+v", terr))
	// } else {
	// 	fllog.Log().Debug(fmt.Sprintf("json unmarshal succ ttJ:%+v string(argBytes)=%s", ttJ, string(argBytes)))
	// }

	var args interface{}
	var reply interface{}

	defer func() {
		rsp := HttpRsp{
			Code:    errno.FLSvrSucc,
			Message: "SUCC",
		}
		if err != nil {
			rsp.Code, rsp.Message = errno.DecodeErr(err)
		} else {
			rspBytes, _ := json.Marshal(reply)
			rsp.RspBody = string(rspBytes)
		}
		// 回包
		ctx.JSON(http.StatusOK, rsp)
		fllog.Log().Debug(fmt.Sprintf("req:%+v rsp:%+v", args, reply))

	}()

	if err := ctx.ShouldBindBodyWith(&args, binding.JSON); err != nil {
		body, _ := ctx.Get(gin.BodyBytesKey)
		body = string(body.([]byte))
		fllog.Log().Error(fmt.Sprintf("invalid parameter.req: %v, err: %v, %+v", body, err, args))
	}
	fllog.Log().Debug(fmt.Sprintf("args=%+v", args))
	// body, _ := ctx.Get(gin.BodyBytesKey)
	// bodyBytes := body.([]byte)
	// err = json.Unmarshal(bodyBytes, &args)
	// if err != nil {
	// 	fllog.Log().Debug(fmt.Sprintf("json unmarshal failed err:%+v, bodyBytes:%+v", err, string(bodyBytes)))
	// 	err = errno.ErrInvalidReq
	// 	return
	// } else {
	// 	fllog.Log().Debug(fmt.Sprintf("json unmarshal succ req:%+v", args))
	// }

	err = trans.TransferReq(fullPath, args, &reply)

}
