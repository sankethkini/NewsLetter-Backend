package admin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
			testName: "incorrect email format",
			req:      SignInRequest{Email: "some", Password: "some password"},
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
