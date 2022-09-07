package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestTeamTypes(t *testing.T) {
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
		expectedResponse []RefDataTeamType
		expectedError    error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("Some team types").
					UponReceiving("A request for team types").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/reference-data"),
						Query: dsl.MapMatcher{
							"filter": dsl.String("teamType"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.Like(map[string]interface{}{
							"teamType": dsl.EachLike(map[string]interface{}{
								"handle": dsl.String("ALLOCATIONS"),
								"label":  dsl.String("Allocations"),
							}, 1),
						}),
					})
			},
			expectedResponse: []RefDataTeamType{
				{
					Handle: "ALLOCATIONS",
					Label:  "Allocations",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				types, err := client.TeamTypes(Context{Context: context.Background()})

				assert.Equal(t, tc.expectedResponse, types)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestTeamTypesStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.TeamTypes(Context{Context: context.Background()})
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/reference-data?filter=teamType",
		Method: http.MethodGet,
	}, err)
}
