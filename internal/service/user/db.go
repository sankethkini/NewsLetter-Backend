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
	validate(ctx context.Context, s SignInRequest) (*UserModel, error)
	getEmail(ctx context.Context, g GetEmailRequest) (string, error)
}

type repository struct {
	db *gorm.DB
}

func NewDB(db *gorm.DB) DB {
	return &repository{db: db}
}

func (r repository) insertUser(ctx context.Context, m *UserModel) (*UserModel, error) {
	err := checkForTable(r.db)
	if err != nil {
		return nil, errors.Wrap(err, errTableCreation)
	}

	err = m.validate()
	if err != nil {
		return nil, errors.Wrap(err, errUserCreation)
	}

	usr := UserModel{Email: m.Email}
	exists, err := checkIfRecordExists(r.db, &usr)
	if err != nil {
		return nil, errors.Wrap(err, errGetEmail)
	}

	if exists {
		return nil, errors.New(errRecordExists)
	}

	tx := r.db.WithContext(ctx).Create(m)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, errUserCreation)
	}

	var ret UserModel
	tx = r.db.Find(&ret)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, errUserCreation)
	}

	return &ret, nil
}

func (r repository) validate(ctx context.Context, s SignInRequest) (*UserModel, error) {
	err := checkForTable(r.db)
	if err != nil {
		return nil, errors.Wrap(err, errTableCreation)
	}

	err = s.validate()
	if err != nil {
		return nil, err
	}

	usr := UserModel{Email: s.Email}
	exists, err := checkIfRecordExists(r.db, &usr)
	if err != nil {
		return nil, errors.Wrap(err, errGetEmail)
	}

	if !exists {
		return nil, errors.New(errRecordNotExists)
	}

	var record *UserModel
	err = r.db.WithContext(ctx).Where("email = ?", s.Email).Limit(1).Find(&record).Error
	if err != nil {
		return nil, errors.Wrap(err, errSignin)
	}

	return record, nil
}

func (r repository) getEmail(ctx context.Context, g GetEmailRequest) (string, error) {
	err := checkForTable(r.db)
	if err != nil {
		return "", errors.Wrap(err, errTableCreation)
	}

	err = g.validate()
	if err != nil {
		return "", err
	}

	usr := UserModel{UserID: g.ID}
	exists, err := checkIfRecordExists(r.db, &usr)
	if err != nil {
		return "", errors.Wrap(err, errGetEmail)
	}

	if !exists {
		return "", errors.New(errRecordNotExists)
	}

	var record *UserModel
	err = r.db.WithContext(ctx).Where("user_id = ?", g.ID).Find(&record).Error
	if err != nil {
		return "", errors.Wrap(err, errGetEmail)
	}

	return record.Email, nil
}

func checkForTable(db *gorm.DB) error {
	if !db.Migrator().HasTable(&UserModel{}) {
		err := db.Migrator().AutoMigrate(&UserModel{})
		return err
	}
	return nil
}

func checkIfRecordExists(db *gorm.DB, usr *UserModel) (bool, error) {
	// check if exists.
	count := int64(0)
	err := db.Model(&UserModel{}).Where(usr).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count != 0, nil
}
