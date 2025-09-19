package listener

import (
	"fmt"
	"reflect"

	"github.com/archine/gin-plus/v4/app"
	"github.com/archine/gin-plus/v4/component/gplog"
	"github.com/archine/grpc-demo-proto/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GrpcClientListener gRPC客户端监听器
type GrpcClientListener struct{}

func NewGrpcClientListener() *GrpcClientListener {
	return &GrpcClientListener{}
}

func (g *GrpcClientListener) Order() int {
	return 0
}

func (g *GrpcClientListener) OnContainerRefreshBefore(ctx app.ApplicationContext) {
	conn, err := grpc.NewClient("127.0.0.1:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		gplog.Fatal(fmt.Sprintf("连接 gRPC 服务器失败: %v", err))
	}

	ctx.RegisterBean("helloClient", hello.NewHelloServiceClient(conn), reflect.TypeFor[hello.HelloServiceClient]())
}
