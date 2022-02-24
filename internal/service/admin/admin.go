package admin

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Metadata struct {
	CreateTime time.Time `gorm:"column:create_time;type:datetime;default:CURRENT_TIMESTAMP;"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP;"`
	Creator    string    `gorm:"column:creator;type:VARCHAR(100);"`
	Updater    string    `gorm:"column:updater;type:VARCHAR(100);"`
	Active     bool      `gorm:"column:active;type:tinyint;default:1;"`
}

// nolint:revive
type AdminModel struct {
	Metadata Metadata `gorm:"embedded"`
	AdminID  string   `gorm:"column:id;primary_key;type:VARCHAR(100)"`
	RecordID int64    `gorm:"column:record_id;AUTOINCREMENT;type:BIGINT"`
	Email    string   `gorm:"column:email;type:CHAR(100);UNIQUE;NOT NULL"`
	Password string   `gorm:"column:password;type:VARCHAR(100)"`
}

type SignInRequest struct {
	Email    string
	Password string
}

// input validations.
func (s SignInRequest) validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Email, validation.Required, is.Email),
		validation.Field(&s.Password, validation.Required, validation.Length(1, 100)),
	)
}
