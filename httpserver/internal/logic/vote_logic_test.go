package logic

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"vote-demo/grpcserve/voteservice"
	"vote-demo/grpcserver/pb/vote-demo/grpcserver/pb"
	"vote-demo/httpserver/internal/svc"
	"vote-demo/httpserver/internal/types"

	"google.golang.org/grpc"
)

type fakeVoteService struct {
	voteFunc         func(ctx context.Context, in *voteservice.VoteRequest) (*voteservice.VoteResponse, error)
	revokeFunc       func(ctx context.Context, in *voteservice.RevokeRequest) (*voteservice.RevokeResponse, error)
	getResultsFunc   func(ctx context.Context, in *voteservice.GetResultsRequest) (*voteservice.GetResultsResponse, error)
	getUserVotesFunc func(ctx context.Context, in *voteservice.GetUserVotesRequest) (*voteservice.GetUserVotesResponse, error)
}

func (f *fakeVoteService) Vote(ctx context.Context, in *voteservice.VoteRequest, _ ...grpc.CallOption) (*voteservice.VoteResponse, error) {
	if f.voteFunc == nil {
		return nil, errors.New("unexpected Vote call")
	}
	return f.voteFunc(ctx, in)
}

func (f *fakeVoteService) Revoke(ctx context.Context, in *voteservice.RevokeRequest, _ ...grpc.CallOption) (*voteservice.RevokeResponse, error) {
	if f.revokeFunc == nil {
		return nil, errors.New("unexpected Revoke call")
	}
	return f.revokeFunc(ctx, in)
}

func (f *fakeVoteService) GetResults(ctx context.Context, in *voteservice.GetResultsRequest, _ ...grpc.CallOption) (*voteservice.GetResultsResponse, error) {
	if f.getResultsFunc == nil {
		return nil, errors.New("unexpected GetResults call")
	}
	return f.getResultsFunc(ctx, in)
}

func (f *fakeVoteService) GetUserVotes(ctx context.Context, in *voteservice.GetUserVotesRequest, _ ...grpc.CallOption) (*voteservice.GetUserVotesResponse, error) {
	if f.getUserVotesFunc == nil {
		return nil, errors.New("unexpected GetUserVotes call")
	}
	return f.getUserVotesFunc(ctx, in)
}

func TestVoteRequiresUserIDHeader(t *testing.T) {
	logic := NewVoteLogic(context.Background(), &svc.ServiceContext{
		VoteService: &fakeVoteService{},
	})

	resp, err := logic.Vote(&types.VoteRequest{Topics: []string{"Go"}})
	if err != nil {
		t.Fatalf("Vote returned error: %v", err)
	}
	if resp.Success {
		t.Fatalf("expected failed response, got success")
	}
	if resp.Message != "missing X-User-Id header" {
		t.Fatalf("message = %q", resp.Message)
	}
}

func TestVoteForwardsUserIDAndMapsResults(t *testing.T) {
	ctx := context.WithValue(context.Background(), svc.UserIdKey, "u1")
	service := &fakeVoteService{
		voteFunc: func(ctx context.Context, in *voteservice.VoteRequest) (*voteservice.VoteResponse, error) {
			if in.UserId != "u1" {
				t.Fatalf("user id = %q, want u1", in.UserId)
			}
			if !reflect.DeepEqual(in.Topics, []string{"Go", "K8s"}) {
				t.Fatalf("topics = %#v", in.Topics)
			}
			return &pb.VoteResponse{
				Success: true,
				Message: "ok",
				Results: []*pb.TopicResult{
					{Topic: "Go", Count: 2},
					{Topic: "K8s", Count: 1},
				},
			}, nil
		},
	}

	resp, err := NewVoteLogic(ctx, &svc.ServiceContext{VoteService: service}).Vote(&types.VoteRequest{
		Topics: []string{
			"Go",
			"K8s",
		},
	})
	if err != nil {
		t.Fatalf("Vote returned error: %v", err)
	}

	want := &types.VoteResponse{
		Success: true,
		Message: "ok",
		Results: []types.TopicResult{
			{Topic: "Go", Count: 2},
			{Topic: "K8s", Count: 1},
		},
	}
	if !reflect.DeepEqual(resp, want) {
		t.Fatalf("response = %#v, want %#v", resp, want)
	}
}

