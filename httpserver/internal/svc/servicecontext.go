// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"vote-demo/grpcserve/voteservice"
	"vote-demo/httpserver/internal/config"
	"github.com/zeromicro/go-zero/zrpc"
)

type contextKey string

const UserIdKey contextKey = "userId"

type ServiceContext struct {
	Config config.Config
	VoteService voteservice.VoteService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		VoteService: voteservice.NewVoteService(zrpc.MustNewClient(c.VoteRpc)),
	}
}
