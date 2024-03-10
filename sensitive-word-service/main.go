package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/wangliujing/foundation-framework/conf/nacos"
	"github.com/wangliujing/foundation-framework/system"
	"github.com/wangliujing/sensitive-word/internal/config"
	"github.com/wangliujing/sensitive-word/internal/jsonrpc"
	"github.com/wangliujing/sensitive-word/internal/listener"
	"github.com/wangliujing/sensitive-word/internal/logic"
	"github.com/wangliujing/sensitive-word/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/test.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	nacos.MustLoad(*configFile, &c)
	logx.MustSetup(c.Log)
	svc.NewServiceContext(&c, func(serviceContext *svc.ServiceContext) {
		sensitiveWordLogic := logic.NewSensitiveWordLogic(context.Background(), serviceContext)
		err := sensitiveWordLogic.InitTrie()
		if err != nil {
			logx.Must(err)
		}
		// 注册监听器
		serviceContext.RabbitmqClient.RegisterListener(listener.NewWordChangeListener(c, serviceContext))
		// 注册jsonRpc
		serviceContext.JsonRpcRegister.RegisterService("SensitiveWordRpcService",
			jsonrpc.NewSensitiveWordRpcService(serviceContext), "")
	})

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	system.Start(s)
}
