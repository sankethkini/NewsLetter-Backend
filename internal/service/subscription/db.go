package subscription

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sankethkini/NewsLetter-Backend/internal/service/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type DB interface {
	addUser(context.Context, UserSchemeRequest) ([]UserSubscription, error)
	removeUser(context.Context, UserSchemeRequest) ([]UserSubscription, error)
	createScheme(context.Context, SchemeRequest) (*SubscriptionModel, error)
	renew(context.Context, UserSchemeRequest) (*UserSubscription, error)
	search(context.Context, string) ([]SubscriptionModel, error)
	sort(context.Context, Field) ([]SubscriptionModel, error)
	filter(context.Context, FilterRequest) ([]SubscriptionModel, error)
	getUsers(context.Context, string) ([]string, error)
}

type repository struct {
	db *gorm.DB
}

func NewSubRepo(db *gorm.DB) DB {
	return &repository{
		db: db,
	}
}

func (r repository) addUser(ctx context.Context, req UserSchemeRequest) ([]UserSubscription, error) {
	err := checkForTable(r.db)
	if err != nil {
		return nil, err
	}

	// check if such scheme exists.
	err = CheckIfSubsExists(ctx, r.db, req.UserID, req.SchemeID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "record not found")
	}

	ok, err := CheckForScheme(ctx, r.db, req.UserID, req.SchemeID)
	if err != nil || ok {
		return nil, status.Errorf(codes.AlreadyExists, "scheme already subscribed")
	}

	var us UserSubscription
	us.UserID = req.UserID
	us.SchemeID = req.SchemeID
	us.Validity = time.Now()

	tx := r.db.WithContext(ctx).Create(&us)

	if tx.Error != nil {
		return nil, tx.Error
	}

	ret, err := r.getSchemesOfUser(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (r repository) removeUser(ctx context.Context, req UserSchemeRequest) ([]UserSubscription, error) {
	err := checkForTable(r.db)
	if err != nil {
		return nil, err
	}

	err = req.validate()
	if err != nil {
		return nil, err
	}

	// check if such scheme exists.
	err = CheckIfSubsExists(ctx, r.db, req.UserID, req.SchemeID)
	if err != nil {
		return nil, err
	}

	ok, err := CheckForScheme(ctx, r.db, req.UserID, req.SchemeID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, status.Errorf(codes.NotFound, "record not found")
	}

	var us UserSubscription
	us.SchemeID = req.SchemeID
	us.UserID = req.UserID
	tx := r.db.Delete(&us)
	if tx.Error != nil {
		return nil, tx.Error
	}

	ret, err := r.getSchemesOfUser(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (r repository) renew(ctx context.Context, req UserSchemeRequest) (*UserSubscription, error) {
	err := checkForTable(r.db)
	if err != nil {
		return nil, err
	}

	err = req.validate()
	if err != nil {
		return nil, err
	}

	// check if such scheme exists.
	err = CheckIfSubsExists(ctx, r.db, req.UserID, req.SchemeID)
	if err != nil {
		return nil, err
	}

	ok, err := CheckForScheme(ctx, r.db, req.UserID, req.SchemeID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, status.Errorf(codes.NotFound, "record not found")
	}

	var us UserSubscription
	us.SchemeID = req.SchemeID
	us.UserID = req.UserID
	tx := r.db.WithContext(ctx).Find(&us)
	if tx.Error != nil {
		return nil, tx.Error
	}
	var sch SubscriptionModel
	sch.SchemeID = req.SchemeID
	tx = r.db.WithContext(ctx).Find(&sch)

	if tx.Error != nil {
		return nil, tx.Error
	}

	usrTime := us.Validity
	curTime := time.Now()
	if usrTime.Sub(curTime) <= 0 {
		us.Validity = time.Now().AddDate(0, 0, sch.Days)
	} else {
		us.Validity = usrTime.AddDate(0, 0, sch.Days)
	}

	us1 := UserSubscription{SchemeID: req.SchemeID, UserID: req.UserID}
	tx = r.db.WithContext(ctx).Model(&us1).Update("validity", us.Validity)
	if tx.Error != nil {
		return nil, tx.Error
	}

	tx = r.db.WithContext(ctx).Find(&us1)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &us1, nil
}

func (r repository) createScheme(ctx context.Context, s SchemeRequest) (*SubscriptionModel, error) {
	err := s.validate()
	if err != nil {
		return nil, err
	}

	var sc SubscriptionModel
	sc.Name = s.name
	sc.Days = s.days
	sc.Price = s.price
	id := uuid.NewString()
	sc.SchemeID = id

	tx := r.db.Create(&sc)
	if tx.Error != nil {
		return nil, tx.Error
	}

	var ret SubscriptionModel
	ret.SchemeID = id
	tx = r.db.Find(&ret)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return &ret, nil
}

// nolint: gocritic
func (r repository) getSchemesOfUser(ctx context.Context, UserID string) ([]UserSubscription, error) {
	exi, err := user.CheckIfRecordExists(r.db, &user.UserModel{UserID: UserID})
	if err != nil {
		return nil, err
	}
	if !exi {
		return nil, status.Errorf(codes.NotFound, "database: record with given id not found")
	}

	var subs []UserSubscription
	var schemeIDs []string
	tx := r.db.WithContext(ctx).Model(&UserSubscription{UserID: UserID}).Find(&subs)

	if tx.Error != nil {
		return nil, err
	}

	tx = r.db.WithContext(ctx).Find(&subs, schemeIDs)
	if tx.Error != nil {
		return nil, err
	}
	return subs, nil
}

func checkForTable(db *gorm.DB) error {
	if !db.Migrator().HasTable(&user.UserModel{}) {
		err := db.Migrator().AutoMigrate(&user.UserModel{})
		return err
	}
	if !db.Migrator().HasTable(&SubscriptionModel{}) {
		err := db.Migrator().AutoMigrate(&SubscriptionModel{})
		return err
	}
	if !db.Migrator().HasTable(&UserSubscription{}) {
		err := db.Migrator().AutoMigrate(&UserSchemeRequest{})
		if err != nil {
			return err
		}
	}
	return nil
}

func CheckIfRecordExists(ctx context.Context, db *gorm.DB, usr *SubscriptionModel) (bool, error) {
	// check if exists.
	count := int64(0)
	err := db.WithContext(ctx).Model(&SubscriptionModel{}).Where(usr).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

// nolint: gocritic
func CheckIfSubsExists(ctx context.Context, db *gorm.DB, UserID, SchemeID string) error {
	// check if such scheme exists.
	var s SubscriptionModel
	s.SchemeID = SchemeID
	exi, err := CheckIfRecordExists(ctx, db, &s)
	if err != nil {
		return err
	}
	if !exi {
		return status.Errorf(codes.NotFound, "database: record with given id not found")
	}

	var u user.UserModel
	u.UserID = UserID
	exi, err = user.CheckIfRecordExists(db, &u)
	if err != nil {
		return err
	}
	if !exi {
		return status.Errorf(codes.NotFound, "database: record with given id not found")
	}
	return nil
}

func CheckForScheme(ctx context.Context, db *gorm.DB, userID, schemeID string) (bool, error) {
	var s UserSubscription
	s.UserID = userID
	s.SchemeID = schemeID
	count := int64(0)
	err := db.WithContext(ctx).Model(&UserSubscription{}).Where(&s).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (r repository) search(ctx context.Context, req string) ([]SubscriptionModel, error) {
	err := checkForTable(r.db)
	if err != nil {
		return nil, err
	}

	var resp []SubscriptionModel
	tx := r.db.WithContext(ctx).Where("name like ?", "%"+req+"%").Find(&resp)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return resp, nil
}

func (r repository) sort(ctx context.Context, req Field) ([]SubscriptionModel, error) {
	err := checkForTable(r.db)
	if err != nil {
		return nil, err
	}

	var tx *gorm.DB
	var resp []SubscriptionModel
	switch req {
	case PRICE:
		tx = r.db.WithContext(ctx).Order("price").Find(&resp)
	case DAYS:
		tx = r.db.WithContext(ctx).Order("days").Find(&resp)
	default:
		tx = r.db.WithContext(ctx).Order("name").Find(&resp)
	}

	if tx.Error != nil {
		return nil, tx.Error
	}
	return resp, nil
}

func (r repository) filter(ctx context.Context, req FilterRequest) ([]SubscriptionModel, error) {
	err := checkForTable(r.db)
	if err != nil {
		return nil, err
	}

	var tx *gorm.DB
	var resp []SubscriptionModel
	switch req.field {
	case PRICE:
		tx = r.db.WithContext(ctx).Where("price between ? and ?", req.min, req.max).Find(&resp)
	case DAYS:
		tx = r.db.WithContext(ctx).Where("days between ? and ?", req.min, req.max).Find(&resp)
	default:
		tx = r.db.WithContext(ctx).Order("name").Find(&resp)
	}

	if tx.Error != nil {
		return nil, tx.Error
	}
	return resp, nil
}

// nolint: gocritic
func (r repository) getUsers(ctx context.Context, SchemeID string) ([]string, error) {
	err := checkForTable(r.db)
	if err != nil {
		return nil, err
	}

	var s SubscriptionModel
	s.SchemeID = SchemeID
	exi, err := CheckIfRecordExists(ctx, r.db, &s)
	if err != nil {
		return nil, err
	}
	if !exi {
		return nil, status.Errorf(codes.NotFound, "database: record with given id not found")
	}
	var subs []UserSubscription
	tx := r.db.WithContext(ctx).Model(&UserSubscription{}).Where(&s).Find(&subs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	res := make([]string, 0, len(subs))
	for _, val := range subs {
		res = append(res, val.UserID)
	}
	return res, nil
}
