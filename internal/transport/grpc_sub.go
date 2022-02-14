package transport

import (
	"context"

	"github.com/go-kit/kit/transport/grpc"
	"github.com/sankethkini/NewsLetter-Backend/internal/service/subscription"
	subscriptionpb "github.com/sankethkini/NewsLetter-Backend/proto/subscriptionpb/v1"
)

type SubscriptionGrpcServer struct {
	subscriptionpb.UnimplementedSubscriptionServiceServer
	adduser      grpc.Handler
	removeuser   grpc.Handler
	createscheme grpc.Handler
	renew        grpc.Handler
	search       grpc.Handler
	sort         grpc.Handler
	filter       grpc.Handler
	getusers     grpc.Handler
}

func NewSubscriptionService(ctx context.Context, svc subscription.Service) subscriptionpb.SubscriptionServiceServer {
	return &SubscriptionGrpcServer{
		adduser: grpc.NewServer(
			subscription.MakeAddUserEndpoint(svc),
			decodeAddUserRequest,
			func(c context.Context, i interface{}) (request interface{}, err error) { return i, nil },
		),
		removeuser: grpc.NewServer(
			subscription.MakeRemoveUserEndpoint(svc),
			decodeRemoveUserRequest,
			func(c context.Context, i interface{}) (response interface{}, err error) { return i, nil },
		),
		createscheme: grpc.NewServer(
			subscription.MakeCreateSchemeEndpoint(svc),
			decodeCreateSchemeRequest,
			func(c context.Context, i interface{}) (response interface{}, err error) { return i, nil },
		),
		renew: grpc.NewServer(
			subscription.MakeRenewEndpoint(svc),
			decodeRenewRequest,
			func(c context.Context, i interface{}) (response interface{}, err error) { return i, nil },
		),
		search: grpc.NewServer(
			subscription.MakeSearchEndpoint(svc),
			decodeSearchRequest,
			func(c context.Context, i interface{}) (response interface{}, err error) { return i, nil },
		),
		sort: grpc.NewServer(
			subscription.MakeSortEndpoint(svc),
			decodeSortRequest,
			func(c context.Context, i interface{}) (response interface{}, err error) { return i, nil },
		),
		filter: grpc.NewServer(
			subscription.MakeFilterEndpoint(svc),
			decodeFilterRequest,
			func(c context.Context, i interface{}) (response interface{}, err error) { return i, nil },
		),
		getusers: grpc.NewServer(
			subscription.MakeGetUsersEndpoint(svc),
			decodeGetUsersRequest,
			func(c context.Context, i interface{}) (response interface{}, err error) { return i, nil },
		),
	}
}

func (g *SubscriptionGrpcServer) AddUser(ctx context.Context, req *subscriptionpb.AddUserRequest) (*subscriptionpb.AddUserResponse, error) {
	_, resp, err := g.adduser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*subscriptionpb.AddUserResponse), nil
}

func (g *SubscriptionGrpcServer) RemoveUser(ctx context.Context, req *subscriptionpb.RemoveUserRequest) (*subscriptionpb.RemoveUserResponse, error) {
	_, resp, err := g.removeuser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*subscriptionpb.RemoveUserResponse), nil
}

func (g *SubscriptionGrpcServer) CreateScheme(ctx context.Context, req *subscriptionpb.CreateSchemeRequest) (*subscriptionpb.Scheme, error) {
	_, resp, err := g.createscheme.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*subscriptionpb.Scheme), nil
}

func (g *SubscriptionGrpcServer) Renew(ctx context.Context, req *subscriptionpb.RenewRequest) (*subscriptionpb.RenewResponse, error) {
	_, resp, err := g.renew.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*subscriptionpb.RenewResponse), nil
}

func (g *SubscriptionGrpcServer) Search(ctx context.Context, req *subscriptionpb.SearchRequest) (*subscriptionpb.SearchResponse, error) {
	_, resp, err := g.search.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*subscriptionpb.SearchResponse), nil
}

func (g *SubscriptionGrpcServer) Sort(ctx context.Context, req *subscriptionpb.SortRequest) (*subscriptionpb.SortResponse, error) {
	_, resp, err := g.sort.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*subscriptionpb.SortResponse), nil
}

func (g *SubscriptionGrpcServer) Filter(ctx context.Context, req *subscriptionpb.FilterRequest) (*subscriptionpb.FilterResponse, error) {
	_, resp, err := g.filter.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*subscriptionpb.FilterResponse), nil
}

func (g *SubscriptionGrpcServer) GetUsers(ctx context.Context, req *subscriptionpb.GetUsersRequest) (*subscriptionpb.GetUsersResponse, error) {
	_, resp, err := g.getusers.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*subscriptionpb.GetUsersResponse), nil
}

func decodeAddUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*subscriptionpb.AddUserRequest)
	return req, nil
}

func decodeRemoveUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*subscriptionpb.RemoveUserRequest)
	return req, nil
}

func decodeCreateSchemeRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*subscriptionpb.CreateSchemeRequest)
	return req, nil
}

func decodeRenewRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*subscriptionpb.RenewRequest)
	return req, nil
}

func decodeSearchRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*subscriptionpb.SearchRequest)
	return req, nil
}

func decodeSortRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*subscriptionpb.SortRequest)
	return req, nil
}

func decodeFilterRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*subscriptionpb.FilterRequest)
	return req, nil
}

func decodeGetUsersRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*subscriptionpb.GetUsersRequest)
	return req, nil
}
