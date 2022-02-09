package user

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/sankethkini/NewsLetter-Backend/proto/userpb/v1"
)

type Metadata struct {
	CreateTime time.Time `gorm:"column:create_time;type:datetime;default:CURRENT_TIMESTAMP;"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP;"`
	Creator    string    `gorm:"column:creator;type:VARCHAR(100);"`
	Updater    string    `gorm:"column:updater;type:VARCHAR(100);"`
	Active     bool      `gorm:"column:active;type:tinyint;default:1;"`
}

type UserModel struct {
	Metadata Metadata `gorm:"embedded"`
	UserID   string   `gorm:"column:user_id;primary_key;type:VARCHAR(100)"`
	RecordID int64    `gorm:"column:record_id;AUTOINCREMENT;type:BIGINT"`
	Name     string   `gorm:"column:name;type:CHAR(100)"`
	Email    string   `gorm:"column:email;type:CHAR(100);UNIQUE;NOT NULL"`
	Password string   `gorm:"column:password;type:VARCHAR(100)"`
}

type SignInRequest struct {
	Email    string
	Password string
}

func (s SignInRequest) validate() error {
	return validation.Validate(s.Email, validation.Required, is.Email)
}

type GetEmailRequest struct {
	ID string
}

func (s GetEmailRequest) validate() error {
	return validation.Validate(s.ID, validation.Required, validation.Length(1, 100))
}

func (m UserModel) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Email, validation.Required, is.Email),
		validation.Field(&m.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&m.Password, validation.Required, validation.Length(6, 100)),
	)
}

func ModelToProto(m *UserModel) userpb.User {
	if m == nil {
		return userpb.User{}
	}
	return userpb.User{
		Email:    m.Email,
		Name:     m.Name,
		Password: m.Password,
	}
}

func ProtoToModel(up *userpb.User) UserModel {

	return UserModel{
		Email:    up.Email,
		Password: up.Password,
		Name:     up.Name,
	}
}
