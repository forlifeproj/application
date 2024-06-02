package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	math_rand "math/rand"
	"time"

	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	fllog "github.com/forlifeproj/msf/log"
	"gorm.io/gorm"

	"github.com/forlifeproj/application/gfriends/order_manager/conf"
	"github.com/forlifeproj/application/gfriends/order_manager/dao"
	"github.com/forlifeproj/protocol/gfriends/json/order_manager"
	"github.com/forlifeproj/protocol/gfriends/json/wx_pay"
)

func GenerateOrderId(ctx context.Context, req *order_manager.GenerateOrderIdReq, rsp *order_manager.GenerateOrderIdRsp) error {
	rsp.OrderId = getEncryptOrderId()
	fllog.Log().Debug("req=", req, "rsp=", rsp)
	//debug
	//checkOrderId(rsp.OrderId)
	return nil
}

func QueryOrderStatus(ctx context.Context, req *order_manager.QueryOrderStatusReq, rsp *order_manager.QueryOrderStatusRsp) error {
	rsp.MaxPollTimes = 60
	rsp.PollInterval = 1000
	rsp.NeedPoll = 1

	if !checkOrderId(req.OrderId) {
		fllog.Log().Errorf("orderId invalid, req:%+v", req)
		rsp.NeedPoll = 0
		return fmt.Errorf("非法的订单id")
	}

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

	groupInfo := GetGroupInfoById(order.GroupID)
	if groupInfo == nil {
		fllog.Log().Errorf("invalid group id:%s", order.GroupID)
		return nil
	}
	bizInfo, _ := json.Marshal(groupInfo)
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

func getEncryptOrderId() string {
	orderId := fmt.Sprintf("%s_%d%d", conf.OrderIdPrefix, time.Now().Unix(), math_rand.Intn(10000))
	ciphertext, err := AesEncrypt([]byte(conf.EncryptKey), []byte(orderId))
	if err != nil {
		fllog.Log().Errorf("err:%v", err)
	}
	base64str := base64.StdEncoding.EncodeToString(ciphertext)
	fllog.Log().Debugf("orderId:%s, encryptId:%s", orderId, base64str)
	return base64str
}

func checkOrderId(orderId string) bool {
	base64Decode, err := base64.StdEncoding.DecodeString(orderId)
	if err != nil {
		fllog.Log().Errorf("err:%v", err)
		return false
	}

	decryptedText, err := AesDecrypt([]byte(conf.EncryptKey), base64Decode)
	if err != nil {
		fllog.Log().Errorf("err:%v", err)
		return false
	}
	if !strings.HasPrefix(string(decryptedText), conf.OrderIdPrefix) {
		fllog.Log().Errorf("decrypt order id [%s] invalid", string(decryptedText))
		return false
	}
	fllog.Log().Debugf("orderId %s check pass", orderId)
	return true
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AES加密,CBC
func AesEncrypt(key, origData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// AES解密
func AesDecrypt(key, crypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}
