package subscription

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sankethkini/NewsLetter-Backend/pkg/cache"
	v1 "github.com/sankethkini/NewsLetter-Backend/proto/subscriptionpb/v1"
	"github.com/stretchr/testify/require"
)

func TestAddUserSuccess(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockDB := NewMockDB(ctrl)
	mockRedis := cache.NewMockService(ctrl)
	mockDB.EXPECT().getSubscription(gomock.Any(), gomock.Any()).Return(&SubscriptionModel{
		Name:     "some name",
		SchemeID: "someid",
		Price:    123,
		Days:     12,
	}, nil).Times(1)
	mockDB.EXPECT().addUser(gomock.Any(), gomock.Any()).Return([]UserSubscription{
		{
			SchemeID: "someid123",
			UserID:   "someid123",
			Validity: time.Now(),
		},
	}, nil).Times(1)
	svc := NewSubService(mockDB, mockRedis)
	resp, err := svc.AddUser(ctx, &v1.AddUserRequest{UserId: "someid123", SchemeId: "someid123"})
	require.NotNil(t, resp)
	require.Nil(t, err)
}

func TestRemoveUserSuccess(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockDB := NewMockDB(ctrl)
	mockRedis := cache.NewMockService(ctrl)
	mockDB.EXPECT().removeUser(gomock.Any(), gomock.Any()).Return([]UserSubscription{
		{
			SchemeID: "someid123",
			UserID:   "someid123",
			Validity: time.Now(),
		},
	}, nil).Times(1)
	svc := NewSubService(mockDB, mockRedis)
	resp, err := svc.RemoveUser(ctx, &v1.RemoveUserRequest{UserId: "some123", SchemeId: "some123"})
	require.NotNil(t, resp)
	require.Nil(t, err)
}

func TestCreateSchemeSuccess(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockDB := NewMockDB(ctrl)
	mockRedis := cache.NewMockService(ctrl)
	mockDB.EXPECT().createScheme(gomock.Any(), gomock.Any()).Return(&SubscriptionModel{
		SchemeID: "someid123",
		Name:     "some name",
		Price:    123,
		Days:     28,
	}, nil).Times(1)
	svc := NewSubService(mockDB, mockRedis)
	resp, err := svc.CreateScheme(ctx, &v1.CreateSchemeRequest{
		Name:  "some name",
		Price: 123,
		Days:  28,
	})
	require.NotNil(t, resp)
	require.Nil(t, err)
}

func TestRenewSuccess(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockDB := NewMockDB(ctrl)
	mockRedis := cache.NewMockService(ctrl)
	mockDB.EXPECT().getUserScheme(gomock.Any(), gomock.Any()).Return(&UserSubscription{
		SchemeID: "some123",
		UserID:   "some123",
		Validity: time.Now(),
	}, nil).Times(1)
	mockDB.EXPECT().getSubscription(gomock.Any(), gomock.Any()).Return(&SubscriptionModel{
		Name:     "some name",
		SchemeID: "some123",
		Price:    123,
		Days:     12,
	}, nil).Times(1)
	mockDB.EXPECT().renew(gomock.Any(), gomock.Any(), gomock.Any()).Return(&UserSubscription{
		SchemeID: "some123",
		UserID:   "some123",
		Validity: time.Now(),
	}, nil).Times(1)
	svc := NewSubService(mockDB, mockRedis)
	resp, err := svc.Renew(ctx, &v1.RenewRequest{SchemeId: "some123", UserId: "some123"})
	require.NotNil(t, resp)
	require.Nil(t, err)
}

func TestSearchSuccess(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockDB := NewMockDB(ctrl)
	mockRedis := cache.NewMockService(ctrl)
	mockDB.EXPECT().search(gomock.Any(), gomock.Any()).Return([]SubscriptionModel{
		{
			SchemeID: "some123",
			Name:     "some name",
			Price:    123,
			Days:     28,
		},
	}, nil).Times(1)
	mockRedis.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("not found")).Times(1)
	mockRedis.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return().Times(1)
	svc := NewSubService(mockDB, mockRedis)
	resp, err := svc.Search(ctx, &v1.SearchRequest{Text: "some"})
	require.NotNil(t, resp)
	require.Nil(t, err)
}

func TestSortSuccess(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockDB := NewMockDB(ctrl)
	mockRedis := cache.NewMockService(ctrl)
	mockDB.EXPECT().sort(gomock.Any(), gomock.Any()).Return([]SubscriptionModel{
		{
			SchemeID: "some123",
			Name:     "some name",
			Price:    123,
			Days:     28,
		},
	}, nil).Times(1)
	mockRedis.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("not found")).Times(1)
	mockRedis.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return().Times(1)
	svc := NewSubService(mockDB, mockRedis)
	resp, err := svc.Sort(ctx, &v1.SortRequest{Field: v1.Field_DAYS})
	require.NotNil(t, resp)
	require.Nil(t, err)
}

func TestFilterSuccess(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockDB := NewMockDB(ctrl)
	mockRedis := cache.NewMockService(ctrl)
	mockDB.EXPECT().filter(gomock.Any(), gomock.Any()).Return([]SubscriptionModel{
		{
			SchemeID: "some123",
			Name:     "some name",
			Price:    123,
			Days:     28,
		},
	}, nil).Times(1)
	mockRedis.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("not found")).Times(1)
	mockRedis.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return().Times(1)
	svc := NewSubService(mockDB, mockRedis)
	resp, err := svc.Filter(ctx, &v1.FilterRequest{Field: v1.Field_DAYS, Min: 12, Max: 30})
	require.NotNil(t, resp)
	require.Nil(t, err)
}

func TestGetUserSuccess(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockDB := NewMockDB(ctrl)
	mockRedis := cache.NewMockService(ctrl)
	mockDB.EXPECT().getUsers(gomock.Any(), gomock.Any()).Return([]string{"one", "two"}, nil)
	svc := NewSubService(mockDB, mockRedis)
	resp, err := svc.GetUsers(ctx, &v1.GetUsersRequest{SchemeId: "1234"})
	require.NotNil(t, resp)
	require.Nil(t, err)
}
