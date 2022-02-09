package user

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/sankethkini/NewsLetter-Backend/proto/userpb/v1"
)

type Endpoints struct {
	CreateUserEndpoint   endpoint.Endpoint
	ValidateUserEndpoint endpoint.Endpoint
	GetEmailEndpoint     endpoint.Endpoint
}

func (e Endpoints) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.User, error) {
	res, err := e.CreateUserEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := res.(*userpb.User)
	return resp, nil
}

func (e Endpoints) ValidateUser(ctx context.Context, req *userpb.ValidateUserRequest) (*userpb.ValidateUserResponse, error) {
	res, err := e.ValidateUserEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := res.(*userpb.ValidateUserResponse)
	return resp, nil
}

func (e Endpoints) GetEmail(ctx context.Context, req *userpb.GetEmailRequest) (*userpb.Email, error) {
	res, err := e.ValidateUserEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := res.(*userpb.Email)
	return resp, nil
}

func MakeCreateUserEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*userpb.CreateUserRequest)
		usr, err := svc.CreateUser(ctx, req)
		if err != nil {
			return nil, err
		}
		return usr, nil
	}
}

func MakeValidateUserEndPoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*userpb.ValidateUserRequest)
		fmt.Println(2, req)
		pass, err := svc.ValidateUser(ctx, req)
		if err != nil {
			return nil, err
		}
		return pass, nil
	}
}

func MakeGetEmailEndPoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*userpb.GetEmailRequest)
		email, err := svc.GetEmail(ctx, req)
		if err != nil {
			return nil, err
		}
		return email, nil
	}
}
