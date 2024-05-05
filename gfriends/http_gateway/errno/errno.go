package errno

import (
	"fmt"

	fllog "github.com/forlifeproj/msf/log"
)

// Errno 错误码定义
type Errno struct {
	Code    int
	Message string
}

// Error 错误接口
func (err *Errno) Error() string {
	return fmt.Sprintf("Errno - code: %d, message: %s", err.Code, err.Message)
}

const (
	FLSvrSucc    = 0
	FLSvrFail    = 1000
	FLInvalidReq = 1001
	//http_gateway
	FLTransferReqFail = 2001
	FLInvalidUrl      = 2002
)

var (
	ErrOK          = &Errno{Code: FLSvrSucc, Message: ""}
	ErrInternal    = &Errno{Code: FLSvrFail, Message: "内部错误"}
	ErrInvalidReq  = &Errno{Code: FLInvalidReq, Message: "无效请求"}
	ErrTransferReq = &Errno{Code: FLTransferReqFail, Message: "网关转发请求失败"}
	ErrInvalidUrl  = &Errno{Code: FLInvalidUrl, Message: "无效URL"}
)

func DecodeErr(err error) (int, string) {
	if nil == err {
		return 0, ""
	}

	switch typed := err.(type) {
	case *Errno:
		return typed.Code, typed.Message
	default:
		fllog.Log().Error("unknown err type=", typed)
	}

	return ErrInternal.Code, err.Error()
}
