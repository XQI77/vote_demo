package logic

import (
	"context"

	"vote-demo/grpcserve/internal/svc"
	"vote-demo/grpcserver/pb/vote-demo/grpcserver/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserVotesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserVotesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserVotesLogic {
	return &GetUserVotesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询某用户已投的话题（用于前端回显状态）
func (l *GetUserVotesLogic) GetUserVotes(in *pb.GetUserVotesRequest) (*pb.GetUserVotesResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.GetUserVotesResponse{}, nil
}
