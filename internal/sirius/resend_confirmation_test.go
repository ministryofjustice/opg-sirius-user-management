package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestResendConfirmation(t *testing.T) {
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
		name          string
		setup         func()
		email         string
		expectedError error
	}{
		{
			name: "Created",
			setup: func() {
				pact.
					AddInteraction().
					Given("An admin user").
					UponReceiving("A request to resend a confirmation email").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/auth/resend-confirmation"),
						Body: "email=system.admin@opgtest.com",
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
					})
			},
			email: "system.admin@opgtest.com",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				err := client.ResendConfirmation(Context{Context: context.Background()}, tc.email)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestResendConfirmationStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.ResendConfirmation(Context{Context: context.Background()}, "")
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/auth/resend-confirmation",
		Method: http.MethodPost,
	}, err)
}
