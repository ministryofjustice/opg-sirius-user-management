package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestSearchUsers(t *testing.T) {
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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/search/users"),
						Query: dsl.MapMatcher{
							"query": dsl.String("anton"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.EachLike(map[string]interface{}{
							"id":          dsl.Like(47),
							"displayName": dsl.String("Anton Mccoy"),
							"surname":     dsl.String("Mccoy"),
							"email":       dsl.String("anton.mccoy@opgtest.com"),
							"suspended":   dsl.Like(false),
							"teams": dsl.EachLike(map[string]interface{}{
								"displayName": dsl.Like("my friendly team"),
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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/search/users"),
						Query: dsl.MapMatcher{
							"query": dsl.String("admin"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.EachLike(map[string]interface{}{
							"id":          dsl.Like(47),
							"displayName": dsl.String("system admin"),
							"surname":     dsl.String("admin"),
							"email":       dsl.String("system.admin@opgtest.com"),
							"suspended":   dsl.Like(false),
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

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

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
		URL:    s.URL + "/api/v1/search/users?query=abc",
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
		URL:    s.URL + "/api/v1/search/users?query=Maria+Fern%C3%A1ndez",
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
