package service

import (
	"context"
	"fmt"
	"time"

	"github.com/forlifeproj/application/gfriends/order_manager/conf"
	"github.com/forlifeproj/application/gfriends/order_manager/dao"
	fllog "github.com/forlifeproj/msf/log"
	"github.com/forlifeproj/protocol/gfriends/json/wx_pay"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

var wxPayClient *core.Client

func InitWxPayClient() error {
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath("../conf/apiclient_key.pem")
	if err != nil {
		fllog.Log().Errorf("load merchant private key error, err:%v", err)
		return err
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(conf.GConf.WxPay.MchId, conf.GConf.WxPay.MchCertificateSerialNumber, mchPrivateKey, conf.GConf.WxPay.MchAPIv3Key),
	}
	wxPayClient, err = core.NewClient(ctx, opts...)
	if err != nil {
		fllog.Log().Errorf("new wechat pay client err:%s", err)
		return err
	}

	fllog.Log().Debugf("wxPay client init success")
	return nil
}

func WxPrePay(ctx context.Context, req *wx_pay.WxPrePayReq, rsp *wx_pay.WxPrePayRsp) error {
	groupInfo := GetGroupInfoById(req.GroupId)
	if groupInfo == nil {
		fllog.Log().Errorf("invalid group id:%s", req.GroupId)
		return fmt.Errorf("invalid param")
	} else if req.Amount != int64(groupInfo.Price) && IsProdEnv() {
		fllog.Log().Errorf("invalid price:%d", req.Amount)
		return fmt.Errorf("invalid param")
	}

	svc := jsapi.JsapiApiService{Client: wxPayClient}
	// 得到prepay_id，以及调起支付所需的参数和签名
	resp, result, err := svc.PrepayWithRequestPayment(ctx,
		jsapi.PrepayRequest{
			Appid:       core.String(conf.GConf.WxPay.Appid),
			Mchid:       core.String(conf.GConf.WxPay.MchId),
			Description: core.String("密友来了"),
			OutTradeNo:  core.String(req.OrderId),
			Attach:      core.String("自定义数据说明"),
			NotifyUrl:   core.String(conf.GConf.WxPay.WxCallbackUrl),
			Amount: &jsapi.Amount{
				Total: core.Int64(req.Amount),
			},
			Payer: &jsapi.Payer{
				Openid: core.String(req.OpenId),
			},
		},
	)

	if err != nil {
		fllog.Log().Errorf("wx pre pay failed, err:%v, req:%+v", err, req)
		return err
	}
	fllog.Log().Debugf("resp:%+v, result:%+v, err:%v", resp, result, err)

	orderInfo := dao.OrderStorage{
		OrderId:    req.OrderId,
		Openid:     req.OpenId,
		PayOrderid: *resp.Package,
		GroupID:    req.GroupId,
		Amount:     req.Amount,
		Status:     wx_pay.WxPayStatus_UnFinished,
		CreateTs:   int(time.Now().Unix()),
	}
	if err := dao.AddOrder(&orderInfo); err != nil {
		fllog.Log().Errorf("add order failed, err:%v, orderInfo:%+v", err, orderInfo)
		return err
	}

	rsp.Appid = *resp.Appid
	rsp.NonceStr = *resp.NonceStr
	rsp.Package = *resp.Package
	rsp.SignType = *resp.SignType
	rsp.PaySign = *resp.PaySign
	rsp.TimeStamp = *resp.TimeStamp
	fllog.Log().Debugf("wx pre pay success, req:%+v", req)
	return nil
}
