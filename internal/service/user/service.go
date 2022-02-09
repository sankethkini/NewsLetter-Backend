package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sankethkini/NewsLetter-Backend/pkg/encryption"
	"github.com/sankethkini/NewsLetter-Backend/proto/userpb/v1"
)

const (
	errEmailPasswordNotMatch = "email or password is inncorect"
)

type UserService interface {
	CreateUser(ctx context.Context, usr *userpb.CreateUserRequest) (*userpb.User, error)
	ValidateUser(ctx context.Context, sgn *userpb.ValidateUserRequest) (*userpb.ValidateUserResponse, error)
	GetEmail(ctx context.Context, ID *userpb.GetEmailRequest) (*userpb.Email, error)
}

type UserServiceImpl struct {
	repo DB
}

func NewUserService() *UserServiceImpl {
	//need wire
	repo := NewDB()
	u := UserServiceImpl{
		repo: repo,
	}
	return &u
}

func (u *UserServiceImpl) CreateUser(ctx context.Context, usr *userpb.CreateUserRequest) (*userpb.User, error) {
	m := ProtoToModel(usr.User)
	fmt.Println(m, u)
	m.UserID = uuid.NewString()
	m.Password = encryption.Encrypt(m.Password)
	ret, err := u.repo.insertUser(ctx, &m)
	val := ModelToProto(ret)
	fmt.Println(m.UserID)
	return &val, err
}

func (u *UserServiceImpl) ValidateUser(ctx context.Context, sgn *userpb.ValidateUserRequest) (*userpb.ValidateUserResponse, error) {
	fmt.Println(3, sgn.Email)
	b, id, err := u.repo.signIn(ctx, SignInRequest{Email: sgn.Email, Password: sgn.Password})
	if err != nil {
		return nil, err
	}
	// val := encryption.Compare(sgn.Password, []byte(b.Password))
	// if val {
	m := ModelToProto(b)

	return &userpb.ValidateUserResponse{UserId: id, Token: "some", Email: m.Email, Name: m.Name}, nil
	//}
	//return nil, errors.Wrap(nil, errEmailPasswordNotMatch)
}

func (u *UserServiceImpl) GetEmail(ctx context.Context, usrID *userpb.GetEmailRequest) (*userpb.Email, error) {
	email, err := u.repo.getEmail(ctx, GetEmailRequest{ID: usrID.Id})
	return &userpb.Email{Email: email}, err
}
