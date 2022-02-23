package admin

import (
	"context"

	"github.com/sankethkini/NewsLetter-Backend/internal/enum"
	"github.com/sankethkini/NewsLetter-Backend/pkg/auth"
	"github.com/sankethkini/NewsLetter-Backend/pkg/encryption"
	adminpb "github.com/sankethkini/NewsLetter-Backend/proto/adminpb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	SingIn(context.Context, *adminpb.SignInRequest) (*adminpb.SignInResponse, error)
}

type service struct {
	repo       DB
	jwtManager *auth.JWTManager
}

func NewAdminService(repo DB, jwt *auth.JWTManager) Service {
	return &service{
		repo:       repo,
		jwtManager: jwt,
	}
}

func (adm *service) SingIn(ctx context.Context, req *adminpb.SignInRequest) (*adminpb.SignInResponse, error) {
	dbreq := SignInRequest{Email: req.Email, Password: req.Password}

	err := dbreq.validate()
	if err != nil {
		return nil, err
	}

	resp, err := adm.repo.getUser(ctx, dbreq)
	if err != nil {
		return nil, err
	}

	ok := encryption.Compare(req.Password, []byte(resp.Password))
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "email and password not matching")
	}

	token, err := adm.jwtManager.Generator(req.Email, enum.ADMIN)
	if err != nil {
		return nil, err
	}
	return &adminpb.SignInResponse{AdminId: resp.AdminID, Token: token}, nil
}
