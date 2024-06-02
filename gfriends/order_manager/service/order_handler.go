package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	fllog "github.com/forlifeproj/msf/log"
	"gorm.io/gorm"

	"github.com/forlifeproj/application/gfriends/order_manager/dao"
	"github.com/forlifeproj/protocol/gfriends/json/order_manager"
	"github.com/forlifeproj/protocol/gfriends/json/wx_pay"
)

func GenerateOrderId(ctx context.Context, req *order_manager.GenerateOrderIdReq, rsp *order_manager.GenerateOrderIdRsp) error {
	rsp.OrderId = fmt.Sprintf("%d%d", time.Now().UnixNano(), rand.Intn(10000))
	fllog.Log().Debug("req=", req, "rsp=", rsp)
	return nil
}

func QueryOrderStatus(ctx context.Context, req *order_manager.QueryOrderStatusReq, rsp *order_manager.QueryOrderStatusRsp) error {
	rsp.MaxPollTimes = 60
	rsp.PollInterval = 1000
	rsp.NeedPoll = 1

	order, err := dao.GetOrderRecord(req.OrderId)
	if err != nil {
		fllog.Log().Errorf("err:%v, req:%+v", err, req)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rsp.NeedPoll = 0
			return nil
		}
		return fmt.Errorf("query db error")
	}
	fllog.Log().Debugf("order:%+v", order)

	// 等待支付回调中
	if order.Status == wx_pay.WxPayStatus_UnFinished {
		fllog.Log().Debugf("order not finished, orderId:%s", req.OrderId)
		rsp.NeedPoll = 1
		return nil
	}

	// 支付失败
	if order.Status == wx_pay.WxPayStatus_Fail {
		fllog.Log().Debugf("order failed, orderId:%s", req.OrderId)
		rsp.NeedPoll = 0
		return nil
	}

	rsp.Status = wx_pay.WxPayStatus_Suc
	rsp.NeedPoll = 0
	groupInfo := order_manager.GroupInfo{
		GroupId:     "1",
		GroupName:   "密友交友群",
		GroupQRCode: "xxx", //TODO
	}
	bizInfo, _ := json.Marshal(&groupInfo)
	rsp.BizInfo = string(bizInfo)
	fllog.Log().Debugf("req:%+v, rsp:%+v", req, rsp)
	return nil
}

func GetOrderList(ctx context.Context, req *order_manager.GetOrderListReq, rsp *order_manager.GetOrderListRsp) error {
	orders, err := dao.GetOrderList(req.Uid, req.StrPassBack, req.PageSize)
	if err != nil {
		fllog.Log().Errorf("err:%v, req:%+v", err, req)
		return fmt.Errorf("query db error")
	}

	if req.Uid == 0 {
		return fmt.Errorf("invalid param")
	}

	if req.PageSize >= 20 || req.PageSize <= 0 {
		req.PageSize = 20
	}

	for _, item := range orders {
		rsp.OrderList = append(rsp.OrderList, order_manager.OrderInfo{
			OrderId: item.OrderId,
			//TODO add more fields
		})
	}

	fllog.Log().Debugf("req:%+v, rsp:%+v", req, rsp)
	return nil
}
