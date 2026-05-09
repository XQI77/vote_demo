package logic

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"vote-demo/grpcserve/internal/config"
	"vote-demo/grpcserve/internal/svc"
	"vote-demo/grpcserver/pb/vote-demo/grpcserver/pb"

	"github.com/zeromicro/go-zero/core/stores/redis/redistest"
)

func newTestServiceContext(t *testing.T) *svc.ServiceContext {
	t.Helper()

	rds, clean := redistest.CreateRedisWithClean(t)
	t.Cleanup(clean)

	return &svc.ServiceContext{
		Config: config.Config{
			Topics: []string{"Go", "K8s", "AI"},
		},
		Redis: rds,
	}
}

func countsByTopic(results []*pb.TopicResult) map[string]int64 {
	counts := make(map[string]int64, len(results))
	for _, result := range results {
		counts[result.Topic] = result.Count
	}
	return counts
}

func TestVoteRejectsInvalidRequests(t *testing.T) {
	ctx := newTestServiceContext(t)
	logic := NewVoteLogic(context.Background(), ctx)

	tests := []struct {
		name string
		req  *pb.VoteRequest
		msg  string
	}{
		{
			name: "missing user",
			req:  &pb.VoteRequest{Topics: []string{"Go"}},
			msg:  "missing user_id",
		},
		{
			name: "missing topics",
			req:  &pb.VoteRequest{UserId: "u1"},
			msg:  "no topics provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := logic.Vote(tt.req)
			if err != nil {
				t.Fatalf("Vote returned error: %v", err)
			}
			if resp.Success {
				t.Fatalf("expected failed response, got success")
			}
			if resp.Message != tt.msg {
				t.Fatalf("message = %q, want %q", resp.Message, tt.msg)
			}
		})
	}
}

func TestVoteIsIdempotentPerUserAndTopic(t *testing.T) {
	ctx := newTestServiceContext(t)
	logic := NewVoteLogic(context.Background(), ctx)

	resp, err := logic.Vote(&pb.VoteRequest{
		UserId: "u1",
		Topics: []string{
			"Go",
			"K8s",
		},
	})
	if err != nil {
		t.Fatalf("Vote returned error: %v", err)
	}
	if !resp.Success {
		t.Fatalf("expected successful vote, got message %q", resp.Message)
	}

	resp, err = logic.Vote(&pb.VoteRequest{
		UserId: "u1",
		Topics: []string{
			"Go",
		},
	})
	if err != nil {
		t.Fatalf("second Vote returned error: %v", err)
	}
	if !resp.Success {
		t.Fatalf("expected idempotent vote success, got message %q", resp.Message)
	}

	wantCounts := map[string]int64{
		"Go":  1,
		"K8s": 1,
		"AI":  0,
	}
	if got := countsByTopic(resp.Results); !reflect.DeepEqual(got, wantCounts) {
		t.Fatalf("counts = %#v, want %#v", got, wantCounts)
	}
}

func TestRevokeSpecificAndAllVotes(t *testing.T) {
	ctx := newTestServiceContext(t)
	voteLogic := NewVoteLogic(context.Background(), ctx)
	revokeLogic := NewRevokeLogic(context.Background(), ctx)

	if _, err := voteLogic.Vote(&pb.VoteRequest{UserId: "u1", Topics: []string{"Go", "K8s"}}); err != nil {
		t.Fatalf("Vote returned error: %v", err)
	}
	if _, err := voteLogic.Vote(&pb.VoteRequest{UserId: "u2", Topics: []string{"Go"}}); err != nil {
		t.Fatalf("Vote returned error: %v", err)
	}

	resp, err := revokeLogic.Revoke(&pb.RevokeRequest{UserId: "u1", Topics: []string{"Go"}})
	if err != nil {
		t.Fatalf("Revoke returned error: %v", err)
	}
	if !resp.Success {
		t.Fatalf("expected successful revoke, got message %q", resp.Message)
	}

	wantCounts := map[string]int64{
		"Go":  1,
		"K8s": 1,
		"AI":  0,
	}
	if got := countsByTopic(resp.Results); !reflect.DeepEqual(got, wantCounts) {
		t.Fatalf("counts after specific revoke = %#v, want %#v", got, wantCounts)
	}

	resp, err = revokeLogic.Revoke(&pb.RevokeRequest{UserId: "u1"})
	if err != nil {
		t.Fatalf("Revoke all returned error: %v", err)
	}

	wantCounts = map[string]int64{
		"Go":  1,
		"K8s": 0,
		"AI":  0,
	}
	if got := countsByTopic(resp.Results); !reflect.DeepEqual(got, wantCounts) {
		t.Fatalf("counts after revoke all = %#v, want %#v", got, wantCounts)
	}

	resp, err = revokeLogic.Revoke(&pb.RevokeRequest{UserId: "u1", Topics: []string{"K8s"}})
	if err != nil {
		t.Fatalf("repeated Revoke returned error: %v", err)
	}
	if got := countsByTopic(resp.Results); !reflect.DeepEqual(got, wantCounts) {
		t.Fatalf("counts after repeated revoke = %#v, want %#v", got, wantCounts)
	}
}

func TestGetResultsAndUserVotes(t *testing.T) {
	ctx := newTestServiceContext(t)

	if _, err := NewVoteLogic(context.Background(), ctx).Vote(&pb.VoteRequest{
		UserId: "u1",
		Topics: []string{
			"Go",
			"AI",
		},
	}); err != nil {
		t.Fatalf("Vote returned error: %v", err)
	}

	resultsResp, err := NewGetResultsLogic(context.Background(), ctx).GetResults(&pb.GetResultsRequest{
		Topics: []string{
			"AI",
			"Go",
		},
	})
	if err != nil {
		t.Fatalf("GetResults returned error: %v", err)
	}

	wantCounts := map[string]int64{
		"AI": 1,
		"Go": 1,
	}
	if got := countsByTopic(resultsResp.Results); !reflect.DeepEqual(got, wantCounts) {
		t.Fatalf("counts = %#v, want %#v", got, wantCounts)
	}

	userVotesResp, err := NewGetUserVotesLogic(context.Background(), ctx).GetUserVotes(&pb.GetUserVotesRequest{UserId: "u1"})
	if err != nil {
		t.Fatalf("GetUserVotes returned error: %v", err)
	}

	gotTopics := append([]string(nil), userVotesResp.VotedTopics...)
	sort.Strings(gotTopics)

	wantTopics := []string{"AI", "Go"}
	if !reflect.DeepEqual(gotTopics, wantTopics) {
		t.Fatalf("voted topics = %#v, want %#v", gotTopics, wantTopics)
	}
}
