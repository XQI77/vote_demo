// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import "github.com/zeromicro/go-zero/rest"
import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	rest.RestConf
	VoteRpc zrpc.RpcClientConf
}
