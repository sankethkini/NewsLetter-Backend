package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/sankethkini/NewsLetter-Backend/internal/enum"
	"github.com/sankethkini/NewsLetter-Backend/pkg/apperrors"
	"github.com/sankethkini/NewsLetter-Backend/pkg/auth"
	"github.com/sankethkini/NewsLetter-Backend/pkg/encryption"
	userpb "github.com/sankethkini/NewsLetter-Backend/proto/userpb/v1"
)

const (
	errEmailPasswordNotMatch = "email or password is inncorect"
	errInputValues           = "error in input values"
	errTokenGen              = "error in generating token"
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
	if err := m.validate(); err != nil {
		return nil, apperrors.E(ctx, err, errInputValues)
	}

	ret, err := u.repo.insertUser(ctx, &m)
	if err != nil {
		return nil, err
	}

	val := ModelToProto(ret)
	return &val, err
}

func (u *service) ValidateUser(ctx context.Context, sgn *userpb.ValidateUserRequest) (*userpb.ValidateUserResponse, error) {
	dbreq := SignInRequest{Email: sgn.Email, Password: sgn.Password}
	if err := dbreq.validate(); err != nil {
		return nil, apperrors.E(ctx, err, errInputValues)
	}

	b, err := u.repo.getUser(ctx, dbreq)
	if err != nil {
		return nil, err
	}

	// compare encrypted password and user provided password.
	val := encryption.Compare(sgn.Password, []byte(b.Password))

	// generate user token.
	if val {
		token, err := u.jwtManager.Generator(b.Email, enum.USER)
		if err != nil {
			return nil, apperrors.E(ctx, err, errTokenGen)
		}
		return &userpb.ValidateUserResponse{UserId: b.UserID, Token: token, Email: b.Email, Name: b.Name}, nil
	}

	return nil, apperrors.E(ctx, errEmailPasswordNotMatch)
}

func (u *service) GetEmail(ctx context.Context, usrID *userpb.GetEmailRequest) (*userpb.Email, error) {
	dbreq := GetEmailRequest{ID: usrID.Name}
	if err := dbreq.validate(); err != nil {
		return nil, apperrors.E(ctx, err, errInputValues)
	}

	email, err := u.repo.getEmail(ctx, dbreq)
	if err != nil {
		return nil, err
	}

	return &userpb.Email{Email: email}, err
}
