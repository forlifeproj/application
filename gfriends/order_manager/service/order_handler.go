package service

import (
	"context"
	"fmt"
	"time"

	fllog "github.com/forlifeproj/msf/log"

	"github.com/forlifeproj/application/gfriends/order_manager/dao"
	"github.com/forlifeproj/protocol/gfriends/json/order_manager"
)

func GenerateOrderId(ctx context.Context, req *order_manager.GenerateOrderIdReq, rsp *order_manager.GenerateOrderIdRsp) error {
	rsp.OrderId = fmt.Sprintf("%d", time.Now().UnixNano())
	fllog.Log().Debug("req=", req, "rsp=", rsp)
	return nil
}

func QueryOrderStatus(ctx context.Context, req *order_manager.QueryOrderStatusReq, rsp *order_manager.QueryOrderStatusRsp) error {
	rsp.MaxPollTimes = 60
	rsp.PollInterval = 1000
	rsp.NeedPoll = 1
	defer func() {
		if rsp.Status == 1 {
			rsp.NeedPoll = 0
			// TODO 给bizInfo赋值
		}
	}()
	order, err := dao.GetOrderRecord(req.OrderId)
	if err != nil {
		fllog.Log().Errorf("err:%v, req:%+v", err, req)
		return fmt.Errorf("query db error")
	}
	if order == nil {
		return nil
	}
	rsp.Status = 1
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
