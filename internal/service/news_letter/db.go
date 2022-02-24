package newsletter

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sankethkini/NewsLetter-Backend/pkg/apperrors"
	"gorm.io/gorm"
)

const (
	errResourceNotFound = "database: error resource not found %s"
	errCreate           = "database: error in creating resource %s"
)

type DB interface {
	addNewsLetter(context.Context, *NewsLetterModel) (*NewsLetterModel, error)
	addSchemeToNews(context.Context, AddSchemeRequest) (*NewsSchemes, error)
	getNewsLetter(context.Context, string) (*NewsLetterModel, error)
}

type repository struct {
	db *gorm.DB
}

func NewNewsRepo(db *gorm.DB) DB {
	return &repository{
		db: db,
	}
}

func (r repository) addNewsLetter(ctx context.Context, req *NewsLetterModel) (*NewsLetterModel, error) {
	tx := r.db.WithContext(ctx).Create(&req)
	if tx.Error != nil {
		return nil, apperrors.E(ctx, errors.Wrapf(tx.Error, errCreate, req.NewsLetterID))
	}

	ret, err := r.getNewsLetter(ctx, req.NewsLetterID)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (r repository) addSchemeToNews(ctx context.Context, req AddSchemeRequest) (*NewsSchemes, error) {
	mod := NewsSchemes{NewsLetterID: req.NewsLetterID, SchemeID: req.SchemeID}
	tx := r.db.WithContext(ctx).Create(&mod)

	if tx.Error != nil {
		return nil, apperrors.E(ctx, errors.Wrapf(tx.Error, errCreate, req.NewsLetterID))
	}

	var resp NewsSchemes
	tx = r.db.WithContext(ctx).Model(&NewsSchemes{}).Where(&NewsSchemes{NewsLetterID: req.NewsLetterID, SchemeID: req.SchemeID}).Take(&resp)
	if tx.Error != nil {
		return nil, apperrors.E(ctx, errors.Wrapf(tx.Error, errResourceNotFound, req.NewsLetterID))
	}
	return &resp, nil
}

func (r repository) getNewsLetter(ctx context.Context, id string) (*NewsLetterModel, error) {
	var resp NewsLetterModel
	tx := r.db.WithContext(ctx).Model(&NewsLetterModel{}).Where(&NewsLetterModel{NewsLetterID: id}).Take(&resp)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, apperrors.E(ctx, tx.Error, errors.Wrapf(tx.Error, errResourceNotFound, id))
		}
		return nil, apperrors.E(ctx, tx.Error)
	}

	return &resp, nil
}
