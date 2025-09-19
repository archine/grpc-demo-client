## gRPC Demo 客户端

### 必须
- Go 1.24 or later

### 项目结构

```
grpc-demo-client
├── controller/           # 控制器
├── listener/           # 监听器，主要用于创建grpc客户端连接
├── main.go               # 项目入口