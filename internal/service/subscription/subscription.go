package subscription

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/sankethkini/NewsLetter-Backend/internal/enum"
	subscriptionpb "github.com/sankethkini/NewsLetter-Backend/proto/subscriptionpb/v1"
	"gorm.io/gorm"
)

type Metadata struct {
	CreateTime time.Time `gorm:"column:create_time;type:datetime;default:CURRENT_TIMESTAMP;"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP;"`
	Creator    string    `gorm:"column:creator;type:VARCHAR(100);"`
	Updater    string    `gorm:"column:updater;type:VARCHAR(100);"`
	Active     bool      `gorm:"column:active;type:tinyint;default:1;"`
}

// nolint:revive
type SubscriptionModel struct {
	Metadata Metadata `gorm:"embedded"`
	SchemeID string   `gorm:"column:id;primary_key;type:VARCHAR(100)"`
	RecordID int64    `gorm:"column:record_id;AUTOINCREMENT;type:BIGINT"`
	Name     string   `gorm:"column:name;type:CHAR(100)"`
	Price    float64  `gorm:"column:price;type:FLOAT"`
	Days     int      `gorm:"column:days;type:INT"`
}

type UserSubscription struct {
	Metadata Metadata  `gorm:"embedded"`
	RecordID int64     `gorm:"column:record_id;AUTOINCREMENT;type:BIGINT;"`
	SchemeID string    `gorm:"primary_key;column:scheme_id;type:VARCHAR(100);"`
	UserID   string    `gorm:"primary_key;column:user_id;type:VARCHAR(100);"`
	Validity time.Time `gorm:"column:validity;type:datetime;"`
}

type AddUserRequest struct {
	UserID   string
	SchemeID string
	Validity time.Time
}

type UserSchemeRequest struct {
	UserID   string
	SchemeID string
}

type SchemeRequest struct {
	name  string
	price float64
	days  int
}

type SearchRequest struct {
	query string
}

type RenewResponse struct {
	Name     string
	Price    float64
	Days     int
	Validity string
}

type FilterRequest struct {
	field enum.Field
	min   float32
	max   float32
}

func (f FilterRequest) whereClause() []func(*gorm.DB) *gorm.DB {
	var clause []func(*gorm.DB) *gorm.DB
	switch f.field {
	case enum.PRICE:
		clause = append(clause, func(d *gorm.DB) *gorm.DB {
			return d.Where("price between ? and ?", f.min, f.max)
		})
	case enum.DAYS:
		clause = append(clause, func(d *gorm.DB) *gorm.DB {
			return d.Where("days between ? and ?", f.min, f.max)
		})
	}
	return clause
}

func (s SearchRequest) whereClause() []func(*gorm.DB) *gorm.DB {
	var clause []func(*gorm.DB) *gorm.DB
	clause = append(clause, func(d *gorm.DB) *gorm.DB {
		return d.Where("name like ?", "%"+s.query+"%")
	})
	return clause
}

func (u UserSchemeRequest) validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.UserID, validation.Required, validation.Length(1, 100)),
		validation.Field(&u.SchemeID, validation.Required, validation.Length(1, 100)),
	)
}

func (s SchemeRequest) validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.days, validation.Required, validation.Min(1)),
		validation.Field(&s.name, validation.Required, validation.Length(1, 200)),
		validation.Field(&s.price, validation.Required, validation.Min(1.0)),
	)
}

func SubModelToProto(mod *SubscriptionModel) *subscriptionpb.Scheme {
	if mod == nil {
		return &subscriptionpb.Scheme{}
	}
	return &subscriptionpb.Scheme{
		SchemeId: mod.SchemeID,
		Name:     mod.Name,
		Price:    float32(mod.Price),
		Days:     int32(mod.Days),
	}
}

func UserSubModelToProto(usr *UserSubscription) *subscriptionpb.UserSubscription {
	if usr == nil {
		return &subscriptionpb.UserSubscription{}
	}
	return &subscriptionpb.UserSubscription{
		SchemeId: usr.SchemeID,
		Validity: usr.Validity.String(),
	}
}
