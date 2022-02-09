package user

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sankethkini/NewsLetter-Backend/pkg/database"
	"gorm.io/gorm"
)

const (
	errResourceCreation = "database: failed to batch create inventoryItems"
	errSignin           = "database: failed in finding record during sign in"
	errGetEmail         = "database: failed in get email"
	errTableCreation    = "database: failed to create table"
)

type DB interface {
	insertUser(ctx context.Context, m *UserModel) (*UserModel, error)
	signIn(ctx context.Context, s SignInRequest) (*UserModel, string, error)
	getEmail(ctx context.Context, g GetEmailRequest) (string, error)
}

type repository struct {
	db *gorm.DB
}

//TODO if email already exists
func (r repository) insertUser(ctx context.Context, m *UserModel) (*UserModel, error) {
	err := checkForTable(r.db)
	if err != nil {
		return nil, errors.Wrap(err, errTableCreation)
	}

	err = m.validate()
	if err != nil {
		return nil, err
	}
	tx := r.db.WithContext(ctx).Create(m)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, errResourceCreation)
	}
	var ret UserModel
	tx = r.db.Find(&ret)

	return &ret, nil
}

func (r repository) signIn(ctx context.Context, s SignInRequest) (*UserModel, string, error) {
	err := checkForTable(r.db)
	if err != nil {
		return nil, "", errors.Wrap(err, errTableCreation)
	}

	err = s.validate()
	if err != nil {
		return nil, "", errors.Wrap(err, errSignin)
	}

	var record *UserModel
	tx := r.db.WithContext(ctx).Where("email = ?", s.Email).Find(&record)
	if tx.Error != nil {
		return nil, "", errors.Wrap(tx.Error, errSignin)
	}

	return record, record.UserID, nil
}

func (r repository) getEmail(ctx context.Context, g GetEmailRequest) (string, error) {

	err := checkForTable(r.db)
	if err != nil {
		return "", errors.Wrap(err, errTableCreation)
	}

	err = g.validate()
	if err != nil {
		return "", errors.Wrap(err, errSignin)
	}
	var record *UserModel
	tx := r.db.WithContext(ctx).Where("user_id = ?", g.ID).Find(&record)
	if tx.Error != nil {
		return "", errors.Wrap(tx.Error, errGetEmail)
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

func NewDB() DB {
	//wire
	op, err := database.Open()
	if err != nil {
		panic(err)
	}
	return &repository{db: op}
}
