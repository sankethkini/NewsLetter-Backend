package newsletter

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	newsletterpb "github.com/sankethkini/NewsLetter-Backend/proto/newsletterpb/v1"
)

type API struct {
	CreateNewsLetterEndpoint endpoint.Endpoint
	AddSchemeToNewsEndpoint  endpoint.Endpoint
}

func (e API) CreateNewsLetter(ctx context.Context, req *newsletterpb.CreateNewsLetterRequest) (*newsletterpb.NewsLetter, error) {
	res, err := e.CreateNewsLetterEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := res.(*newsletterpb.NewsLetter)
	return resp, nil
}

func (e API) AddSchemeToNews(ctx context.Context, req *newsletterpb.NewsScheme) (*newsletterpb.NewsScheme, error) {
	res, err := e.AddSchemeToNewsEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := res.(*newsletterpb.NewsScheme)
	return resp, nil
}

func MakeCreateNewsLetterEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*newsletterpb.CreateNewsLetterRequest)
		resp, err := svc.CreateNewsLetter(ctx, req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

func MakeAddSchemeToNewsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*newsletterpb.NewsScheme)
		resp, err := svc.AddSchemeToNews(ctx, req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}
