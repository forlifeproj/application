package service

import (
	"context"
	"fmt"

	"github.com/forlifeproj/application/gfriends/login/dao"
	"github.com/forlifeproj/application/gfriends/login/errno"
	"github.com/forlifeproj/application/gfriends/login/ticket"
	fllog "github.com/forlifeproj/msf/log"
	"github.com/forlifeproj/protocol/gfriends/json/login"
	"gorm.io/gorm"
)

func LoginAuth(ctx context.Context, req *login.LoginAuthReq, rsp *login.LoginAuthRsp) (err error) {
	fllog.Log().Debug(fmt.Sprintf("login auth recvreq=%+v", req))
	defer func() {
		if err != nil {
			code, msg := errno.DecodeErr(err)
			rsp.Code = code
			rsp.ErrMsg = msg
		}
		fllog.Log().Error(fmt.Sprintf("LoginIn req:%+v rsp:%+v", req, rsp))
	}()
	authTicket := ticket.Ticket{}
	authTicket.SetTicket(req.Token)
	authTicket.SetOpenid(req.Openid)
	authTicket.SetAppid(req.Appid)
	authTicket.SetLoginType(req.LoginType)

	if !authTicket.IsValidTicket() {
		fllog.Log().Error(fmt.Sprintf("invalid token :ticket:%+v req:%+v ", authTicket, req))
		err = errno.ErrInvalidTicket
		return
	}

	loginToken, err := dao.GetLoginToken(req.Token)
	if err != nil || err == gorm.ErrRecordNotFound || loginToken == nil {
		fllog.Log().Error(fmt.Sprintf("get login token failed. err:%+v token:%s", err, req.Token))
		err = errno.ErrAuthToken
		return
	}

	if loginToken.Uid != req.RegisterUid {
		fllog.Log().Error(fmt.Sprintf("invalid uid loginToken:%+v req:%+v ", loginToken, req))
		err = errno.ErrInvalidTicket
		return
	}
	return nil
}
