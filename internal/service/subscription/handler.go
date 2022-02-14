package subscription

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	subscriptionpb "github.com/sankethkini/NewsLetter-Backend/proto/subscriptionpb/v1"
)

type Endpoints struct {
	AddUserEndpoint      endpoint.Endpoint
	RemoveUserEndpoint   endpoint.Endpoint
	CreateSchemeEndpoint endpoint.Endpoint
	RenewEndpoint        endpoint.Endpoint
	SearchEndpoint       endpoint.Endpoint
	SortEndpoint         endpoint.Endpoint
	FilterEndpoint       endpoint.Endpoint
	GetUsersEndpoint     endpoint.Endpoint
}

func (e Endpoints) AddUser(ctx context.Context, req *subscriptionpb.AddUserRequest) (*subscriptionpb.AddUserResponse, error) {
	res, err := e.AddUserEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := res.(*subscriptionpb.AddUserResponse)
	return resp, nil
}

func (e Endpoints) RemoveUser(ctx context.Context, req *subscriptionpb.RemoveUserRequest) (*subscriptionpb.RemoveUserResponse, error) {
	res, err := e.RemoveUserEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := res.(*subscriptionpb.RemoveUserResponse)
	return resp, nil
}

func (e Endpoints) CreateScheme(ctx context.Context, req *subscriptionpb.CreateSchemeRequest) (*subscriptionpb.Scheme, error) {
	res, err := e.CreateSchemeEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := res.(*subscriptionpb.Scheme)
	return resp, nil
}

func (e Endpoints) Renew(ctx context.Context, req *subscriptionpb.RenewRequest) (*subscriptionpb.RenewResponse, error) {
	res, err := e.RenewEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := res.(*subscriptionpb.RenewResponse)
	return resp, nil
}

func (e Endpoints) Search(ctx context.Context, req *subscriptionpb.SearchRequest) (*subscriptionpb.SearchResponse, error) {
	res, err := e.RenewEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := res.(*subscriptionpb.SearchResponse)
	return resp, nil
}

func (e Endpoints) Sort(ctx context.Context, req *subscriptionpb.SortRequest) (*subscriptionpb.SortResponse, error) {
	res, err := e.RenewEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := res.(*subscriptionpb.SortResponse)
	return resp, nil
}

func (e Endpoints) Filter(ctx context.Context, req *subscriptionpb.FilterRequest) (*subscriptionpb.FilterResponse, error) {
	res, err := e.RenewEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := res.(*subscriptionpb.FilterResponse)
	return resp, nil
}

func (e Endpoints) GetUsers(ctx context.Context, req *subscriptionpb.GetUsersRequest) (*subscriptionpb.GetUsersResponse, error) {
	res, err := e.GetUsersEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := res.(*subscriptionpb.GetUsersResponse)
	return resp, nil
}

func MakeAddUserEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*subscriptionpb.AddUserRequest)
		usr, err := svc.AddUser(ctx, req)
		if err != nil {
			return nil, err
		}
		return usr, nil
	}
}

func MakeRemoveUserEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*subscriptionpb.RemoveUserRequest)
		usr, err := svc.RemoveUser(ctx, req)
		if err != nil {
			return nil, err
		}
		return usr, nil
	}
}

func MakeCreateSchemeEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*subscriptionpb.CreateSchemeRequest)
		usr, err := svc.CreateScheme(ctx, req)
		if err != nil {
			return nil, err
		}
		return usr, nil
	}
}

func MakeRenewEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*subscriptionpb.RenewRequest)
		usr, err := svc.Renew(ctx, req)
		if err != nil {
			return nil, err
		}
		return usr, nil
	}
}

func MakeSearchEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*subscriptionpb.SearchRequest)
		usr, err := svc.Search(ctx, req)
		if err != nil {
			return nil, err
		}
		return usr, nil
	}
}

func MakeSortEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*subscriptionpb.SortRequest)
		usr, err := svc.Sort(ctx, req)
		if err != nil {
			return nil, err
		}
		return usr, nil
	}
}

func MakeFilterEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*subscriptionpb.FilterRequest)
		usr, err := svc.Filter(ctx, req)
		if err != nil {
			return nil, err
		}
		return usr, nil
	}
}

func MakeGetUsersEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*subscriptionpb.GetUsersRequest)
		usr, err := svc.GetUsers(ctx, req)
		if err != nil {
			return nil, err
		}
		return usr, nil
	}
}
