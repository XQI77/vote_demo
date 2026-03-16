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
	userId, _ := l.ctx.Value(svc.UserIdKey).(string)
	if userId == "" {
		return &types.VoteResponse{Success: false, Message: "missing X-User-Id header"}, nil
	}

	rpcResp, err := l.svcCtx.VoteService.Vote(l.ctx, &voteservice.VoteRequest{
		UserId: userId,
		Topics: req.Topics,
	})
	if err != nil {
		return nil, err
	}

	result := &types.VoteResponse{
		Success: rpcResp.Success,
		Message: rpcResp.Message,
	}
	for _, r := range rpcResp.Results {
		result.Results = append(result.Results, types.TopicResult{
			Topic: r.Topic,
			Count: r.Count,
		})
	}
	return result, nil
}
