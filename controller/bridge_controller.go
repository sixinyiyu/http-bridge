package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sixinyiyu/http-bridge/logger"
	"github.com/sixinyiyu/http-bridge/util"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"time"
)

var (
	customerHeaderName = "headers"
	customerUrlName = "url"
)

// 自定义请求头 支持放置在请求参数['headers]里或者放在请求头里
func IndexHttpHandle(ctx *gin.Context) {
	startReqTime := time.Now()
	targetUrl := ctx.Query(customerUrlName)
	logger.Sugar.Infof("请求Content-type: %v", ctx.ContentType())

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()

	/**设置请求参数*/
	req.Header.SetContentType(ctx.ContentType())
	req.Header.SetMethod(ctx.Request.Method)
	req.SetRequestURI(targetUrl)

	// 表单数据
	for k, v := range ctx.Request.PostForm {
		logger.Sugar.Infof("PostForm 请求参数 %s = %s", k, v)
		for _, _v := range v {
			req.PostArgs().Add(k, _v)
		}
	}

	// 文件上传 暂时不处理

	// query url 解析的请求参数
	for k, v := range  ctx.Request.URL.Query() {
		if k == customerUrlName {
			continue
		}
		logger.Sugar.Infof("Query URL 参数: %s = %s", k, v)
		if k == customerHeaderName {
			if len(v) > 0 {
				var headerMap map[string] string
				if err := ffjson.Unmarshal([]byte(v[0]), &headerMap); err != nil {
					logger.Sugar.Error(err.Error())
				}
				for k, v := range  headerMap {
					logger.Sugar.Infof("自定义请求头; %v=%v", k, v)
					req.Header.Set(k, v)
				}
			}
		} else {
			for _, _v := range  v {
				req.URI().QueryArgs().Add(k, _v)
			}
		}
	}

	logger.Sugar.Infof("请求地址: %s, 请求方法: %s", req.URI().String(), ctx.Request.Method)

	postBody, err := ioutil.ReadAll(ctx.Request.Body)
	if err == nil {
		logger.Sugar.Infof("ctx.Request.body: %v", string(postBody))
		req.SetBody(postBody)
	}

	//获取透传的请求头 暂时必须都为字符串
	for k, v := range  ctx.Request.Header {
		logger.Sugar.Infof("请求头：%s = %s", k,v)
		for _, _v := range v {
			req.Header.Add(k, _v)
		}
	}

	// 发送请求
	logger.Sugar.Infof("远程地址: %s", req.URI().String())
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
