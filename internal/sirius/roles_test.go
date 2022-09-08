package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestRoles(t *testing.T) {
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
		expectedResponse []string
		expectedError    error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("Some roles").
					UponReceiving("A request for roles").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/roles"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body:   dsl.EachLike("System Admin", 1),
					})
			},
			expectedResponse: []string{"System Admin"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				roles, err := client.Roles(Context{Context: context.Background()})

				assert.Equal(t, tc.expectedResponse, roles)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestRolesStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.Roles(Context{Context: context.Background()})
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/roles",
		Method: http.MethodGet,
	}, err)
}
