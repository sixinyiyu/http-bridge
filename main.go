package  main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sixinyiyu/http-bridge/controller"
	"github.com/sixinyiyu/http-bridge/logger"
	_ "go.uber.org/automaxprocs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	router := gin.Default()
	router.Use(CorsMiddleware())
	router.Any("/", controller.IndexHttpHandle )

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Sugar.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Sugar.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Sugar.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
		logger.Sugar.Info("timeout of 5 seconds.")
	}
	logger.Sugar.Info("Server exiting")
}

func CorsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger.Sugar.Infof("请求IP", ctx.ClientIP())
		if ctx.Request.Method == "OPTIONS" {
			origin := ctx.GetHeader("Origin")
			if (origin != "") && (origin != "null") {
				ctx.Header("Access-Control-Allow-Origin", origin)
			}else {
				ctx.Header("Access-Control-Allow-Origin", "*")
			}
			ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			ctx.Header("Access-Control-Allow-Credentials", "true")
			ctx.Header("Access-Control-Allow-Headers",
				"Accept, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Token, token")
			ctx.Header("Content-Type","application/json; charset=utf-8")
			ctx.JSON(http.StatusOK, "Options Request!")
		}
		ctx.Next()
	}
}