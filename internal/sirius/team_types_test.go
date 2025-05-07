package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestTeamTypes(t *testing.T) {
	pact, err := newPact()
	assert.NoError(t, err)

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
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/v1/reference-data"),
						Query: matchers.MapMatcher{
							"filter": matchers.String("teamType"),
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"teamType": matchers.EachLike(map[string]interface{}{
								"handle": matchers.String("ALLOCATIONS"),
								"label":  matchers.String("Allocations"),
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

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

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
		URL:    s.URL + "/v1/reference-data?filter=teamType",
		Method: http.MethodGet,
	}, err)
}
