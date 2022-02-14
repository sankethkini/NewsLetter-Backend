package newsletter

import (
	"context"

	"github.com/sankethkini/NewsLetter-Backend/internal/service/subscription"
	"gorm.io/gorm"
)

type DB interface {
	addNewsLetter(context.Context, *NewsLetterModel) (*NewsLetterModel, error)
	addSchemeToNews(context.Context, AddSchemeRequest) (*NewsSchemes, error)
	getNewsLetter(context.Context, string) (NewsLetterModel, error)
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
	err := checkForTable(r.db)
	if err != nil {
		return nil, err
	}

	err = req.validate()
	if err != nil {
		return nil, err
	}

	tx := r.db.WithContext(ctx).Create(req)
	if tx.Error != nil {
		return nil, tx.Error
	}

	var ret NewsLetterModel
	ret.NewsLetterID = req.NewsLetterID
	tx = r.db.WithContext(ctx).Find(&ret)
	if tx.Error != nil {
		return nil, err
	}

	return &ret, nil
}

func (r repository) addSchemeToNews(ctx context.Context, req AddSchemeRequest) (*NewsSchemes, error) {
	err := checkForTable(r.db)
	if err != nil {
		return nil, err
	}
	err = req.validate()
	if err != nil {
		return nil, err
	}

	ok, err := CheckIfRecordExists(ctx, r.db, req.NewsLetterID)
	if !ok {
		return nil, err
	}

	ok, err = subscription.CheckIfRecordExists(ctx, r.db, &subscription.SubscriptionModel{SchemeID: req.SchemeID})
	if !ok {
		return nil, err
	}

	var mod NewsSchemes
	mod.NewsLetterID = req.NewsLetterID
	mod.SchemeID = req.SchemeID
	tx := r.db.WithContext(ctx).Create(&mod)

	if tx.Error != nil {
		return nil, tx.Error
	}

	var resp NewsSchemes
	tx = r.db.WithContext(ctx).Where("news_letter_id =? and scheme_id=?", req.NewsLetterID, req.SchemeID).Find(&resp)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &resp, nil
}

func (r repository) getNewsLetter(ctx context.Context, id string) (NewsLetterModel, error) {
	err := checkForTable(r.db)
	if err != nil {
		return NewsLetterModel{}, err
	}

	var data NewsLetterModel
	tx := r.db.WithContext(ctx).Model(&NewsLetterModel{}).Where("news_letter_id=?", id).Find(&data)

	if tx.Error != nil {
		return NewsLetterModel{}, err
	}

	return data, nil
}

func checkForTable(db *gorm.DB) error {
	if !db.Migrator().HasTable(&NewsLetterModel{}) {
		err := db.Migrator().AutoMigrate(&NewsLetterModel{})
		return err
	}
	if !db.Migrator().HasTable(&NewsSchemes{}) {
		err := db.Migrator().AutoMigrate(&NewsSchemes{})
		return err
	}
	return nil
}

func CheckIfRecordExists(ctx context.Context, db *gorm.DB, newsID string) (bool, error) {
	count := int64(0)
	err := db.WithContext(ctx).Model(&NewsLetterModel{}).Where(NewsLetterModel{NewsLetterID: newsID}).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count != 0, nil
}
