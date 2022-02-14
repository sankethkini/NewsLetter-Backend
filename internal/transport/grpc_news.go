package transport

import (
	"context"

	"github.com/go-kit/kit/transport/grpc"
	newsletter "github.com/sankethkini/NewsLetter-Backend/internal/service/news_letter"
	newsletterpb "github.com/sankethkini/NewsLetter-Backend/proto/newsletterpb/v1"
)

type NewsGrpcServer struct {
	newsletterpb.UnimplementedNewsLetterServiceServer
	createnews      grpc.Handler
	addschemetonews grpc.Handler
}

func NewNewsGrpcServer(ctx context.Context, svc newsletter.Service) newsletterpb.NewsLetterServiceServer {
	return &NewsGrpcServer{
		createnews: grpc.NewServer(
			newsletter.MakeCreateNewsLetterEndpoint(svc),
			decodeCreateNewsRequest,
			func(c context.Context, i interface{}) (response interface{}, err error) { return i, nil },
		),
		addschemetonews: grpc.NewServer(
			newsletter.MakeAddSchemeToNewsEndpoint(svc),
			decodeAddSchemeRequest,
			func(c context.Context, i interface{}) (response interface{}, err error) { return i, nil },
		),
	}
}

func (g *NewsGrpcServer) CreateNewsLetter(ctx context.Context, req *newsletterpb.CreateNewsLetterRequest) (*newsletterpb.NewsLetter, error) {
	_, resp, err := g.createnews.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*newsletterpb.NewsLetter), nil
}

func (g *NewsGrpcServer) AddSchemeToNews(ctx context.Context, req *newsletterpb.NewsScheme) (*newsletterpb.NewsScheme, error) {
	_, resp, err := g.addschemetonews.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*newsletterpb.NewsScheme), nil
}

func decodeCreateNewsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*newsletterpb.CreateNewsLetterRequest)
	return req, nil
}

func decodeAddSchemeRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*newsletterpb.NewsScheme)
	return req, nil
}
