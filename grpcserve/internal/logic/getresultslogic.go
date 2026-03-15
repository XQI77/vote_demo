package logic

import (
	"context"

	"vote-demo/grpcserve/internal/svc"
	"vote-demo/grpcserver/pb/vote-demo/grpcserver/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetResultsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetResultsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetResultsLogic {
	return &GetResultsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询所有话题的实时票数
func (l *GetResultsLogic) GetResults(in *pb.GetResultsRequest) (*pb.GetResultsResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.GetResultsResponse{}, nil
}
