package main

import (
	"flag"
	"fmt"

	conf "github.com/forlifeproj/application/gfriends/http_gateway/config"
	"github.com/forlifeproj/application/gfriends/http_gateway/router"
	"github.com/forlifeproj/msf/consul"

	fllog "github.com/forlifeproj/msf/log"
)

var (
	cfg string
)

func init() {
	flag.StringVar(&cfg, "c", "../conf/http_gateway.toml", "config file path, default ../conf/http_gateway.toml")
}

// 测试证书生成工具：https://keymanager.org/#
// 中间件对应的包：github.com/unrolled/secure
func main() {
	flag.Parse()
	// log init
	fllog.Init(cfg)
	// consul init
	consul.Init(cfg)
	fllog.Log().Debug(fmt.Sprintf("consulAddr=%s consulEnvir=%s", consul.GetConsulAddr(), consul.GetConsulEnvironment()))
	// config init
	conf.Init(cfg)
	
	// r := gin.Default()
	// r.Use(httpsHandler()) //https对应的中间件
	// r.GET("/svr/mul", func(c *gin.Context) {
	// 	testMulRpcCall()
	// 	fmt.Println(c.Request.Host)
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"code":   http.StatusOK,
	// 		"result": "test mul rpc succ",
	// 	})
	// })

	// r.GET("/svr/add", func(c *gin.Context) {
	// 	testAddRpcCall()
	// 	fmt.Println(c.Request.Host)
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"code":   http.StatusOK,
	// 		"result": "test add rpc succ",
	// 	})
	// })

	r, err := router.RegisterRouters()
	if err != nil {
		fllog.Log().Error(fmt.Sprintf("register routers failed. err:%+v", err))
		return
	}

	path := "./CA/"                                  //证书的路径
	r.RunTLS(":18080", path+"ca.crt", path+"ca.key") //开启HTTPS服务
}

/*
func httpsHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		secureMiddle := secure.New(secure.Options{
			SSLRedirect: true, //只允许https请求
			//SSLHost:"" //http到https的重定向
			STSSeconds:           1536000, //Strict-Transport-Security header的时效:1年
			STSIncludeSubdomains: true,    //includeSubdomains will be appended to the Strict-Transport-Security header
			STSPreload:           true,    //STS Preload(预加载)
			FrameDeny:            true,    //X-Frame-Options 有三个值:DENY（表示该页面不允许在 frame 中展示，即便是在相同域名的页面中嵌套也不允许）、SAMEORIGIN、ALLOW-FROM uri
			ContentTypeNosniff:   true,    //禁用浏览器的类型猜测行为,防止基于 MIME 类型混淆的攻击
			BrowserXssFilter:     true,    //启用XSS保护,并在检查到XSS攻击时，停止渲染页面
			//IsDevelopment:true,  //开发模式
		})
		err := secureMiddle.Process(context.Writer, context.Request)
		// 如果不安全，终止.
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, "数据不安全")
			return
		}
		// 如果是重定向，终止
		if status := context.Writer.Status(); status > 300 && status < 399 {
			context.Abort()
			return
		}
		context.Next()
	}
}

func testMulRpcCall() {

	args := &demo.Args{
		A: 10,
		B: 20,
	}
	reply := &demo.Reply{}

	callDesc := flcli.CallDesc{
		ServiceName: "gfriends.Arith.Mul",
		Timeout:     time.Second,
	}
	flC := flcli.NewClient(callDesc)
	defer flC.Close()

	flC.DoRequest(context.Background(), args, reply)
	fllog.Log().Debug(fmt.Sprintf("%d * %d = %d", args.A, args.B, reply.C))
}

func testAddRpcCall() {
	args := &demo.Args{
		A: 10,
		B: 20,
	}
	reply := &demo.Reply{}
	callDesc := flcli.CallDesc{
		ServiceName: "gfriends.Demo.Add",
		Timeout:     time.Second,
	}
	flC := flcli.NewClient(callDesc)
	defer flC.Close()

	flC.DoRequest(context.Background(), args, reply)

	fllog.Log().Debug(fmt.Sprintf("%d + %d = %d", args.A, args.B, reply.C))
}
*/
