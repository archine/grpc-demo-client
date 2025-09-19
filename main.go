package main

import (
	_ "grpc-demo-client/controller"
	"grpc-demo-client/listener"

	ginplus "github.com/archine/gin-plus/v4"
)

func main() {
	ginplus.New().
		With(
			ginplus.WithEvent(
				listener.NewGrpcClientListener(),
			),
		).
		Run(ginplus.ServerMode)
}
