package plugin

import (
	"fmt"

	"github.com/forlifeproj/application/gfriends/login/errno"
	fllog "github.com/forlifeproj/msf/log"
	"github.com/forlifeproj/protocol/gfriends/json/login"
)

type AuthRsp struct {
	ErrCode           int    `json:"errcode"`
	ErrMsg            string `json:"errmsg"`
	StrOpenid         string `json:"openid"`
	StrUnionid        string `json:"unionid"`
	StrRefreshToken   string `json:"refresh_token"`
	StrAccessToken    string `json:"access_token"`
	AccessTokenExpire int64  `json:"expires_in"`
}

type RefreshRsp struct {
	ErrCode           int    `json:"errcode"`
	ErrMsg            string `json:"errmsg"`
	StrOpenid         string `json:"openid"`
	StrUnionid        string `json:"unionid"`
	StrAccessToken    string `json:"access_token"`
	AccessTokenExpire int64  `json:"expires_in"`
}

type LoginHandler interface {
	AuthorziationCode(strCode string) (*AuthRsp, error)
	RefreshAccessToken(refreshToken string) (*RefreshRsp, error)
}

func GetLoginHandler(loginType int) (LoginHandler, error) {
	var handler LoginHandler
	switch loginType {
	case login.WEIXIN_LOGIN_TYPE:
		handler = NewWeiXin()
	default:
		fllog.Log().Error(fmt.Sprintf("invalid logintype:%d", loginType))
		return handler, errno.ErrInvalidReq
	}
	return handler, nil
}
