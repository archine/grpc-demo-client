package controller

import (
	"context"
	"time"

	"github.com/archine/gin-plus/v4/component/ioc"
	"github.com/archine/gin-plus/v4/resp"
	"github.com/archine/grpc-demo-proto/hello"
	"github.com/gin-gonic/gin"
)

func init() {
	ioc.RegisterBeanDef(&HelloController{})
}

// HelloController 调用 gRPC Hello 服务
type HelloController struct {
	helloClient hello.HelloServiceClient `autowire:""` // gRPC Hello 服务客户端
}

func (h *HelloController) SetRoutes(router *gin.RouterGroup) {
	router.GET("/hello", h.SayHello)
}

// SayHello 调用 gRPC Hello 服务的 SayHello 方法
func (h *HelloController) SayHello(ctx *gin.Context) {
	name := ctx.Query("name")

	apiCtx, cancelFunc := context.WithTimeout(ctx, 3*time.Second)
	defer cancelFunc()

	response, err := h.helloClient.SayHello(apiCtx, &hello.HelloRequest{Name: name})
	if err != nil {
		resp.BadRequest(ctx, err.Error())
		return
	}

	resp.Json(ctx, response.Message)
}
