// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"vote-demo/httpserver/internal/svc"
	"vote-demo/httpserver/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VoteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VoteLogic {
	return &VoteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VoteLogic) Vote(req *types.VoteRequest) (resp *types.VoteResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
