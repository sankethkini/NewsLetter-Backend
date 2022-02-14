package admin

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

const (
	errTableCreation = "database: cannot create the table"
)

type DB interface {
	signIn(ctx context.Context, req SignInRequest) (*AdminModel, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) DB {
	return &repository{
		db: db,
	}
}

func (r repository) signIn(ctx context.Context, req SignInRequest) (*AdminModel, error) {
	err := checkForTable(r.db)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, errTableCreation)
	}

	err = req.validate()
	if err != nil {
		return nil, err
	}

	mod := AdminModel{Email: req.Email}
	ok, err := CheckIfRecordExists(ctx, r.db, &mod)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, status.Errorf(codes.NotFound, "record not found")
	}
	var res AdminModel
	tx := r.db.WithContext(ctx).Where("email=?", req.Email).Find(&res)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &res, nil
}

func checkForTable(db *gorm.DB) error {
	if !db.Migrator().HasTable(&AdminModel{}) {
		err := db.Migrator().AutoMigrate(&AdminModel{})
		return err
	}
	return nil
}

func CheckIfRecordExists(ctx context.Context, db *gorm.DB, req *AdminModel) (bool, error) {
	count := int64(0)
	err := db.Model(&AdminModel{}).Where(req).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count != 0, nil
}
