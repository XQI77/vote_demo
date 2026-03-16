// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"vote-demo/grpcserve/voteservice"
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
	userId, _ := l.ctx.Value(svc.UserIdKey).(string)
	if userId == "" {
		return &types.GetUserVotesResponse{}, nil
	}

	rpcResp, err := l.svcCtx.VoteService.GetUserVotes(l.ctx, &voteservice.GetUserVotesRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	return &types.GetUserVotesResponse{VotedTopics: rpcResp.VotedTopics}, nil
}
