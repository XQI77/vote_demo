// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"vote-demo/httpserver/internal/svc"
	"vote-demo/httpserver/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserVotesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserVotesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserVotesLogic {
	return &GetUserVotesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserVotesLogic) GetUserVotes() (resp *types.GetUserVotesResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
