package logic

import (
	"context"
	"fmt"

	"vote-demo/grpcserve/internal/svc"
	"vote-demo/grpcserver/pb/vote-demo/grpcserver/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

// Lua 脚本：原子地完成"检查用户是否已投 + 写入 + 计数+1"
// KEYS[1] = vote:user:{userId}  (Set，记录该用户已投的话题)
// KEYS[2] = vote:count:{topic}  (String，话题票数)
// ARGV[1] = topic
// 返回 1 表示投票成功，0 表示已投过（幂等）
const voteScript = `
local userKey   = KEYS[1]
local countKey  = KEYS[2]
local topic     = ARGV[1]
if redis.call('SISMEMBER', userKey, topic) == 1 then
    return 0
end
redis.call('SADD', userKey, topic)
redis.call('INCR', countKey)
return 1
`

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

func (l *VoteLogic) Vote(in *pb.VoteRequest) (*pb.VoteResponse, error) {
	if in.UserId == "" {
		return &pb.VoteResponse{Success: false, Message: "missing user_id"}, nil
	}
	if len(in.Topics) == 0 {
		return &pb.VoteResponse{Success: false, Message: "no topics provided"}, nil
	}

	userKey := fmt.Sprintf("vote:user:%s", in.UserId)

	for _, topic := range in.Topics {
		countKey := fmt.Sprintf("vote:count:%s", topic)
		_, err := l.svcCtx.Redis.EvalCtx(l.ctx, voteScript,
			[]string{userKey, countKey}, topic)
		if err != nil {
			l.Errorf("vote lua script error, topic=%s err=%v", topic, err)
			return &pb.VoteResponse{Success: false, Message: "internal error"}, nil
		}
	}

	results, err := getResults(l.ctx, l.svcCtx, l.svcCtx.Config.Topics)
	if err != nil {
		return &pb.VoteResponse{Success: false, Message: "failed to get results"}, nil
	}
	return &pb.VoteResponse{Success: true, Message: "ok", Results: results}, nil
}
