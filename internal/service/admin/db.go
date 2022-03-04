package admin

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sankethkini/NewsLetter-Backend/pkg/apperrors"
	"gorm.io/gorm"
)

const (
	errResourceNotFound = "database: error resource not found %s"
)

//go:generate mockgen -destination db_mock.go -package admin github.com/sankethkini/NewsLetter-Backend/internal/service/admin DB
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
			return nil, apperrors.E(ctx, errors.Wrapf(tx.Error, errResourceNotFound, req.Email))
		}
		return nil, apperrors.E(ctx, tx.Error)
	}
	return &res, nil
}
