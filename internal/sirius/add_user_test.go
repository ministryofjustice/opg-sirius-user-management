package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

type addUserBadRequestResponse struct {
	ErrorMessages *struct {
		Email *struct {
			EmailAddressLengthExceeded string `json:"emailAddressLengthExceeded" pact:"example=The input is more than 255 characters long"`
		} `json:"email"`
	} `json:"validation_errors"`
}

func TestAddUser(t *testing.T) {
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
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/api/v1/users"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"firstname": "John",
							"surname":   "Doe",
							"email":     "john.doe@example.com",
							"roles":     []string{"COP User", "other1", "other2"},
						},
					}).
					WillRespondWith(dsl.Response{
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
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/api/v1/users"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"firstname": "John",
							"surname":   "Doe",
							"email":     "john.doefkhjgerhergjgerjkrgejgerjgerjegrjhkgrehjergjgerhjkgerhjkegrhjkgerhjkegrhjkegrhjkegrhjkgerhjkgerhjkgerhjkgerhjkgerhjkgerhjkegrhjkgerhjkgerhjkgerhjkgerhjkerghjkgerhjkgerhjkgerhjkgrhjkgrehjgerhjkgerhjkegrhjkgerhjkgrerghger@example.com",
							"roles":     []string{"COP User", "other1", "other2"},
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusBadRequest,
						Body:   dsl.Match(addUserBadRequestResponse{}),
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

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

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
