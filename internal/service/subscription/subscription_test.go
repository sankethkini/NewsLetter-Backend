package subscription

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAddRequestValidate(t *testing.T) {
	tests := []struct {
		testName string
		req      AddUserRequest
		expErr   bool
	}{
		{
			testName: "everything is correct",
			req:      AddUserRequest{UserID: "248wdhdnakjhq23qndia", SchemeID: "w32874du82342"},
			expErr:   false,
		},
		{
			testName: "user id field is blank",
			req:      AddUserRequest{UserID: "", SchemeID: "sjfwuieh23je23bdkj2"},
			expErr:   true,
		},
		{
			testName: "scheme id field is blank",
			req:      AddUserRequest{UserID: "8237jndwuyasdt123", SchemeID: ""},
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

func TestSubscriptionModelValidate(t *testing.T) {
	tests := []struct {
		testName string
		req      SubscriptionModel
		expErr   bool
	}{
		{
			testName: "everything is correct",
			req:      SubscriptionModel{Name: "some", Price: 200, Days: 21},
			expErr:   false,
		},
		{
			testName: "name field is blank",
			req:      SubscriptionModel{Name: "", Price: 200, Days: 21},
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

func TestUserSchemeRequestValidate(t *testing.T) {
	tests := []struct {
		testName string
		req      UserSchemeRequest
		expErr   bool
	}{
		{
			testName: "everything is correct",
			req:      UserSchemeRequest{UserID: "248wdhdnakjhq23qndia", SchemeID: "w32874du82342"},
			expErr:   false,
		},
		{
			testName: "user id field is blank",
			req:      UserSchemeRequest{UserID: "", SchemeID: "sjfwuieh23je23bdkj2"},
			expErr:   true,
		},
		{
			testName: "scheme id field is blank",
			req:      UserSchemeRequest{UserID: "8237jndwuyasdt123", SchemeID: ""},
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

func TestSubModelToProto(t *testing.T) {
	t.Logf("Name:%s", t.Name())
	name := "some"
	price := 500
	days := 30
	got := SubModelToProto(&SubscriptionModel{Price: float64(price), Days: days, Name: name})
	if got.Name != name {
		t.Errorf("expected %v got %v", name, got.Name)
	}
	if got.Price != float32(price) {
		t.Errorf("expected %v got %v", price, got.Price)
	}
	if got.Days != int32(days) {
		t.Errorf("expected %v got %v", days, got.Days)
	}
}

func TestUserSubModelToProto(t *testing.T) {
	schemeID := "2fjkhasdj121"
	val := time.Now()
	got := UserSubModelToProto(&UserSubscription{SchemeID: schemeID, Validity: val})
	if got.SchemeId != schemeID {
		t.Errorf("expected %v got %v", schemeID, got.SchemeId)
	}
	if got.Validity != val.String() {
		t.Errorf("expected %v got %v", val, got.Validity)
	}
}
