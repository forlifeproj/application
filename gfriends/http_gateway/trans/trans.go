package trans

import (
	"context"
	"fmt"
	"strings"
	"time"

	conf "github.com/forlifeproj/application/gfriends/http_gateway/config"
	"github.com/forlifeproj/application/gfriends/http_gateway/errno"
	flcli "github.com/forlifeproj/msf/client"
	"github.com/forlifeproj/msf/consul"
	fllog "github.com/forlifeproj/msf/log"
	"github.com/smallnest/rpcx/protocol"
)

func TransferReq(url string, args interface{}, reply interface{}) error {
	serviceName := GetServiceNameByUrl(url)
	if len(serviceName) == 0 {
		fllog.Log().Error(fmt.Sprintf("unsupport url;%s transfer", url))
		return errno.ErrInvalidUrl
	}
	callDesc := flcli.CallDesc{
		ServiceName: serviceName,
		Timeout:     time.Second,
		CodecType:   protocol.JSON,
	}
	flC := flcli.NewClient(callDesc)
	defer flC.Close()

	if err := flC.DoRequest(context.Background(), &args, reply); err != nil {
		fllog.Log().Error(fmt.Sprintf("transfer req failed. err:%+v req:%+v", err, args))
		return errno.ErrTransferReq
	}
	fllog.Log().Debug(fmt.Sprintf("args:%+v reply:%+v", args, reply))
	return nil
}

func GetServiceNameByUrl(url string) string {
	basePath := fmt.Sprintf("/svr/%s", consul.GetConsulEnvironment())
	urlKey := strings.TrimPrefix(url, basePath)
	routerCfg, ok := conf.Conf.RouterTableMap[urlKey]
	if !ok {
		return ""
	}
	return routerCfg.ServiceName
}
