package conf

const (
	WxTradeSuccess         = "SUCCESS"
	HTTP_STATUS_CODE       = "http_status_code"
	HTTP_SERVER_ERROR_CODE = "500"
	EncryptKey             = "6baa188dc4a01374"
	OrderIdPrefix          = "forlife@2024"
)

var GConf = struct {
	Global struct {
		Env string `default:"exp"`
	}
	DB struct {
		User     string
		PassWord string
		StrIp    string
		Port     int
		Database string
		MaxOpen  int
		MaxIdol  int
	}
	WxPay struct {
		Appid                      string `default:"wx3891b8ec2b4e21d3"`                          // 公众号appid
		MchId                      string `default:"1674845257"`                                  // 商户号
		MchCertificateSerialNumber string `default:"3C65C96036169AD7E05EA15DB66ADDF6B8E54E5C"`    // 商户证书序列号
		MchAPIv3Key                string `default:"B12F3E0641E3149469B661BAF693B3D9"`            // 商户APIv3密钥
		WxCallbackUrl              string `default:"http://forlifejj.cn/svr/proc/wxpay_callback"` // 支付结果异步通知地址
	}
	GroupConf struct {
		FriendGroupUrl string
		NewsGroupUrl   string
	}
}{}
