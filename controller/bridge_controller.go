package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sixinyiyu/http-bridge/logger"
	"github.com/sixinyiyu/http-bridge/util"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var (
	customerHeaderName = "headers"
	customerUrlName = "url"
)

// 自定义请求头 支持放置在请求参数['headers]里或者放在请求头里
func IndexHttpHandle(ctx *gin.Context) {
	startReqTime := time.Now()
	headerName := ctx.QueryMap(customerHeaderName)
	remoteURL, err := url.ParseRequestURI(ctx.Query(customerUrlName))
	logger.Sugar.Infof("ssssssssssss %s", remoteURL)
	if nil != err || remoteURL == nil  {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "被代理地址不是一个有效的URL",
		})
		ctx.Abort()
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()

	// 表单数据
	for k, v := range ctx.Request.PostForm {
		logger.Sugar.Infof("PostForm 请求参数 %s = %s", k, v)
		for _, _v := range v {
			req.PostArgs().Add(k, _v)
		}
	}

	// 文件上传 暂时不处理<TODO>
	// query url 解析的请求参数
	for k, v := range  ctx.Request.URL.Query() {
		if k == customerUrlName || k == customerHeaderName {
			continue
		}
		logger.Sugar.Infof("Query URL 参数: %s = %s", k, v)
		if nil != remoteURL {
			remoteURL.Query().Set(k, v[0])
		}
	}

	postBody, err := ioutil.ReadAll(ctx.Request.Body)
	if err == nil {
		logger.Sugar.Infof("ctx.Request.body: %v", string(postBody))
		req.SetBody(postBody)
	}

	//获取透传的请求头 暂时必须都为字符串
	for k, v := range  ctx.Request.Header {
		if k != "Cache-Control" && len(v) > 0 {
			req.Header.Set(k, v[0])
		}
	}
	if len(headerName) > 0 {
		for k, v := range headerName {
			req.Header.Set(k, v)
		}
	}

	/**设置请求参数*/
	req.Header.SetMethod(ctx.Request.Method)
	req.SetRequestURI(remoteURL.String())

	//req.Header.VisitAll(func(key, value []byte) {
	//	logger.Sugar.Infof("=======请求头：%s = %s", string(key), string(value))
	//})

	logger.Sugar.Infof("请求地址: %s, 请求方法: %s", req.URI().String(), ctx.Request.Method)

	// 发送请求
	if err := fasthttp.Do(req, resp); err != nil {
		logger.Sugar.Errorf("请求失败 %v" , err.Error())
	}
	logger.Sugar.Infof("响应结果: %s", util.B2S(resp.Body()))
	ctx.Header("X-Request-Time", time.Since(startReqTime).String())
	resp.Header.VisitAll(func(key, value []byte) {
		ctx.Header(util.B2S(key), util.B2S(value))
	})
	ctx.Data(200, util.B2S(resp.Header.ContentType()), resp.Body())
}
