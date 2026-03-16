// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"vote-demo/httpserver/internal/logic"
	"vote-demo/httpserver/internal/svc"
	"vote-demo/httpserver/internal/types"
)

func RevokeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RevokeRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		ctx := context.WithValue(r.Context(), svc.UserIdKey, r.Header.Get("X-User-Id"))
		l := logic.NewRevokeLogic(ctx, svcCtx)
		resp, err := l.Revoke(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
