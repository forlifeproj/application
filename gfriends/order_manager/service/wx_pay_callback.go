package service

import (
	"context"

	"encoding/base64"
	"encoding/json"

	"github.com/forlifeproj/application/gfriends/order_manager/conf"
	"github.com/forlifeproj/application/gfriends/order_manager/dao"
	fllog "github.com/forlifeproj/msf/log"
	flsvr "github.com/forlifeproj/msf/server"

	"crypto/aes"
	"crypto/cipher"

	"github.com/forlifeproj/protocol/gfriends/json/wx_pay"
)

type PayCallbackHandler struct {
	Req *wx_pay.WxPayCallbackReq

	BizInfo wx_pay.Transaction
}

func WxPayCallback(ctx context.Context, req *wx_pay.WxPayCallbackReq, rsp *wx_pay.WxPayCallbackRsp) (err error) {
	reqStr, _ := json.Marshal(req)
	fllog.Log().Debugf("req:%s", reqStr)
	resMeta := flsvr.GetResMetaDataMap(ctx)
	defer func() {
		if err != nil {
			resMeta[conf.HTTP_STATUS_CODE] = conf.HTTP_SERVER_ERROR_CODE
			rsp.Code = wx_pay.WxCallbackFail
			rsp.Message = "系统繁忙"
		}
	}()

	// TODO check sign

	h := &PayCallbackHandler{
		Req: req,
	}

	// 解析数据
	if err := h.ParseBizData(); err != nil {
		fllog.Log().Errorf("ParseBizData failed, err:%v, req:%+v", err, req)
		return err
	}

	// 入库
	state := wx_pay.WxPayStatus_Suc
	if h.BizInfo.TradeState != conf.WxTradeSuccess {
		state = wx_pay.WxPayStatus_Fail
	}
	if err := dao.UpdateOrderRecord(h.BizInfo.OutTradeNo, state, h.BizInfo.TransactionID); err != nil {
		fllog.Log().Errorf("update order record failed, err:%v, req:%+v", err, req)
		return err
	}

	fllog.Log().Debugf("success, req:%+v, bizInfo:%+v", req, h.BizInfo)
	return nil
}

func (h *PayCallbackHandler) CheckSign() error {
	//TODO
	return nil
}

func (h *PayCallbackHandler) ParseBizData() error {
	ciphertext, noncestr, associatedData := h.Req.Resource.Ciphertext, h.Req.Resource.Nonce, h.Req.Resource.AssociatedData
	// 解码回调数据
	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		fllog.Log().Errorf("err:%v", err)
		return err
	}
	fllog.Log().Debugf("decodedCiphertext:%s", string(decodedCiphertext))

	// 创建AES cipher.Block，使用给定的密钥
	block, err := aes.NewCipher([]byte(conf.GConf.WxPay.MchAPIv3Key))
	if err != nil {
		fllog.Log().Errorf("创建密钥错误：%v", err)
		return err
	}

	// 创建GCM模式的解密器
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		fllog.Log().Errorf("创建解密器错误：%v", err)
		return err
	}

	// 解密回调数据
	plaintext, err := aesgcm.Open(nil, []byte(noncestr), decodedCiphertext, []byte(associatedData))
	if err != nil {
		fllog.Log().Errorf("解密错误：%v", err)
		return err
	}
	fllog.Log().Debugf("plaintext:%s", string(plaintext))

	// 解析解密后的业务信息
	err = json.Unmarshal(plaintext, &h.BizInfo)
	if err != nil {
		fllog.Log().Errorf("解析业务信息错误：%v, data:%s", err, plaintext)
		return err
	}

	fllog.Log().Debugf("businessData:%+v", h.BizInfo)
	return nil
}
