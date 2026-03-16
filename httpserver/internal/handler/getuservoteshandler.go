// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"vote-demo/httpserver/internal/logic"
	"vote-demo/httpserver/internal/svc"
)

func GetUserVotesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), svc.UserIdKey, r.Header.Get("X-User-Id"))
		l := logic.NewGetUserVotesLogic(ctx, svcCtx)
		resp, err := l.GetUserVotes()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
