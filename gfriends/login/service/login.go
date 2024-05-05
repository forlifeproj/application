package service

import (
	"context"
	"fmt"

	"github.com/forlifeproj/application/gfriends/login/errno"
	"github.com/forlifeproj/application/gfriends/login/plugin"
	fllog "github.com/forlifeproj/msf/log"
	"github.com/forlifeproj/protocol/gfriends/json/login"
)

func LoginIn(ctx context.Context, req *login.LoginReq, rsp *login.LoginRsp) error {
	loginHandler, err := getLoginHandler(req.LoginType)
	if err != nil {
		fllog.Log().Error(fmt.Sprintf("invalid req:%+v", req))
		return err
	}

	authRsp, err := loginHandler.AuthorziationCode(req.Token)
	if err != nil || authRsp.ErrCode != 0 {
		fllog.Log().Error(fmt.Sprintf("login auth code failed. err:%+v req:%+v authRsp:%+v",
			err, req, authRsp))
		return errno.ErrLoginAuthFail
	}

	// to-do getUid

	return nil
}

func getLoginHandler(loginType int) (plugin.LoginHandler, error) {
	var handler plugin.LoginHandler
	switch loginType {
	case login.WEIXIN_LOGIN_TYPE:
		handler = plugin.NewWeiXin()
	default:
		fllog.Log().Error(fmt.Sprintf("invalid logintype:%d", loginType))
		return handler, errno.ErrInvalidReq
	}
	return handler, nil
}
