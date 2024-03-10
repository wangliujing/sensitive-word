package config

import (
	"github.com/wangliujing/foundation-framework/rabbitmq"
	"github.com/wangliujing/foundation-framework/reg/nacos"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Nacos    nacos.Conf
	RabbitMq rabbitmq.Conf
	/*SensitiveWord SensitiveWord*/
}

/*type SensitiveWord struct {
	ReconstructionDelay time.Duration `json:",default=60s"` // 单位秒
}*/
