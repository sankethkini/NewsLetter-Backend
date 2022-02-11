package transport

import (
	"context"

	"github.com/sankethkini/NewsLetter-Backend/internal/service/user"

	"github.com/go-kit/kit/transport/grpc"

	userpb "github.com/sankethkini/NewsLetter-Backend/proto/userpb/v1"
)

type UserGrpcServer struct {
	userpb.UnimplementedUserServiceServer
	createuser   grpc.Handler
	validateuser grpc.Handler
	getemail     grpc.Handler
}

func NewUserGrpcServer(ctx context.Context, svc user.UserService) userpb.UserServiceServer {
	return &UserGrpcServer{
		createuser: grpc.NewServer(
			user.MakeCreateUserEndpoint(svc),
			decodeCreateUserRequest,
			func(c context.Context, i interface{}) (request interface{}, err error) { return i, nil },
		),
		validateuser: grpc.NewServer(
			user.MakeValidateUserEndPoint(svc),
			decodeValidateUserRequest,
			func(c context.Context, i interface{}) (request interface{}, err error) { return i, nil },
		),
		getemail: grpc.NewServer(
			user.MakeGetEmailEndPoint(svc),
			decodeGetEmailRequest,
			func(c context.Context, i interface{}) (request interface{}, err error) { return i, nil },
		),
	}
}

func (g *UserGrpcServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.User, error) {
	_, resp, err := g.createuser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*userpb.User), nil
}

func (g *UserGrpcServer) ValidateUser(ctx context.Context, req *userpb.ValidateUserRequest) (*userpb.ValidateUserResponse, error) {
	_, resp, err := g.validateuser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*userpb.ValidateUserResponse), nil
}

func (g *UserGrpcServer) GetEmail(ctx context.Context, req *userpb.GetEmailRequest) (*userpb.Email, error) {
	_, resp, err := g.getemail.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*userpb.Email), nil
}

func decodeCreateUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*userpb.CreateUserRequest)
	return req, nil
}

func decodeValidateUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*userpb.ValidateUserRequest)
	return req, nil
}

func decodeGetEmailRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*userpb.GetEmailRequest)
	return req, nil
}
