package newsletter

import (
	"reflect"
	"testing"

	v1 "github.com/sankethkini/NewsLetter-Backend/proto/newsletterpb/v1"
	"github.com/stretchr/testify/require"
)

func TestAddSchemeRequestValidate(t *testing.T) {
	tests := []struct {
		testName string
		req      AddSchemeRequest
		expErr   bool
	}{
		{
			testName: "everything is correct",
			req:      AddSchemeRequest{SchemeID: "2131c12w", NewsLetterID: "asdad1212"},
			expErr:   false,
		},
		{
			testName: "scheme id is blank",
			req:      AddSchemeRequest{SchemeID: "", NewsLetterID: "123123sadsad"},
			expErr:   true,
		},
		{
			testName: "user id is blank",
			req:      AddSchemeRequest{SchemeID: "asdab123b12", NewsLetterID: ""},
			expErr:   true,
		},
	}

	for _, val := range tests {
		t.Logf("Name:%s", t.Name())
		t.Logf("Description:%s", val.testName)

		err := val.req.validate()
		if val.expErr {
			require.NotNil(t, err)
		} else {
			require.Nil(t, err)
		}
	}
}

func TestNewsLetterModelValidate(t *testing.T) {
	tests := []struct {
		testName string
		req      NewsLetterModel
		expErr   bool
	}{
		{
			testName: "everything is correct",
			req:      NewsLetterModel{Title: "some title", Body: "sadsada asdas"},
			expErr:   false,
		},
		{
			testName: "title is blank",
			req:      NewsLetterModel{Title: "", Body: "soesnskdf sdkfnsd"},
			expErr:   true,
		},
		{
			testName: "user id is blank",
			req:      NewsLetterModel{Title: "sda sadsadas", Body: ""},
			expErr:   true,
		},
	}

	for _, val := range tests {
		t.Logf("Name:%s", t.Name())
		t.Logf("Description:%s", val.testName)

		err := val.req.validate()
		if val.expErr {
			require.NotNil(t, err)
		} else {
			require.Nil(t, err)
		}
	}
}

func TestModelToProto(t *testing.T) {
	body := "sjdnfs sdfnskdfn"
	title := "some sdfjn sdfn"
	nID := "someid1212"
	got := ModelToProto(&NewsLetterModel{NewsLetterID: nID, Body: body, Title: title})
	if got.Body != body {
		t.Errorf("expected %v got %v", body, got.Body)
	}
	if got.Title != title {
		t.Errorf("expected %v got %v", title, got.Title)
	}
	if got.NewsLetterId != nID {
		t.Errorf("expected %v got %v", nID, got.NewsLetterId)
	}
}

func TestSchemeToProto(t *testing.T) {
	nID := "someid1212"
	sID := "someid6767"
	got := SchemeToProto(&NewsSchemes{NewsLetterID: nID, SchemeID: sID})
	if got.NewsLetterId != nID {
		t.Errorf("expected %v got %v", nID, got.NewsLetterId)
	}
	if got.SchemeId != sID {
		t.Errorf("expected %v got %v", sID, got.SchemeId)
	}
}

// nolint:govet
func TestToAndFromJSON(t *testing.T) {
	data := EmailData{Letter: v1.NewsLetter{Title: "sssd asdasd"}, Scheme: v1.NewsScheme{NewsLetterId: "som12", SchemeId: "another123"}}
	str, err := toJSON(data)
	require.Nil(t, err)
	got, err := ToModel(str)
	require.Nil(t, err)
	if !reflect.DeepEqual(data, got) {
		t.Errorf("not equal expected %v got %v", data, got)
	}
}
