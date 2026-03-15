package logic

import (
	"context"
	"fmt"

	"vote-demo/grpcserve/internal/svc"
	"vote-demo/grpcserver/pb/vote-demo/grpcserver/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

// Lua 脚本：原子地完成"检查用户是否已投该话题 + 移除 + 计数-1"
// KEYS[1] = vote:user:{userId}
// KEYS[2] = vote:count:{topic}
// ARGV[1] = topic
// 返回 1 表示撤销成功，0 表示本来就没投
const revokeScript = `
local userKey  = KEYS[1]
local countKey = KEYS[2]
local topic    = ARGV[1]
if redis.call('SISMEMBER', userKey, topic) == 0 then
    return 0
end
redis.call('SREM', userKey, topic)
local cur = tonumber(redis.call('GET', countKey) or "0")
if cur > 0 then
    redis.call('DECR', countKey)
end
return 1
`

type RevokeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRevokeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeLogic {
	return &RevokeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RevokeLogic) Revoke(in *pb.RevokeRequest) (*pb.RevokeResponse, error) {
	if in.UserId == "" {
		return &pb.RevokeResponse{Success: false, Message: "missing user_id"}, nil
	}

	userKey := fmt.Sprintf("vote:user:%s", in.UserId)

	// Topics 为空则撤销该用户所有已投话题
	topics := in.Topics
	if len(topics) == 0 {
		var err error
		topics, err = l.svcCtx.Redis.SmembersCtx(l.ctx, userKey)
		if err != nil {
			return &pb.RevokeResponse{Success: false, Message: "internal error"}, nil
		}
	}

	for _, topic := range topics {
		countKey := fmt.Sprintf("vote:count:%s", topic)
		_, err := l.svcCtx.Redis.EvalCtx(l.ctx, revokeScript,
			[]string{userKey, countKey}, topic)
		if err != nil {
			l.Errorf("revoke lua script error, topic=%s err=%v", topic, err)
			return &pb.RevokeResponse{Success: false, Message: "internal error"}, nil
		}
	}

	results, err := getResults(l.ctx, l.svcCtx, l.svcCtx.Config.Topics)
	if err != nil {
		return &pb.RevokeResponse{Success: false, Message: "failed to get results"}, nil
	}
	return &pb.RevokeResponse{Success: true, Message: "ok", Results: results}, nil
}
