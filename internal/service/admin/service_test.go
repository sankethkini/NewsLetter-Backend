package admin

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/sankethkini/NewsLetter-Backend/pkg/auth"
	"github.com/sankethkini/NewsLetter-Backend/pkg/encryption"
	v1 "github.com/sankethkini/NewsLetter-Backend/proto/adminpb/v1"
	"github.com/stretchr/testify/require"
)

func TestSignInSuccess(t *testing.T) {
	ctx := context.Background()
	t.Logf("Name:%s", t.Name())
	ctrl := gomock.NewController(t)
	mockDB := NewMockDB(ctrl)
	mockJWT := auth.NewMockJWTManager(ctrl)
	mockDB.EXPECT().getUser(gomock.Any(), gomock.Any()).Return(&AdminModel{Email: "someuser@some.com", Password: encryption.Encrypt("some password")}, nil).Times(1)
	mockJWT.EXPECT().Generator(gomock.Any(), gomock.Any()).Return("sometoken123", nil)
	svc := NewAdminService(mockDB, mockJWT)
	resp, err := svc.SingIn(ctx, &v1.SignInRequest{Email: "someuser@some.com", Password: "some password"})
	require.Nil(t, err)
	require.NotNil(t, resp)
}
