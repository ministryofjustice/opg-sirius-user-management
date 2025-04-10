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

func TestSearchUsers(t *testing.T) {
	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		searchTerm       string
		expectedResponse []User
		expectedError    error
	}{
		{
			name: "User in a team",
			setup: func() {
				pact.
					AddInteraction().
					Given("A user called Anton exists who is in a team").
					UponReceiving("A search for Anton").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/api/v1/search/users"),
						Query: matchers.MapMatcher{
							"includeSuspended": matchers.String("1"),
							"query":            matchers.String("anton"),
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.EachLike(map[string]interface{}{
							"id":          matchers.Like(47),
							"displayName": matchers.String("Anton Mccoy"),
							"surname":     matchers.String("Mccoy"),
							"email":       matchers.String("anton.mccoy@opgtest.com"),
							"suspended":   matchers.Like(false),
							"teams": matchers.EachLike(map[string]interface{}{
								"displayName": matchers.Like("my friendly team"),
							}, 1),
						}, 1),
					})
			},
			searchTerm: "anton",
			expectedResponse: []User{
				{
					ID:          47,
					DisplayName: "Anton Mccoy",
					Email:       "anton.mccoy@opgtest.com",
					Status:      "Active",
					Team:        "my friendly team",
				},
			},
		},
		{
			name: "User not in a team",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A search for admin users").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/api/v1/search/users"),
						Query: matchers.MapMatcher{
							"includeSuspended": matchers.String("1"),
							"query":            matchers.String("admin"),
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.EachLike(map[string]interface{}{
							"id":          matchers.Like(47),
							"displayName": matchers.String("system admin"),
							"surname":     matchers.String("admin"),
							"email":       matchers.String("system.admin@opgtest.com"),
							"suspended":   matchers.Like(false),
						}, 1),
					})
			},
			searchTerm: "admin",
			expectedResponse: []User{
				{
					ID:          47,
					DisplayName: "system admin",
					Email:       "system.admin@opgtest.com",
					Status:      "Active",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				users, err := client.SearchUsers(Context{Context: context.Background()}, tc.searchTerm)
				assert.Equal(t, tc.expectedResponse, users)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestSearchUsersStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.SearchUsers(Context{Context: context.Background()}, "abc")
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/search/users?includeSuspended=1&query=abc",
		Method: http.MethodGet,
	}, err)
}

func TestSearchUsersEscapesQuery(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.SearchUsers(Context{Context: context.Background()}, "Maria Fernández")
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/search/users?includeSuspended=1&query=Maria+Fern%C3%A1ndez",
		Method: http.MethodGet,
	}, err)
}

func TestSearchUsersTooShort(t *testing.T) {
	client, _ := NewClient(http.DefaultClient, "")

	users, err := client.SearchUsers(Context{Context: context.Background()}, "ad")
	assert.Nil(t, users)
	assert.Equal(t, ClientError("Search term must be at least three characters"), err)
}

func TestUserStatus(t *testing.T) {
	assert.Equal(t, "string", UserStatus("string").String())

	assert.Equal(t, "", UserStatus("string").TagColour())
	assert.Equal(t, "govuk-tag--grey", UserStatus("Suspended").TagColour())
}
