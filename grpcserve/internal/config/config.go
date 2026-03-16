package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Store redis.RedisConf

	// 支持的话题列表，启动时从配置读取
	Topics []string
}
