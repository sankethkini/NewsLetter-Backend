package admin

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	adminpb "github.com/sankethkini/NewsLetter-Backend/proto/adminpb/v1"
)

type Endpoints struct {
	SignIn endpoint.Endpoint
}

func (e Endpoints) SingIn(ctx context.Context, req *adminpb.SignInRequest) (*adminpb.SignInResponse, error) {
	res, err := e.SignIn(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := res.(*adminpb.SignInResponse)
	return resp, nil
}

func MakeSingInEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*adminpb.SignInRequest)
		resp, err := svc.SingIn(ctx, req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}
