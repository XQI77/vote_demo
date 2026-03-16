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

type RevokeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRevokeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeLogic {
	return &RevokeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RevokeLogic) Revoke(req *types.RevokeRequest) (resp *types.RevokeResponse, err error) {
	userId, _ := l.ctx.Value(svc.UserIdKey).(string)
	if userId == "" {
		return &types.RevokeResponse{Success: false, Message: "missing X-User-Id header"}, nil
	}

	rpcResp, err := l.svcCtx.VoteService.Revoke(l.ctx, &voteservice.RevokeRequest{
		UserId: userId,
		Topics: req.Topics,
	})
	if err != nil {
		return nil, err
	}

	result := &types.RevokeResponse{
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
