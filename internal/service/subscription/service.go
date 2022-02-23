package subscription

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sankethkini/NewsLetter-Backend/internal/enum"
	"github.com/sankethkini/NewsLetter-Backend/pkg/cache"
	subscriptionpb "github.com/sankethkini/NewsLetter-Backend/proto/subscriptionpb/v1"
)

type Service interface {
	AddUser(context.Context, *subscriptionpb.AddUserRequest) (*subscriptionpb.AddUserResponse, error)
	RemoveUser(context.Context, *subscriptionpb.RemoveUserRequest) (*subscriptionpb.RemoveUserResponse, error)
	CreateScheme(context.Context, *subscriptionpb.CreateSchemeRequest) (*subscriptionpb.Scheme, error)
	Renew(context.Context, *subscriptionpb.RenewRequest) (*subscriptionpb.RenewResponse, error)
	Search(context.Context, *subscriptionpb.SearchRequest) (*subscriptionpb.SearchResponse, error)
	Sort(context.Context, *subscriptionpb.SortRequest) (*subscriptionpb.SortResponse, error)
	Filter(context.Context, *subscriptionpb.FilterRequest) (*subscriptionpb.FilterResponse, error)
	GetUsers(context.Context, *subscriptionpb.GetUsersRequest) (*subscriptionpb.GetUsersResponse, error)
}

type service struct {
	repo  DB
	redis cache.Service
}

func NewSubService(repo DB, redis cache.Service) Service {
	return &service{repo: repo, redis: redis}
}

// nolint: gosec
func (svc *service) AddUser(ctx context.Context, req *subscriptionpb.AddUserRequest) (*subscriptionpb.AddUserResponse, error) {

	dbreq := AddUserRequest{UserID: req.UserId, SchemeID: req.SchemeId}

	sub, err := svc.repo.getSubscription(ctx, req.SchemeId)
	if err != nil {
		return nil, err
	}

	val := time.Now().AddDate(0, 0, sub.Days)
	dbreq.Validity = val
	resp, err := svc.repo.addUser(ctx, dbreq)
	if err != nil {
		return nil, err
	}

	res := make([]*subscriptionpb.UserSubscription, 0, len(resp))
	for _, val := range resp {
		res = append(res, UserSubModelToProto(&val))
	}

	res1 := subscriptionpb.AddUserResponse{Subs: res}
	return &res1, nil
}

// nolint: gosec
func (svc *service) RemoveUser(ctx context.Context, req *subscriptionpb.RemoveUserRequest) (*subscriptionpb.RemoveUserResponse, error) {

	dbreq := UserSchemeRequest{UserID: req.UserId, SchemeID: req.SchemeId}

	if err := dbreq.validate(); err != nil {
		return nil, err
	}

	resp, err := svc.repo.removeUser(ctx, dbreq)
	if err != nil {
		return nil, err
	}

	res := make([]*subscriptionpb.UserSubscription, 0, len(resp))

	for _, val := range resp {
		res = append(res, UserSubModelToProto(&val))
	}

	res1 := subscriptionpb.RemoveUserResponse{Subs: res}
	return &res1, nil
}

func (svc *service) CreateScheme(ctx context.Context, req *subscriptionpb.CreateSchemeRequest) (*subscriptionpb.Scheme, error) {
	mod := SubscriptionModel{SchemeID: uuid.NewString(), Name: req.Name, Days: int(req.Days), Price: float64(req.Price)}
	resp, err := svc.repo.createScheme(ctx, &mod)
	if err != nil {
		return nil, err
	}

	scm := SubModelToProto(resp)
	return scm, nil
}

func (svc *service) Renew(ctx context.Context, req *subscriptionpb.RenewRequest) (*subscriptionpb.RenewResponse, error) {
	dbreq := UserSchemeRequest{UserID: req.UserId, SchemeID: req.SchemeId}

	if err := dbreq.validate(); err != nil {
		return nil, err
	}
	mod, err := svc.repo.getUserScheme(ctx, dbreq)
	if err != nil {
		return nil, err
	}

	sub, err := svc.repo.getSubscription(ctx, req.SchemeId)
	if err != nil {
		return nil, err
	}

	usrTime := mod.Validity
	curTime := time.Now()
	var val time.Time
	if usrTime.Sub(curTime) <= 0 {
		val = time.Now().AddDate(0, 0, sub.Days)
	} else {
		val = usrTime.AddDate(0, 0, sub.Days)
	}

	resp, err := svc.repo.renew(ctx, dbreq, val)
	if err != nil {
		return nil, err
	}

	res := UserSubModelToProto(resp)
	return &subscriptionpb.RenewResponse{Sub: res}, nil
}

// nolint: gosec
func (svc *service) Search(ctx context.Context, req *subscriptionpb.SearchRequest) (*subscriptionpb.SearchResponse, error) {
	data, err := svc.redis.Get(ctx, "search"+req.String())

	if err == nil {
		return &subscriptionpb.SearchResponse{Subs: data}, nil
	}

	resp, err := svc.repo.search(ctx, SearchRequest{req.Text})
	if err != nil {
		return nil, err
	}

	res := make([]*subscriptionpb.Scheme, 0, len(resp))

	for _, val := range resp {
		res = append(res, SubModelToProto(&val))
	}

	svc.redis.Set(ctx, "search"+req.String(), res)

	return &subscriptionpb.SearchResponse{Subs: res}, nil
}

// nolint:gosec
func (svc *service) Sort(ctx context.Context, req *subscriptionpb.SortRequest) (*subscriptionpb.SortResponse, error) {
	data, err := svc.redis.Get(ctx, "sort"+req.String())

	if err == nil {
		return &subscriptionpb.SortResponse{Subs: data}, nil
	}

	resp, err := svc.repo.sort(ctx, req.Field.String())
	if err != nil {
		return nil, err
	}
	res := make([]*subscriptionpb.Scheme, 0, len(resp))

	for _, val := range resp {
		res = append(res, SubModelToProto(&val))
	}

	svc.redis.Set(ctx, "sort"+req.String(), res)

	return &subscriptionpb.SortResponse{Subs: res}, nil
}

// nolint: gosec
func (svc *service) Filter(ctx context.Context, req *subscriptionpb.FilterRequest) (*subscriptionpb.FilterResponse, error) {
	data, err := svc.redis.Get(ctx, "filter"+req.String())
	if err == nil {
		return &subscriptionpb.FilterResponse{Subs: data}, nil
	}

	var mod FilterRequest
	mod.field = enum.Field(req.Field)
	mod.min = req.Min
	mod.max = req.Max
	resp, err := svc.repo.filter(ctx, mod)
	if err != nil {
		return nil, err
	}

	res := make([]*subscriptionpb.Scheme, 0, len(resp))
	for _, val := range resp {
		res = append(res, SubModelToProto(&val))
	}

	svc.redis.Set(ctx, "filter"+req.String(), res)

	return &subscriptionpb.FilterResponse{Subs: res}, nil
}

func (svc *service) GetUsers(ctx context.Context, req *subscriptionpb.GetUsersRequest) (*subscriptionpb.GetUsersResponse, error) {
	resp, err := svc.repo.getUsers(ctx, req.SchemeId)
	if err != nil {
		return nil, err
	}
	return &subscriptionpb.GetUsersResponse{UserIds: resp}, nil
}
