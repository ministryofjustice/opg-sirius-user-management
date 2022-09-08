package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

type errorsResponse struct {
	Errors string `json:"errors" pact:"example=oops"`
}

func TestChangePassword(t *testing.T) {
	pact := &dsl.Pact{
		Consumer:          "sirius-user-management",
		Provider:          "sirius",
		Host:              "localhost",
		PactFileWriteMode: "merge",
		LogDir:            "../../logs",
		PactDir:           "../../pacts",
	}
	defer pact.Teardown()

	testCases := []struct {
		name             string
		setup            func()
		existingPassword string
		password         string
		confirmPassword  string
		expectedError    error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("User has password").
					UponReceiving("A request to change the password").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/auth/change-password"),
						Body: "existingPassword=Password1&password=Password1&confirmPassword=Password1",
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
					})
			},
			existingPassword: "Password1",
			password:         "Password1",
			confirmPassword:  "Password1",
		},
		{
			name: "Errors",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists with password").
					UponReceiving("A request to change password errors").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/auth/change-password"),
						Body: "existingPassword=x&password=y&confirmPassword=z",
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusBadRequest,
						Body:   dsl.Match(errorsResponse{}),
					})
			},
			existingPassword: "x",
			password:         "y",
			confirmPassword:  "z",
			expectedError:    ClientError("oops"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				err := client.ChangePassword(Context{Context: context.Background()}, tc.existingPassword, tc.password, tc.confirmPassword)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestChangePasswordStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.ChangePassword(Context{Context: context.Background()}, "", "", "")
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/auth/change-password",
		Method: http.MethodPost,
	}, err)
}
