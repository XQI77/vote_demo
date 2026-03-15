package logic

import (
	"context"

	"vote-demo/grpcserve/internal/svc"
	"vote-demo/grpcserver/pb/vote-demo/grpcserver/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type VoteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VoteLogic {
	return &VoteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 投票（幂等：同一用户同一话题重复调用不会多计）
func (l *VoteLogic) Vote(in *pb.VoteRequest) (*pb.VoteResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.VoteResponse{}, nil
}
