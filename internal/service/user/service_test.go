package user

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sankethkini/NewsLetter-Backend/pkg/auth"
	"github.com/sankethkini/NewsLetter-Backend/pkg/encryption"
	v1 "github.com/sankethkini/NewsLetter-Backend/proto/userpb/v1"
	"github.com/stretchr/testify/require"
)

func TestCreateUserSuccess(t *testing.T) {
	ctx := context.Background()
	t.Logf("Name:%s", t.Name())
	ctrl := gomock.NewController(t)
	mockDB := NewMockDB(ctrl)
	mockJWT := auth.NewMockJWTManager(ctrl)
	mockDB.EXPECT().insertUser(gomock.Any(), gomock.Any()).Return(
		&UserModel{UserID: "someid123", Email: "someuse@some.com", Name: "some"}, nil).Times(1)
	svc := NewUserService(mockDB, mockJWT)
	resp, err := svc.CreateUser(ctx, &v1.CreateUserRequest{User: &v1.User{Email: "someuser@some.com", Name: "some", Password: "some password"}})
	require.Nil(t, err)
	require.NotNil(t, resp)
}

func TestValidateUserSuccess(t *testing.T) {
	ctx := context.Background()
	t.Logf("Name:%s", t.Name())
	ctrl := gomock.NewController(t)
	mockDB := NewMockDB(ctrl)
	mockJWT := auth.NewMockJWTManager(ctrl)
	mockDB.EXPECT().getUser(gomock.Any(), gomock.Any()).Return(
		&UserModel{Email: "someuser@some.com", Password: encryption.Encrypt("some password")}, nil).Times(1)
	mockJWT.EXPECT().Generator(gomock.Any(), gomock.Any()).Return("some token123", nil).Times(1)
	svc := NewUserService(mockDB, mockJWT)
	resp, err := svc.ValidateUser(ctx, &v1.ValidateUserRequest{Email: "someuser@some.com", Password: "some password"})
	require.Nil(t, err)
	require.NotNil(t, resp)
}

func TestGetEmailSuccess(t *testing.T) {
	ctx := context.Background()
	t.Logf("Name:%s", t.Name())
	ctrl := gomock.NewController(t)
	mockDB := NewMockDB(ctrl)
	mockJWT := auth.NewMockJWTManager(ctrl)
	mockDB.EXPECT().getEmail(gomock.Any(), gomock.Any()).Return("someuser@some.com", nil).Times(1)
	svc := NewUserService(mockDB, mockJWT)
	resp, err := svc.GetEmail(ctx, &v1.GetEmailRequest{Name: "userid123"})
	require.Nil(t, err)
	require.NotNil(t, resp)
}
