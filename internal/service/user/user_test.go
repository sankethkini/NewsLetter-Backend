package user

import (
	"testing"

	v1 "github.com/sankethkini/NewsLetter-Backend/proto/userpb/v1"
	"github.com/stretchr/testify/require"
)

func TestUserValidate(t *testing.T) {
	tests := []struct {
		testName string
		req      UserModel
		expErr   bool
	}{
		{
			testName: "everything is correct",
			req:      UserModel{Name: "some", Email: "someuser@some.com", Password: "some password"},
			expErr:   false,
		},
		{
			testName: "not provided the right email",
			req:      UserModel{Name: "some", Email: "some", Password: "some"},
			expErr:   true,
		},
		{
			testName: "name field is blank",
			req:      UserModel{Name: "", Email: "some", Password: "some"},
			expErr:   true,
		},
		{
			testName: "password field is short",
			req:      UserModel{Name: "some", Email: "some@some.com", Password: "some"},
			expErr:   true,
		},
	}

	for _, val := range tests {
		t.Logf("Name:%s", val.testName)

		err := val.req.validate()
		if val.expErr {
			require.NotNil(t, err)
		} else {
			require.Nil(t, err)
		}
	}
}

func TestSignInRequestValidate(t *testing.T) {
	tests := []struct {
		testName string
		req      SignInRequest
		expErr   bool
	}{
		{
			testName: "everything is correct",
			req:      SignInRequest{Email: "some@some.com", Password: "some password"},
			expErr:   false,
		},
		{
			testName: "not provided the right email",
			req:      SignInRequest{Email: "some", Password: "some password"},
			expErr:   true,
		},
	}

	for _, val := range tests {
		t.Logf("Name:%s", val.testName)

		err := val.req.validate()
		if val.expErr {
			require.NotNil(t, err)
		} else {
			require.Nil(t, err)
		}
	}
}

func TestEmailRequestValidate(t *testing.T) {
	tests := []struct {
		testName string
		req      GetEmailRequest
		expErr   bool
	}{
		{
			testName: "everything is correct",
			req:      GetEmailRequest{ID: "248wdhdnakjhq23qndia"},
			expErr:   false,
		},
		{
			testName: "id field is blank",
			req:      GetEmailRequest{ID: ""},
			expErr:   true,
		},
	}

	for _, val := range tests {
		t.Logf("Name:%s", val.testName)

		err := val.req.validate()
		if val.expErr {
			require.NotNil(t, err)
		} else {
			require.Nil(t, err)
		}
	}
}

func TestModelToProto(t *testing.T) {
	t.Logf("Name:%s", t.Name())
	name := "some"
	email := "someuser@some.com"
	uID := "someid12344"
	got := ModelToProto(&UserModel{Name: name, Email: email, UserID: uID})
	if got.Name != name {
		t.Errorf("expected %v got %v", name, got.Name)
	}
	if got.Email != email {
		t.Errorf("expected %v got %v", email, got.Email)
	}
	if got.UserId != uID {
		t.Errorf("expected %v got %v", uID, got.UserId)
	}
}

func TestProtoToModel(t *testing.T) {
	t.Logf("Name:%s", t.Name())
	name := "some"
	email := "someuser@some.com"
	password := "some password"
	got := ProtoToModel(&v1.User{Email: email, Password: password, Name: name})
	if got.Name != name {
		t.Errorf("expected %v got %v", name, got.Name)
	}
	if got.Email != email {
		t.Errorf("expected %v got %v", email, got.Email)
	}
	if got.Password != password {
		t.Errorf("expected %v got %v", password, got.Password)
	}
}
