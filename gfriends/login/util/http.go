package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	// gcfig "git.code.oa.com/going/config"
	"github.com/forlifeproj/application/gfriends/login/errno"
	fllog "github.com/forlifeproj/msf/log"
	// "github.com/smallnest/rpcx/log"
)

// HttpConf http属性配置
type HttpConf struct {
	Address string
	Timeout time.Duration `default:"2s"`
	Path    string
}

// HttpClient 配置+错误信息
type HttpClient struct {
	conf HttpConf
}

// NewClient httpclient
func NewClient(conf *HttpConf) *HttpClient {
	return &HttpClient{conf: *conf}
}

// PostJson 发送post+json请求
func (c *HttpClient) PostJson(req interface{}, rsp interface{},
	headerMap map[string]string, user, password string) (err error) {
	reqData, _ := json.Marshal(req)
	// start := time.Now()
	url := fmt.Sprintf("%s%s", c.conf.Address, c.conf.Path)

	httpStart := time.Now()
	rspData, rspCode, err := HttpPost(url, user, password, c.conf.Timeout, reqData, headerMap)
	if err != nil {
		fllog.Log().Error(fmt.Sprintf("httpPost err: %v, cost:%+v, req: %v, rsp: %s ",
			err, time.Since(httpStart), req, string(rspData)))
		err = errno.ErrHttpRequest
		return
	}
	fllog.Log().Debug(fmt.Sprintf("http post return: req:%+v headermap:%+v rspData:%s rspCode:%d cost:%+v",
		string(reqData), headerMap, string(rspData), rspCode, time.Since(httpStart)))
	if rspCode != http.StatusOK {
		fllog.Log().Error(fmt.Sprintf("httpPost rsp status err: %v, req: %v, rsp: %s", rspCode, req, string(rspData)))
		err = errno.ErrHttpRspFail
		return
	}
	// 回包
	err = json.Unmarshal(rspData, rsp)
	if err != nil {
		fllog.Log().Error(fmt.Sprintf("json.Unmarshal ERR:%+v", err))
		err = errno.ErrReadHttpFail
		return
	}
	return
}

func (c *HttpClient) HttpGet(params map[string]string, timeout time.Duration, rsp interface{}) error {
	// start := time.Now()
	fullUrl := fmt.Sprintf("%s%s", c.conf.Address, c.conf.Path)
	getParams := url.Values{}
	for k, v := range params {
		getParams.Add(k, v)
	}
	fullUrl = fmt.Sprintf("%s?%s", fullUrl, getParams.Encode())

	fllog.Log().Debug(fmt.Sprintf("fullUrl:%s", fullUrl))

	client := &http.Client{
		Timeout: timeout,
	}

	httpReq, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		fllog.Log().Error(fmt.Sprintf("NewRequest err:%+v fullUrl:%s", err, fullUrl))
		return errno.ErrHttpRequest
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		fllog.Log().Error(fmt.Sprintf("DoRequest err:%+v fullUrl:%s", err, fullUrl))
		return errno.ErrHttpRequest
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fllog.Error("ReadHttpRsp err:%+v fullUrl:%s", err, fullUrl)
		return errno.ErrReadHttpFail
	}

	if resp.StatusCode != http.StatusOK {
		fllog.Log().Error(fmt.Sprintf("ERR resp.StatusCode:%d", resp.StatusCode))
		return errno.ErrHttpRspFail
	}
	err = json.Unmarshal(body, rsp)
	fllog.Log().Debug(fmt.Sprintf("fullUrl:%s rsp:%+v", fullUrl, rsp))
	return nil
}

func HttpPost(url, user, password string, timeout time.Duration,
	reqBody []byte, headerMap map[string]string) (rspBody []byte, httpCode int, err error) {

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		err = fmt.Errorf("new http request fail, err:%v, url: %s", err, url)
		return
	}

	for k, v := range headerMap {
		httpReq.Header.Set(k, v)
	}

	if len(user) > 0 && len(password) > 0 {
		httpReq.SetBasicAuth(user, password)
	}

	// 关闭keep alive
	// httpReq.Close = true

	client := &http.Client{Timeout: timeout}
	httpRsp, err := client.Do(httpReq)
	if err != nil {
		err = fmt.Errorf("do http req fail, err:%v, rsp:%v", err, httpRsp)
		return
	}
	defer httpRsp.Body.Close()

	httpCode = httpRsp.StatusCode
	rspBody, err = ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		err = fmt.Errorf("read rspbody fail, err:%v, body: %v", err, string(rspBody))
	}

	return
}
