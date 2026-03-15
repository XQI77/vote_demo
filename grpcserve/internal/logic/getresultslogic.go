package logic

import (
	"context"
	"fmt"
	"strconv"

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

func (l *GetResultsLogic) GetResults(in *pb.GetResultsRequest) (*pb.GetResultsResponse, error) {
	topics := in.Topics
	if len(topics) == 0 {
		topics = l.svcCtx.Config.Topics
	}

	results, err := getResults(l.ctx, l.svcCtx, topics)
	if err != nil {
		return nil, err
	}
	return &pb.GetResultsResponse{Results: results}, nil
}

// getResults 是被多个 logic 复用的公共函数
func getResults(ctx context.Context, svcCtx *svc.ServiceContext, topics []string) ([]*pb.TopicResult, error) {
	results := make([]*pb.TopicResult, 0, len(topics))
	for _, topic := range topics {
		countKey := fmt.Sprintf("vote:count:%s", topic)
		val, err := svcCtx.Redis.GetCtx(ctx, countKey)
		if err != nil {
			// key 不存在时 go-zero redis 返回 "", err=nil；其他错误才中断
			return nil, err
		}
		count, _ := strconv.ParseInt(val, 10, 64)
		results = append(results, &pb.TopicResult{Topic: topic, Count: count})
	}
	return results, nil
}