func TestRevokeRequiresUserIDHeader(t *testing.T) {
	logic := NewRevokeLogic(context.Background(), &svc.ServiceContext{
		VoteService: &fakeVoteService{},
	})

	resp, err := logic.Revoke(&types.RevokeRequest{Topics: []string{"Go"}})
	if err != nil {
		t.Fatalf("Revoke returned error: %v", err)
	}
	if resp.Success {
		t.Fatalf("expected failed response, got success")
	}
	if resp.Message != "missing X-User-Id header" {
		t.Fatalf("message = %q", resp.Message)
	}
}

func TestRevokeForwardsUserIDAndMapsResults(t *testing.T) {
	ctx := context.WithValue(context.Background(), svc.UserIdKey, "u1")
	service := &fakeVoteService{
		revokeFunc: func(ctx context.Context, in *voteservice.RevokeRequest) (*voteservice.RevokeResponse, error) {
			if in.UserId != "u1" {
				t.Fatalf("user id = %q, want u1", in.UserId)
			}
			if !reflect.DeepEqual(in.Topics, []string{"AI"}) {
				t.Fatalf("topics = %#v", in.Topics)
			}
			return &pb.RevokeResponse{
				Success: true,
				Message: "ok",
				Results: []*pb.TopicResult{
					{Topic: "AI", Count: 0},
				},
			}, nil
		},
	}

	resp, err := NewRevokeLogic(ctx, &svc.ServiceContext{VoteService: service}).Revoke(&types.RevokeRequest{
		Topics: []string{"AI"},
	})
	if err != nil {
		t.Fatalf("Revoke returned error: %v", err)
	}

	want := &types.RevokeResponse{
		Success: true,
		Message: "ok",
		Results: []types.TopicResult{
			{Topic: "AI", Count: 0},
		},
	}
	if !reflect.DeepEqual(resp, want) {
		t.Fatalf("response = %#v, want %#v", resp, want)
	}
}

func TestGetUserVotesHandlesMissingUserIDAndMapsResponse(t *testing.T) {
	missingResp, err := NewGetUserVotesLogic(context.Background(), &svc.ServiceContext{
		VoteService: &fakeVoteService{},
	}).GetUserVotes()
	if err != nil {
		t.Fatalf("GetUserVotes without user returned error: %v", err)
	}
	if len(missingResp.VotedTopics) != 0 {
		t.Fatalf("missing user topics = %#v, want empty", missingResp.VotedTopics)
	}

	ctx := context.WithValue(context.Background(), svc.UserIdKey, "u2")
	service := &fakeVoteService{
		getUserVotesFunc: func(ctx context.Context, in *voteservice.GetUserVotesRequest) (*voteservice.GetUserVotesResponse, error) {
			if in.UserId != "u2" {
				t.Fatalf("user id = %q, want u2", in.UserId)
			}
			return &pb.GetUserVotesResponse{
				VotedTopics: []string{"Go", "AI"},
			}, nil
		},
	}

	resp, err := NewGetUserVotesLogic(ctx, &svc.ServiceContext{VoteService: service}).GetUserVotes()
	if err != nil {
		t.Fatalf("GetUserVotes returned error: %v", err)
	}
	if !reflect.DeepEqual(resp.VotedTopics, []string{"Go", "AI"}) {
		t.Fatalf("topics = %#v", resp.VotedTopics)
	}
}

func TestGetResultsMapsRPCResults(t *testing.T) {
	service := &fakeVoteService{
		getResultsFunc: func(ctx context.Context, in *voteservice.GetResultsRequest) (*voteservice.GetResultsResponse, error) {
			if in == nil {
				t.Fatalf("request is nil")
			}
			return &pb.GetResultsResponse{
				Results: []*pb.TopicResult{
					{Topic: "Go", Count: 3},
					{Topic: "K8s", Count: 4},
				},
			}, nil
		},
	}

	resp, err := NewGetResultsLogic(context.Background(), &svc.ServiceContext{VoteService: service}).GetResults()
	if err != nil {
		t.Fatalf("GetResults returned error: %v", err)
	}

	want := &types.GetResultsResponse{
		Results: []types.TopicResult{
			{Topic: "Go", Count: 3},
			{Topic: "K8s", Count: 4},
		},
	}
	if !reflect.DeepEqual(resp, want) {
		t.Fatalf("response = %#v, want %#v", resp, want)
	}
}
