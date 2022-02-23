package user

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	errUserCreation    = "database: failed to create user"
	errSignin          = "database: failed in finding record during sign in"
	errGetEmail        = "database: failed in get email"
	errTableCreation   = "database: failed to create table"
	errRecordExists    = "database: same email already exists"
	errRecordNotExists = "database: record not exists"
)

type DB interface {
	insertUser(ctx context.Context, m *UserModel) (*UserModel, error)
	getUser(ctx context.Context, s SignInRequest) (*UserModel, error)
	getEmail(ctx context.Context, g GetEmailRequest) (string, error)
}

type repository struct {
	db *gorm.DB
}

func NewDB(db *gorm.DB) DB {
	return &repository{db: db}
}

func (r repository) insertUser(ctx context.Context, m *UserModel) (*UserModel, error) {

	tx := r.db.WithContext(ctx).Create(m)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, errUserCreation)
	}
	var ret UserModel
	tx = r.db.WithContext(ctx).Model(&UserModel{}).Where(&UserModel{UserID: m.UserID}).Take(&ret)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, errUserCreation)
	}

	return &ret, nil
}

func (r repository) getUser(ctx context.Context, s SignInRequest) (*UserModel, error) {
	var resp UserModel
	tx := r.db.WithContext(ctx).Model(&UserModel{}).Where(&UserModel{Email: s.Email}).Take(&resp)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &resp, nil
}

func (r repository) getEmail(ctx context.Context, g GetEmailRequest) (string, error) {

	usr := UserModel{UserID: g.ID}

	var record UserModel
	tx := r.db.WithContext(ctx).Model(&UserModel{}).Where(&usr).Take(&record)
	if tx.Error != nil {
		return "", nil
	}

	return record.Email, nil
}
