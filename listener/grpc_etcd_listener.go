package listener

import (
	"fmt"
	"reflect"
	"time"

	"github.com/archine/gin-plus/v4/app"
	"github.com/archine/gin-plus/v4/component/gplog"
	"github.com/archine/grpc-demo-proto/hello"
	"github.com/archine/grpc-demo-proto/user"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	gresolver "google.golang.org/grpc/resolver"
)

// conf etcd配置
type conf struct {
	Endpoints   []string      `mapstructure:"endpoints"`    // etcd地址列表
	TTL         int64         `mapstructure:"ttl"`          // 租约时间，单位秒
	DialTimeout time.Duration `mapstructure:"dial-timeout"` // 连接超时时间
}

func (c *conf) verify() {
	if len(c.Endpoints) == 0 {
		c.Endpoints = []string{"127.0.0.1:2379"}
	}
	if c.TTL == 0 {
		c.TTL = 10
	}
	if c.DialTimeout == 0 {
		c.DialTimeout = 3 * time.Second
	}
}

// GrpcEtcdClientListener gRPC etcd客户端监听器
type GrpcEtcdClientListener struct{}

func NewGrpcEtcdClientListener() *GrpcEtcdClientListener {
	return &GrpcEtcdClientListener{}
}

func (g *GrpcEtcdClientListener) Order() int {
	return 0
}

func (g *GrpcEtcdClientListener) OnContainerRefreshBefore(ctx app.ApplicationContext) {
	var cfg conf
	if err := ctx.GetConfigProvider().Unmarshal("etcd", &cfg); err != nil {
		gplog.Fatal(fmt.Sprintf("Init etcd client failed, unable to parse configuration: %v", err))
	}
	cfg.verify()

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: cfg.DialTimeout,
	})
	if err != nil {
		gplog.Fatal(fmt.Sprintf("Init etcd client failed, unable to connect to etcd: %v", err))
	}

	builder, err := resolver.NewBuilder(cli)
	if err != nil {
		gplog.Fatal(fmt.Sprintf("Init etcd client failed, unable to create resolver: %v", err))
	}

	newTestServerConn(ctx, builder)
}

// newTestServerConn 创建到 testServer 的连接，并注册相关的客户端Bean
// testServer 是服务提供者在 etcd 上注册的服务名称
// 注意：同一个服务名称里面是会有多个服务实例的，每个实例对应一个不同的 proto
func newTestServerConn(ctx app.ApplicationContext, builder gresolver.Builder) {
	conn, err := grpc.NewClient("etcd:///testServer",
		grpc.WithResolvers(builder),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		gplog.Fatal(fmt.Sprintf("Init grpc client failed, unable to connect to testServer: %v", err))
	}

	ctx.RegisterBean("helloClient", hello.NewHelloServiceClient(conn), reflect.TypeFor[hello.HelloServiceClient]())
	ctx.RegisterBean("userClient", user.NewUserServiceClient(conn), reflect.TypeFor[user.UserServiceClient]())
	// ... 注册本 gRPC 服务里其他的客户端Bean

	gplog.Info("Init grpc client success, connected to testServer")
}
