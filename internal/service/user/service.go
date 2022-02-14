package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sankethkini/NewsLetter-Backend/pkg/auth"
	"github.com/sankethkini/NewsLetter-Backend/pkg/encryption"
	"github.com/sankethkini/NewsLetter-Backend/pkg/role"
	userpb "github.com/sankethkini/NewsLetter-Backend/proto/userpb/v1"
)

const (
	errEmailPasswordNotMatch = "email or password is inncorect"
)

type Service interface {
	CreateUser(ctx context.Context, usr *userpb.CreateUserRequest) (*userpb.User, error)
	ValidateUser(ctx context.Context, sgn *userpb.ValidateUserRequest) (*userpb.ValidateUserResponse, error)
	GetEmail(ctx context.Context, ID *userpb.GetEmailRequest) (*userpb.Email, error)
}

type service struct {
	repo       DB
	jwtManager *auth.JWTManager
}

func NewUserService(repo DB, man *auth.JWTManager) Service {
	u := service{
		repo:       repo,
		jwtManager: man,
	}
	return &u
}

func (u *service) CreateUser(ctx context.Context, usr *userpb.CreateUserRequest) (*userpb.User, error) {
	m := ProtoToModel(usr.User)
	m.UserID = uuid.NewString()
	m.Password = encryption.Encrypt(m.Password)
	ret, err := u.repo.insertUser(ctx, &m)
	if err != nil {
		return nil, err
	}

	val := ModelToProto(ret)
	return &val, err
}

func (u *service) ValidateUser(ctx context.Context, sgn *userpb.ValidateUserRequest) (*userpb.ValidateUserResponse, error) {
	b, err := u.repo.validate(ctx, SignInRequest{Email: sgn.Email, Password: sgn.Password})
	if err != nil {
		return nil, err
	}

	val := encryption.Compare(sgn.Password, []byte(b.Password))
	if val {
		token, err := u.jwtManager.Generator(b.Email, role.USER)
		if err != nil {
			return nil, err
		}
		return &userpb.ValidateUserResponse{UserId: b.UserID, Token: token, Email: b.Email, Name: b.Name}, nil
	}

	return nil, errors.New(errEmailPasswordNotMatch)
}

func (u *service) GetEmail(ctx context.Context, usrID *userpb.GetEmailRequest) (*userpb.Email, error) {
	email, err := u.repo.getEmail(ctx, GetEmailRequest{ID: usrID.Name})
	if err != nil {
		return nil, err
	}

	return &userpb.Email{Email: email}, err
}
