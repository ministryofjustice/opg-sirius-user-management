package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

type PermissionRequest struct {
	group  string
	method string
}

func TestPermissions(t *testing.T) {
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
		cookies          []*http.Cookie
		expectedResponse bool
		expectedError    error
		permission       PermissionRequest
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request to get my permissions").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/permission"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Like(map[string]interface{}{
							"data": map[string]interface{}{
								"user": map[string]interface{}{
									"permissions": dsl.EachLike("PATCH", 1),
								},
								"team": map[string]interface{}{
									"permissions": dsl.EachLike("POST", 1),
								},
							},
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			permission: PermissionRequest{
				group:  "user",
				method: "PATCH",
			},
			expectedResponse: true,
		},
		{
			name: "Unauthorized",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request to get my permissions without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/permission"),
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				hasPermission, err := client.HasPermission(getContext(tc.cookies), tc.permission.group, tc.permission.method)
				assert.Equal(t, tc.expectedResponse, hasPermission)
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

	_, err := client.HasPermission(getContext(nil), "", "")
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/permission",
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
	assert.False(t, permissions.HasPermission("team", "PATCHs"))
}
