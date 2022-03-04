package user

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sankethkini/NewsLetter-Backend/pkg/apperrors"
	"gorm.io/gorm"
)

const (
	errUserCreation = "database: failed to create user %s"
	errFind         = "database: error in finding %s"
)

//go:generate mockgen -destination db_mock.go -package user github.com/sankethkini/NewsLetter-Backend/internal/service/user DB
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
		return nil, apperrors.E(ctx, errors.Wrapf(tx.Error, errUserCreation, m.Email))
	}
	var ret UserModel
	tx = r.db.WithContext(ctx).Model(&UserModel{}).Where(&UserModel{UserID: m.UserID}).Take(&ret)
	if tx.Error != nil {
		return nil, apperrors.E(ctx, errors.Wrapf(tx.Error, errFind, m.Email))
	}

	return &ret, nil
}

func (r repository) getUser(ctx context.Context, s SignInRequest) (*UserModel, error) {
	var resp UserModel
	tx := r.db.WithContext(ctx).Model(&UserModel{}).Where(&UserModel{Email: s.Email}).Take(&resp)
	if tx.Error != nil {
		return nil, apperrors.E(ctx, errors.Wrapf(tx.Error, errFind, s.Email))
	}
	return &resp, nil
}

func (r repository) getEmail(ctx context.Context, g GetEmailRequest) (string, error) {
	usr := UserModel{UserID: g.ID}

	var record UserModel
	tx := r.db.WithContext(ctx).Model(&UserModel{}).Where(&usr).Take(&record)
	if tx.Error != nil {
		return "", apperrors.E(ctx, errors.Wrapf(tx.Error, errFind, g.ID))
	}

	return record.Email, nil
}
