package newsletter

import (
	"encoding/json"
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	newsletterpb "github.com/sankethkini/NewsLetter-Backend/proto/newsletterpb/v1"
)

var (
	errTitle = errors.New("title is empty")
	errBody  = errors.New("body is empty")
)

type Metadata struct {
	CreateTime time.Time `gorm:"column:create_time;type:datetime;default:CURRENT_TIMESTAMP;"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP;"`
	Creator    string    `gorm:"column:creator;type:VARCHAR(100);"`
	Updater    string    `gorm:"column:updater;type:VARCHAR(100);"`
	Active     bool      `gorm:"column:active;type:tinyint;default:1;"`
}

// nolint:revive
type NewsLetterModel struct {
	Metadata     Metadata `gorm:"embedded"`
	NewsLetterID string   `gorm:"column:id;primary_key;type:VARCHAR(100)"`
	RecordID     int64    `gorm:"column:record_id;AUTOINCREMENT;type:BIGINT"`
	Title        string   `gorm:"column:title;type:TEXT"`
	Body         string   `gorm:"column:body;type:TEXT"`
}

type NewsSchemes struct {
	Metadata     Metadata `gorm:"embedded"`
	RecordID     int64    `gorm:"column:record_id;AUTOINCREMENT;type:BIGINT"`
	NewsLetterID string   `gorm:"column:news_letter_id;type:VARCHAR(100)"`
	SchemeID     string   `gorm:"column:scheme_id;type:VARCHAR(100)"`
}

type AddSchemeRequest struct {
	NewsLetterID string
	SchemeID     string
}

type EmailData struct {
	Letter newsletterpb.NewsLetter
	Scheme newsletterpb.NewsScheme
}

func (a AddSchemeRequest) validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.NewsLetterID, validation.Required, validation.Length(1, 100)),
		validation.Field(&a.SchemeID, validation.Required, validation.Length(1, 100)),
	)
}

func (mod NewsLetterModel) validate() error {
	err := validation.IsEmpty(mod.Title)
	if err {
		return errTitle
	}
	err = validation.IsEmpty(mod.Body)
	if err {
		return errBody
	}
	return nil
}

func ModelToProto(mod *NewsLetterModel) *newsletterpb.NewsLetter {
	if mod == nil {
		return &newsletterpb.NewsLetter{}
	}
	var resp newsletterpb.NewsLetter
	resp.Body = mod.Body
	resp.Title = mod.Title
	resp.NewsLetterId = mod.NewsLetterID
	return &resp
}

// nolint:govet
func SchemeToProto(mod *NewsSchemes) newsletterpb.NewsScheme {
	if mod == nil {
		return newsletterpb.NewsScheme{}
	}
	var resp newsletterpb.NewsScheme
	resp.NewsLetterId = mod.NewsLetterID
	resp.SchemeId = mod.SchemeID
	return resp
}

// nolint: govet
func tojson(n EmailData) (string, error) {
	js, err := json.Marshal(n)
	if err != nil {
		return "", err
	}
	return string(js), nil
}

// nolint:govet
func ToModel(val string) (EmailData, error) {
	var data EmailData
	err := json.Unmarshal([]byte(val), &data)
	if err != nil {
		return data, err
	}
	return data, nil
}
