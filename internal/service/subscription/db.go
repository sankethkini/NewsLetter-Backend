package subscription

import (
	"context"
	"time"

	"gorm.io/gorm"
)

const (
	errDuplicateRecord  = "record already exists"
	errResourceNotFound = "resource not found"
)

type DB interface {
	addUser(context.Context, AddUserRequest) ([]UserSubscription, error)
	removeUser(context.Context, UserSchemeRequest) ([]UserSubscription, error)
	createScheme(context.Context, *SubscriptionModel) (*SubscriptionModel, error)
	getUserScheme(context.Context, UserSchemeRequest) (*UserSubscription, error)
	getSubscription(context.Context, string) (*SubscriptionModel, error)
	renew(context.Context, UserSchemeRequest, time.Time) (*UserSubscription, error)
	search(context.Context, SearchRequest) ([]SubscriptionModel, error)
	sort(context.Context, string) ([]SubscriptionModel, error)
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

func (r repository) addUser(ctx context.Context, req AddUserRequest) ([]UserSubscription, error) {

	var us UserSubscription
	us.UserID = req.UserID
	us.SchemeID = req.SchemeID
	us.Validity = req.Validity

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

	us := UserSubscription{SchemeID: req.SchemeID, UserID: req.UserID}
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

func (r repository) getUserScheme(ctx context.Context, req UserSchemeRequest) (*UserSubscription, error) {
	mod := UserSubscription{UserID: req.UserID, SchemeID: req.SchemeID}
	var resp UserSubscription
	tx := r.db.WithContext(ctx).Model(UserSubscription{}).Where(&mod).Take(&resp)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &resp, nil
}

func (r repository) getSubscription(ctx context.Context, id string) (*SubscriptionModel, error) {
	var resp SubscriptionModel
	tx := r.db.WithContext(ctx).Model(&SubscriptionModel{}).Where(SubscriptionModel{SchemeID: id}).Take(&resp)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &resp, nil
}

func (r repository) renew(ctx context.Context, req UserSchemeRequest, val time.Time) (*UserSubscription, error) {

	us1 := UserSubscription{SchemeID: req.SchemeID, UserID: req.UserID}
	tx := r.db.WithContext(ctx).Model(&us1).Update("validity", val)
	if tx.Error != nil {
		return nil, tx.Error
	}

	tx = r.db.WithContext(ctx).Find(&us1)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &us1, nil
}

func (r repository) createScheme(ctx context.Context, req *SubscriptionModel) (*SubscriptionModel, error) {

	tx := r.db.Create(&req)
	if tx.Error != nil {
		return nil, tx.Error
	}
	var ret SubscriptionModel
	tx = r.db.WithContext(ctx).Model(&SubscriptionModel{}).Where(&SubscriptionModel{SchemeID: req.SchemeID}).Take(&ret)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return &ret, nil
}

// nolint: gocritic
func (r repository) getSchemesOfUser(ctx context.Context, UserID string) ([]UserSubscription, error) {

	var subs []UserSubscription
	var schemeIDs []string
	tx := r.db.WithContext(ctx).Model(&UserSubscription{UserID: UserID}).Find(&subs)

	if tx.Error != nil {
		return nil, tx.Error
	}

	tx = r.db.WithContext(ctx).Find(&subs, schemeIDs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return subs, nil
}

func (r repository) search(ctx context.Context, req SearchRequest) ([]SubscriptionModel, error) {

	var resp []SubscriptionModel
	tx := r.db.WithContext(ctx).Model(&SubscriptionModel{}).Scopes(req.whereClause()...).Find(&resp)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return resp, nil
}

func (r repository) sort(ctx context.Context, req string) ([]SubscriptionModel, error) {

	var resp []SubscriptionModel

	tx := r.db.WithContext(ctx).Order(req).Find(&resp)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return resp, nil
}

func (r repository) filter(ctx context.Context, req FilterRequest) ([]SubscriptionModel, error) {

	var resp []SubscriptionModel
	tx := r.db.WithContext(ctx).Scopes(req.whereClause()...).Find(&resp)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return resp, nil
}

// nolint: gocritic
func (r repository) getUsers(ctx context.Context, SchemeID string) ([]string, error) {

	var subs []UserSubscription
	tx := r.db.WithContext(ctx).Model(&UserSubscription{}).Where(&UserSubscription{SchemeID: SchemeID}).Take(&subs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	res := make([]string, 0, len(subs))
	for _, val := range subs {
		res = append(res, val.UserID)
	}
	return res, nil
}

// func checkForMysqlError(err error) error {
// 	var mysqlErr *mysql.MySQLError
// 	check := errors.As(err, mysqlErr)
// 	if !check {
// 		return nil
// 	}
// 	switch mysqlErr.Number {
// 	case 1062:
// 		return status.Errorf(codes.AlreadyExists, errDuplicateRecord)
// 	case 1452:
// 		return status.Errorf(codes.NotFound, errResourceNotFound)
// 	default:
// 		return err
// 	}
// }
