package sirius

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestDeleteTeam(t *testing.T) {
	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		setup         func()
		teamID        int
		expectedError error
	}{
		{
			name:   "OK",
			teamID: 461,
			setup: func() {
				pact.
					AddInteraction().
					Given("A team that can be deleted").
					UponReceiving("A request to delete the team").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodDelete,
						Path:   matchers.String("/v1/teams/461"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusNoContent,
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				err := client.DeleteTeam(Context{Context: context.Background()}, tc.teamID)

				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestDeleteTeamClientError(t *testing.T) {
	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, `{"detail":"oops"}`, http.StatusBadRequest)
		}),
	)
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.DeleteTeam(Context{Context: context.Background()}, 461)
	assert.Equal(t, ClientError("oops"), err)
}

func TestDeleteTeamStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.DeleteTeam(Context{Context: context.Background()}, 461)
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/v1/teams/461",
		Method: http.MethodDelete,
	}, err)
}
