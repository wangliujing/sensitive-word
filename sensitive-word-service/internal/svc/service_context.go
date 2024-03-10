package svc

import (
	"github.com/wangliujing/foundation-framework/jsonrpc"
	"github.com/wangliujing/foundation-framework/rabbitmq"
	"github.com/wangliujing/foundation-framework/reg/nacos"
	"github.com/wangliujing/sensitive-word/internal/config"
	"github.com/wangliujing/sensitive-word/internal/core"
)

type ServiceContext struct {
	Config          *config.Config
	JsonRpcClient   *jsonrpc.Client
	JsonRpcRegister *jsonrpc.Register
	RabbitmqClient  *rabbitmq.Client
	Trie            *core.Trie
}

func NewServiceContext(c *config.Config, callBack func(context *ServiceContext)) *ServiceContext {
	rabbitmqClient := rabbitmq.NewClient(c.RabbitMq, nil, nil)
	jsonRpcService := jsonrpc.Start(c.Nacos.JsonRpcConf.ListenOn)
	registry := nacos.NewRegistry(c.Nacos, c.RpcServerConf)
	jsonRpcRegister := registry.NewJsonRpcRegister(jsonRpcService)
	sc := &ServiceContext{
		Config:          c,
		JsonRpcRegister: jsonRpcRegister,
		JsonRpcClient:   registry.NewJsonRpcClient(),
		RabbitmqClient:  rabbitmqClient,
	}
	go callBack(sc)
	return sc
}
