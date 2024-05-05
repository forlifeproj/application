package service

import (
	"context"
	"fmt"

	"github.com/forlifeproj/application/gfriends/account/dao"
	"github.com/forlifeproj/application/gfriends/account/errno"
	fllog "github.com/forlifeproj/msf/log"
	"github.com/forlifeproj/protocol/gfriends/json/account"
)

func GetUid(ctx context.Context, req *account.GetUidReq, rsp *account.GetUidRsp) error {
	fllog.Log().Debug("getuid recv req=", req)
	registerUid, err := getUidByOpenidUnionid(req.OpenId, req.UnionId, req.OpenType)
	if err != nil {
		return err
	}

	if registerUid > 0 {
		rsp.RegisterUid = registerUid
		fllog.Log().Debug(fmt.Sprintf("has registered uid req:%+v rsp:%+v", req, rsp))
		return nil
	}

	registerUid, err = dao.RegisterUid(req.OpenId, req.UnionId, req.OpenType)
	if err != nil {
		fllog.Log().Debug(fmt.Sprintf("register uid failed. err:%+v req:%+v", err, req))
		return err
	}
	rsp.RegisterUid = registerUid

	fllog.Log().Debug(fmt.Sprintf("getuid succ req:%+v rsp:%+v", req, rsp))
	return nil
}

func getUidByOpenidUnionid(strOpenId, strUnionId string, openType int) (int64, error) {
	if len(strOpenId) > 0 {
		openID, err := dao.GetOpenID(strOpenId, openType)
		if err != nil {
			fllog.Log().Error(fmt.Sprintf("get OpenID from db failed. err:%+v openId:%s openType:%d",
				err, strOpenId, openType))
			return 0, errno.ErrDBQueryail
		}
		if openID.Uid > 0 {
			return openID.Uid, nil
		}
	}

	if len(strUnionId) > 0 {
		unionID, err := dao.GetUnionID(strUnionId, openType)
		if err != nil {
			fllog.Log().Error(fmt.Sprintf("get UnionID from db failed. err:%+v unionId:%s openType:%d",
				err, strUnionId, openType))
			return 0, errno.ErrDBQueryail
		}
		if unionID.Uid > 0 {
			return unionID.Uid, nil
		}
	}

	return 0, nil
}
