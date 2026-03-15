// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return
}
