package controller

import (
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sixinyiyu/http-bridge/logger"
	"github.com/sixinyiyu/http-bridge/util"
	"github.com/valyala/fasthttp"
	"time"
)

// 自定义请求头 支持放置在请求参数['headers]里或者放在请求头里
func IndexHttpHandle(ctx *fasthttp.RequestCtx) {
	startReqTime := time.Now()
	queryArgs := ctx.QueryArgs()
	redirectUrl := string(queryArgs.Peek("url"))
	method := util.B2S(ctx.Method())
	logger.Sugar.Infof("请求地址: %s, 请求方法: %s", redirectUrl, method)

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	// 获取透传的请求头 暂时必须都为字符串
	customerHeaders := util.B2S(queryArgs.Peek("headers"))
	if customerHeaders != "" {
		logger.Sugar.Infof("自定义请求头", customerHeaders)
		var headerMap map[string] string
		if err := ffjson.Unmarshal(queryArgs.Peek("headers"), &headerMap); err != nil {
			logger.Sugar.Error(err.Error())
		}
	}

	ctx.Request.Header.VisitAll(func(key, value []byte) {
		logger.Sugar.Infof("key: %s, value: %s", util.B2S(key), util.B2S(value))
		//req.Header.SetBytesKV(key, value)
	})

	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()
        if method != fasthttp.MethodGet {
		if requestBody := ctx.PostBody(); requestBody != nil {
			req.SetBody(requestBody)
		}
	}

	/**设置请求参数*/
	req.Header.SetContentType(util.B2S(ctx.Request.Header.Peek("Content-Type")))
	req.Header.SetMethod(method)
	req.SetRequestURI(redirectUrl)


	// 发送请求
	if err := fasthttp.Do(req, resp); err != nil {
		logger.Sugar.Errorf("请求失败 %s" , err.Error())
	}

	respText := util.B2S(resp.Body())

	logger.Sugar.Infof("响应结果: %s", respText)

	costTime := time.Since(startReqTime)
	ctx.Response.Header.Set("X-Request-Time", costTime.String())
	resp.Header.VisitAll(func(key, value []byte) {
		ctx.Response.Header.SetBytesKV(key, value)
	})

	_, _ = ctx.WriteString(respText)
}

// 跨域
func CrosHttpHandle(ctx *fasthttp.RequestCtx) {
	logger.Sugar.Infof("请求IP", ctx.RemoteAddr())
	origin := util.B2S(ctx.Request.Header.Peek("Origin"))
	if (origin != "") && (origin != "null") {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", origin)
	}else {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	}
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	ctx.Response.Header.Set("Access-Control-Allow-Headers",
		"Accept, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Token, token")
	ctx.Response.Header.Set("Content-Type","application/json; charset=utf-8")
	ctx.SetStatusCode(fasthttp.StatusOK)
}

// 404
func NotFoundHttpHandle(ctx * fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	_, _ = ctx.WriteString("{\"code\": \"500\", \"message\": \"请求路径错误\"}")
}

