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
	FLSvrSucc       = 0
	FLSvrFail       = 1000
	FLInvalidReq    = 1001
	HttpRequestFail = 1002
	HttpRspFail     = 1003
	HttpReadRspFail = 1004

	//http_gateway
	FLTransferReqFail = 2001
	FLInvalidUrl      = 2002
	//login
	LoginAuthCodeFail = 3001
	CreateTicketFail  = 3002
	InvalidTicket     = 3003
	AuthTokenFail     = 3004

	//account
	GetUidFail = 4001
)

var (
	// 通用
	ErrOK           = &Errno{Code: FLSvrSucc, Message: ""}
	ErrInternal     = &Errno{Code: FLSvrFail, Message: "内部错误"}
	ErrInvalidReq   = &Errno{Code: FLInvalidReq, Message: "无效请求"}
	ErrHttpRequest  = &Errno{Code: HttpRequestFail, Message: "HTTP请求失败"}
	ErrHttpRspFail  = &Errno{Code: HttpRspFail, Message: "HTTP返回错误"}
	ErrReadHttpFail = &Errno{Code: HttpReadRspFail, Message: "读取HTTP返回失败"}
	// http网关
	ErrTransferReq = &Errno{Code: FLTransferReqFail, Message: "网关转发请求失败"}
	ErrInvalidUrl  = &Errno{Code: FLInvalidUrl, Message: "无效URL"}
	// 登录服务
	ErrLoginAuthFail = &Errno{Code: LoginAuthCodeFail, Message: "第三方授权登录失败"}
	ErrCreateTicket  = &Errno{Code: CreateTicketFail, Message: "创建票据失败"}
	ErrInvalidTicket = &Errno{Code: InvalidTicket, Message: "无效票据"}
	ErrAuthToken     = &Errno{Code: AuthTokenFail, Message: "鉴权失败"}

	// 账号服务
	ErrGetUidFail = &Errno{Code: GetUidFail, Message: "查询UID失败"}
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
