package newsletter

import (
	"context"

	"github.com/google/uuid"
	"github.com/sankethkini/NewsLetter-Backend/internal/kproducer"
	newsletterpb "github.com/sankethkini/NewsLetter-Backend/proto/newsletterpb/v1"
)

type Service interface {
	CreateNewsLetter(ctx context.Context, req *newsletterpb.CreateNewsLetterRequest) (*newsletterpb.NewsLetter, error)
	AddSchemeToNews(ctx context.Context, req *newsletterpb.NewsScheme) (*newsletterpb.NewsScheme, error)
}

type service struct {
	repo DB
	kaf  kproducer.Producer
}

func NewNewsService(repo DB, kaf kproducer.Producer) Service {
	return &service{repo: repo, kaf: kaf}
}

func (svc *service) CreateNewsLetter(ctx context.Context, req *newsletterpb.CreateNewsLetterRequest) (*newsletterpb.NewsLetter, error) {
	var dbreq = NewsLetterModel{NewsLetterID: uuid.NewString(), Title: req.Title, Body: req.Body}
	if err := dbreq.validate(); err != nil {
		return nil, err
	}

	resp, err := svc.repo.addNewsLetter(ctx, &dbreq)
	if err != nil {
		return nil, err
	}
	return ModelToProto(resp), nil
}

// nolint:govet
func (svc *service) AddSchemeToNews(ctx context.Context, req *newsletterpb.NewsScheme) (*newsletterpb.NewsScheme, error) {
	dbreq := AddSchemeRequest{NewsLetterID: req.NewsLetterId, SchemeID: req.SchemeId}
	if err := dbreq.validate(); err != nil {
		return nil, err
	}

	resp, err := svc.repo.addSchemeToNews(ctx, dbreq)
	if err != nil {
		return nil, err
	}

	news, err := svc.repo.getNewsLetter(ctx, req.NewsLetterId)
	if err != nil {
		return nil, err
	}

	r1 := ModelToProto(news)
	res := SchemeToProto(resp)

	data := EmailData{Letter: *r1, Scheme: res}
	msg, err := tojson(data)
	if err != nil {
		return nil, err
	}
	err = svc.kaf.Produce(ctx, []byte(req.NewsLetterId), []byte(msg))
	if err != nil {
		return nil, err
	}
	return &res, nil
}
