package transport

import (
	"context"

	"github.com/go-kit/kit/transport/grpc"
	"github.com/sankethkini/NewsLetter-Backend/internal/service/admin"
	adminpb "github.com/sankethkini/NewsLetter-Backend/proto/adminpb/v1"
)

type AdminGrpcServer struct {
	adminpb.UnimplementedAdminServiceServer
	signin grpc.Handler
}

func NewAdminGrpcServer(ctx context.Context, svc admin.Service) adminpb.AdminServiceServer {
	return &AdminGrpcServer{
		signin: grpc.NewServer(
			admin.MakeSingInEndpoint(svc),
			decodeSiginRequest,
			func(c context.Context, i interface{}) (response interface{}, err error) { return i, nil },
		),
	}
}

func (g *AdminGrpcServer) SignIn(ctx context.Context, req *adminpb.SignInRequest) (*adminpb.SignInResponse, error) {
	_, resp, err := g.signin.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*adminpb.SignInResponse), nil
}

func decodeSiginRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*adminpb.SignInRequest)
	return req, nil
}
