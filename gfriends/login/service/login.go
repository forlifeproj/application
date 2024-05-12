package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	ac "github.com/forlifeproj/application/gfriends/login/account"
	"github.com/forlifeproj/application/gfriends/login/dao"
	"github.com/forlifeproj/application/gfriends/login/errno"
	"github.com/forlifeproj/application/gfriends/login/plugin"
	"github.com/forlifeproj/application/gfriends/login/ticket"
	fllog "github.com/forlifeproj/msf/log"
	"github.com/forlifeproj/protocol/gfriends/json/account"
	"github.com/forlifeproj/protocol/gfriends/json/login"
)

func LoginIn(ctx context.Context, req *login.LoginReq, rsp *login.LoginRsp) (err error) {
	fllog.Log().Debug(fmt.Sprintf("req=%+v rsp+%+v", req, rsp))
	rsp.Code = 0
	rsp.ErrMsg = "成功"
	defer func() {
		if err != nil {
			code, msg := errno.DecodeErr(err)
			rsp.Code = code
			rsp.ErrMsg = msg
		}
		fllog.Log().Error(fmt.Sprintf("LoginIn req:%+v rsp:%+v", req, rsp))
	}()

	// step 1 plugin auth
	authRsp, err := pluginAuth(req)
	if err != nil {
		return err
	}

	// step 2 get uid
	registerUid, err := getLoginUid(req, authRsp)
	if err != nil {
		return err
	}

	// step 3 create ticket
	createTicket := createTicket(registerUid, req, authRsp)
	if len(createTicket) == 0 {
		fllog.Log().Error("create ticket failed.")
		return errno.ErrCreateTicket
	}
	rsp.Ticket = createTicket
	rsp.RegisterUid = registerUid
	rsp.Openid = authRsp.StrOpenid

	// step 4 async save ticket
	go saveTicket(registerUid, createTicket, req, authRsp)

	return nil
}

func pluginAuth(req *login.LoginReq) (*plugin.AuthRsp, error) {
	loginHandler, err := plugin.GetLoginHandler(req.LoginType)
	if err != nil {
		fllog.Log().Error(fmt.Sprintf("invalid req:%+v", req))
		return nil, err
	}

	type AuthRsp struct {
		ErrCode           int    `json:"errcode"`
		ErrMsg            string `json:"errmsg"`
		StrOpenid         string `json:"openid"`
		StrUnionid        string `json:"unionid"`
		StrRefreshToken   string `json:"refresh_token"`
		StrAccessToken    string `json:"access_token"`
		AccessTokenExpire int64  `json:"expires_in"`
	}
	// 自测
	if req.Token == "1234" {
		authRsp := &plugin.AuthRsp{
			StrOpenid:         "1111",
			StrUnionid:        "2222",
			StrRefreshToken:   "refresh_token",
			StrAccessToken:    "access_token",
			AccessTokenExpire: 3600,
		}
		return authRsp, nil
	}
	authRsp, err := loginHandler.AuthorziationCode(req.Token)
	if err != nil || authRsp.ErrCode != 0 {
		fllog.Log().Error(fmt.Sprintf("login auth code failed. err:%+v req:%+v authRsp:%+v",
			err, req, authRsp))
		return authRsp, errno.ErrLoginAuthFail
	}
	return authRsp, nil
}

func getLoginUid(req *login.LoginReq, authRsp *plugin.AuthRsp) (int64, error) {
	getUidReq := &account.GetUidReq{
		Appid:        req.Appid,
		OpenType:     req.LoginType,
		OpenId:       authRsp.StrOpenid,
		UnionId:      authRsp.StrUnionid,
		AutoRegister: 1,
	}
	getUidRsp := &account.GetUidRsp{}
	err := ac.GetUid(getUidReq, getUidRsp)
	if err != nil {
		fllog.Log().Error(fmt.Sprintf("get uid failed. err:%+v req:%+v", err, req))
		return 0, err
	}
	if getUidRsp.Code != 0 {
		fllog.Log().Error(fmt.Sprintf("get uid return failed. req:%+v rsp:%+v", getUidReq, getUidRsp))
		return 0, errno.ErrGetUidFail
	}
	return getUidRsp.RegisterUid, nil
}

func createTicket(uid int64, req *login.LoginReq, authRsp *plugin.AuthRsp) string {
	createTicket := ticket.Ticket{}
	createTicket.SetAppid(req.Appid)
	createTicket.SetVersion(ticket.Version1_0)
	createTicket.SetLoginType(req.LoginType)
	createTicket.SetOpenid(authRsp.StrOpenid)
	rand.Seed(time.Now().UnixNano())
	createTicket.SetRandom(rand.Intn(1000000000))
	nowTime := time.Now().Unix()
	createTicket.SetCreateTime(int(nowTime))
	return createTicket.CreateTicket()
}

func saveTicket(uid int64, ticket string, req *login.LoginReq, authRsp *plugin.AuthRsp) {
	nowTime := time.Now().Unix()
	loginToken := dao.LoginToken{
		Token:              ticket,
		Uid:                uid,
		OpenId:             authRsp.StrOpenid,
		UnionId:            authRsp.StrUnionid,
		LoginType:          req.LoginType,
		AccessToken:        authRsp.StrAccessToken,
		RefreshToken:       authRsp.StrRefreshToken,
		AccessTokenExpire:  authRsp.AccessTokenExpire,
		RefreshTokenExpire: 0,
		Status:             1,
		CreateTime:         nowTime,
		ExpireTime:         nowTime + 3600*24*30,
	}
	if err := dao.AddLoginToken(&loginToken); err != nil {
		fllog.Log().Error("insert into logintoken failed. err=", err, " loginToken=", loginToken)
	}
}
