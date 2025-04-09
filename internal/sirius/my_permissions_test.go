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

func TestPermissions(t *testing.T) {
	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse PermissionSet
		expectedError    error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request to get my permissions").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/api/v1/permissions"),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.Like(map[string]interface{}{
							"v1-users": map[string]interface{}{
								"permissions": matchers.EachLike("PATCH", 1),
							},
							"v1-teams": map[string]interface{}{
								"permissions": matchers.EachLike("POST", 1),
							},
						}),
					})
			},
			expectedResponse: PermissionSet{
				"v1-users": PermissionGroup{Permissions: []string{"PATCH"}},
				"v1-teams": PermissionGroup{Permissions: []string{"POST"}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				myPermissions, err := client.MyPermissions(Context{Context: context.Background()})
				assert.Equal(t, tc.expectedResponse, myPermissions)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestPermissionsIgnoredPact(t *testing.T) {
	// We need this test to produce a specific array of permissions in the
	// response so that Cypress tests will pass. Since Pact won't let us return
	// multiple array entries from `matchers.EachLike` we have to write a separate
	// test with the specific output.
	pact, err := newIgnoredPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse PermissionSet
		expectedError    error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request to get all the permissions I need").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/api/v1/permissions"),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.Like(map[string]interface{}{
							"v1-users-updatetelephonenumber": map[string]interface{}{
								"permissions": []string{"PUT"},
							},
							"v1-users": map[string]interface{}{
								"permissions": []string{"PUT", "POST", "DELETE"},
							},
							"v1-teams": map[string]interface{}{
								"permissions": []string{"GET", "POST", "PUT", "DELETE"},
							},
							"v1-random-review-settings": map[string]interface{}{
								"permissions": []string{"GET", "POST"},
							},
						}),
					})
			},
			expectedResponse: PermissionSet{
				"v1-users-updatetelephonenumber": PermissionGroup{Permissions: []string{"PUT"}},
				"v1-users":                       PermissionGroup{Permissions: []string{"PUT", "POST", "DELETE"}},
				"v1-teams":                       PermissionGroup{Permissions: []string{"GET", "POST", "PUT", "DELETE"}},
				"v1-random-review-settings":      PermissionGroup{Permissions: []string{"GET", "POST"}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				myPermissions, err := client.MyPermissions(Context{Context: context.Background()})
				assert.Equal(t, tc.expectedResponse, myPermissions)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestHasPermissionStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.MyPermissions(Context{Context: context.Background()})
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/permissions",
		Method: http.MethodGet,
	}, err)
}

func TestPermissionSetChecksPermission(t *testing.T) {
	permissions := PermissionSet{
		"user": {
			Permissions: []string{"GET", "PATCH"},
		},
		"team": {
			Permissions: []string{"GET"},
		},
	}

	assert.True(t, permissions.HasPermission("user", "PATCH"))
	assert.True(t, permissions.HasPermission("team", "GET"))
	assert.True(t, permissions.HasPermission("team", "get"))
	assert.False(t, permissions.HasPermission("team", "PATCHs"))
}
