package plugin

import (
	"fmt"

	conf "github.com/forlifeproj/application/gfriends/login/config"
	"github.com/forlifeproj/application/gfriends/login/util"
	fllog "github.com/forlifeproj/msf/log"
)

type WeiXinLogin struct {
}

func NewWeiXin() *WeiXinLogin {
	weiXinLogin := &WeiXinLogin{}
	return weiXinLogin
}

func (w *WeiXinLogin) AuthorziationCode(strCode string) (*AuthRsp, error) {
	paramsMap := make(map[string]string)
	paramsMap["appid"] = conf.GConf.WeiXinCfg.StrAppid
	paramsMap["secret"] = conf.GConf.WeiXinCfg.StrSecret
	paramsMap["code"] = strCode
	paramsMap["grant_type"] = "authorization_code"
	paramsMap["simple_get_token"] = "1"

	rsp := &AuthRsp{}
	httpConf := &util.HttpConf{
		Address: conf.GConf.WeiXinCfg.Address,
		Path:    conf.GConf.WeiXinCfg.AuthPath,
	}
	httpCli := util.NewClient(httpConf)

	if err := httpCli.HttpGet(paramsMap, conf.GConf.WeiXinCfg.Timeout, rsp); err != nil {
		fllog.Log().Error(fmt.Sprintf("weixin authozization code failed. err:%+v paramsMap:%+v", err, paramsMap))
		return rsp, err
	}

	return rsp, nil
}

func (w *WeiXinLogin) RefreshAccessToken(refreshToken string) (*RefreshRsp, error) {
	paramsMap := make(map[string]string)
	paramsMap["appid"] = conf.GConf.WeiXinCfg.StrAppid
	paramsMap["refresh_token"] = refreshToken
	paramsMap["grant_type"] = "refresh_token"
	paramsMap["simple_get_token"] = "1"

	rsp := &RefreshRsp{}
	httpConf := &util.HttpConf{
		Address: conf.GConf.WeiXinCfg.Address,
		Path:    conf.GConf.WeiXinCfg.RefreshPath,
	}
	httpCli := util.NewClient(httpConf)

	if err := httpCli.HttpGet(paramsMap, conf.GConf.WeiXinCfg.Timeout, rsp); err != nil {
		fllog.Log().Error(fmt.Sprintf("weixin refresh accesstoken code failed. err:%+v paramsMap:%+v", err, paramsMap))
		return rsp, err
	}

	return rsp, nil
}
