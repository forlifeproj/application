package account

import (
	"context"
	"fmt"
	"time"

	"github.com/forlifeproj/application/gfriends/login/errno"
	flcli "github.com/forlifeproj/msf/client"
	fllog "github.com/forlifeproj/msf/log"
	"github.com/forlifeproj/protocol/gfriends/json/account"
	"github.com/smallnest/rpcx/protocol"
)

func GetUid(req *account.GetUidReq, rsp *account.GetUidRsp) error {
	callDesc := flcli.CallDesc{
		ServiceName: "gfriends.account.GetUid",
		Timeout:     time.Second,
		CodecType:   protocol.JSON,
	}
	flC := flcli.NewClient(callDesc)
	defer flC.Close()

	if err := flC.DoRequest(context.Background(), req, rsp); err != nil {
		fllog.Log().Error(fmt.Sprintf("getuid failed. err:%+v req:%+v", err, req))
		return errno.ErrGetUidFail
	}
	return nil
}
