package logic

import (
	"context"
	"fmt"

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

func (l *GetUserVotesLogic) GetUserVotes(in *pb.GetUserVotesRequest) (*pb.GetUserVotesResponse, error) {
	if in.UserId == "" {
		return &pb.GetUserVotesResponse{}, nil
	}

	userKey := fmt.Sprintf("vote:user:%s", in.UserId)
	topics, err := l.svcCtx.Redis.SmembersCtx(l.ctx, userKey)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserVotesResponse{VotedTopics: topics}, nil
}
