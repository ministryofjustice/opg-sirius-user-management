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

	testCases := map[string]struct {
		setup            func()
		cookies          []*http.Cookie
		existingPassword string
		expectedError    error
	}{
		"OK": {
			setup: func() {
				pact.
					AddInteraction().
					Given("User has password").
					UponReceiving("A request to change the password").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/auth/change-password"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
		},

		"Unauthorized": {
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists with password").
					UponReceiving("A request to change password without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/auth/change-password"),
						Headers: dsl.MapMatcher{
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusUnauthorized,
					})
			},
			expectedError: ErrUnauthorized,
		},

		"Errors": {
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists with password").
					UponReceiving("A request to change password errors").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/auth/change-password"),
						Headers: dsl.MapMatcher{
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusBadRequest,
						Body:   dsl.Match(errorsResponse{}),
					})
			},
			expectedError: ClientError("oops"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				err := client.ChangePassword(context.Background(), tc.cookies, "a", "b", "c")
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}
