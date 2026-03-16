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

type GetResultsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetResultsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetResultsLogic {
	return &GetResultsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetResultsLogic) GetResults() (resp *types.GetResultsResponse, err error) {
	rpcResp, err := l.svcCtx.VoteService.GetResults(l.ctx, &voteservice.GetResultsRequest{})
	if err != nil {
		return nil, err
	}

	result := &types.GetResultsResponse{}
	for _, r := range rpcResp.Results {
		result.Results = append(result.Results, types.TopicResult{
			Topic: r.Topic,
			Count: r.Count,
		})
	}
	return result, nil
}
