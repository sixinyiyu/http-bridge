package  main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/sixinyiyu/http-bridge/logger"
	"github.com/valyala/fasthttp"
	"github.com/sixinyiyu/http-bridge/controller"
)

//type Response struct {
//	Code string  `json:"code"`
//	Message string `json:"message"`
//	Data interface{} `json:"data"`
//}

//var (
//	port  = flag.Int("port", 8080, "listen port ")
//)

func main()  {
	//flag.Parse()
	logger.Sugar.Info("------------------ 服务器启动 ------------------")

	// 创建路由
	router := fasthttprouter.New()
	router.HandleOPTIONS = true
	router.GET("/", controller.IndexHttpHandle)
	router.PUT("/", controller.IndexHttpHandle)
	router.POST("/", controller.IndexHttpHandle)
	router.DELETE("/", controller.IndexHttpHandle)
	router.OPTIONS("/", controller.CrosHttpHandle)
	router.NotFound =  controller.NotFoundHttpHandle

	if err := fasthttp.ListenAndServe("0.0.0.0:4455", router.Handler); err != nil {
		logger.Sugar.Errorf("start fastHttp fail", err.Error())
	}
}
