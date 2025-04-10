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

func TestAddUser(t *testing.T) {
	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		setup         func()
		email         string
		firstName     string
		lastName      string
		organisation  string
		roles         []string
		expectedError error
	}{
		{
			name: "Created",
			setup: func() {
				pact.
					AddInteraction().
					Given("An admin user").
					UponReceiving("A request to add a new user").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/api/v1/users"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"firstname": "John",
							"surname":   "Doe",
							"email":     "john.doe@example.com",
							"roles":     []string{"COP User", "other1", "other2"},
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusCreated,
					})
			},
			firstName:    "John",
			lastName:     "Doe",
			email:        "john.doe@example.com",
			organisation: "COP User",
			roles:        []string{"other1", "other2"},
		},
		{
			name: "Errors",
			setup: func() {
				pact.
					AddInteraction().
					Given("An admin user").
					UponReceiving("A request to add a new user errors").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/api/v1/users"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"firstname": "John",
							"surname":   "Doe",
							"email":     "john.doefkhjgerhergjgerjkrgejgerjgerjegrjhkgrehjergjgerhjkgerhjkegrhjkgerhjkegrhjkegrhjkegrhjkgerhjkgerhjkgerhjkgerhjkgerhjkgerhjkegrhjkgerhjkgerhjkgerhjkgerhjkerghjkgerhjkgerhjkgerhjkgrhjkgrehjgerhjkgerhjkegrhjkgerhjkgrerghger@example.com",
							"roles":     []string{"COP User", "other1", "other2"},
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusBadRequest,
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/problem+json"),
						},
						Body: matchers.Like(map[string]interface{}{
							"validation_errors": matchers.Like(map[string]interface{}{
								"email": matchers.Like(map[string]interface{}{
									"emailAddressLengthExceeded": matchers.Like("The input is more than 255 characters long"),
								}),
							}),
						}),
					})
			},
			firstName:    "John",
			lastName:     "Doe",
			email:        "john.doefkhjgerhergjgerjkrgejgerjgerjegrjhkgrehjergjgerhjkgerhjkegrhjkgerhjkegrhjkegrhjkegrhjkgerhjkgerhjkgerhjkgerhjkgerhjkgerhjkegrhjkgerhjkgerhjkgerhjkgerhjkerghjkgerhjkgerhjkgerhjkgrhjkgrehjgerhjkgerhjkegrhjkgerhjkgrerghger@example.com",
			organisation: "COP User",
			roles:        []string{"other1", "other2"},
			expectedError: ValidationError{
				Errors: ValidationErrors{
					"email": {
						"emailAddressLengthExceeded": "The input is more than 255 characters long",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				err := client.AddUser(Context{Context: context.Background()}, tc.email, tc.firstName, tc.lastName, tc.organisation, tc.roles)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestAddUserStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.AddUser(Context{Context: context.Background()}, "", "", "", "", nil)
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/users",
		Method: http.MethodPost,
	}, err)
}
