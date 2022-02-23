package admin

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	errResourceNotFound = "error resource not found"
)

type DB interface {
	getUser(ctx context.Context, req SignInRequest) (*AdminModel, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) DB {
	return &repository{
		db: db,
	}
}

func (r repository) getUser(ctx context.Context, req SignInRequest) (*AdminModel, error) {
	var res AdminModel
	tx := r.db.WithContext(ctx).Model(&AdminModel{}).Where(&AdminModel{Email: req.Email}).Take(&res)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, errors.Wrap(tx.Error, errResourceNotFound)
		}

	}
	return &res, nil
}
