package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestTeam(t *testing.T) {
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
		id               int
		name             string
		setup            func()
		cookies          []*http.Cookie
		expectedResponse Team
		expectedError    func(int) error
	}{
		{
			name: "OK",
			id:   65,
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists and supervision team exists").
					UponReceiving("A request for a team").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/team/65"),
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
								"id":          dsl.Like(65),
								"displayName": dsl.Like("Cool Team"),
								"email":       dsl.Like("coolteam@opgtest.com"),
								"phoneNumber": dsl.Like("01818118181"),
								"members": dsl.EachLike(map[string]interface{}{
									"displayName": dsl.Like("John"),
									"email":       dsl.Like("john@opgtest.com"),
								}, 1),
								"teamType": dsl.Like(map[string]interface{}{
									"label": "Very Cool",
								}),
							},
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedResponse: Team{
				ID:          65,
				DisplayName: "Cool Team",
				Email:       "coolteam@opgtest.com",
				PhoneNumber: "01818118181",
				Members: []TeamMember{
					{
						DisplayName: "John",
						Email:       "john@opgtest.com",
					},
				},
				Type: "Supervision — Very Cool",
			},
			expectedError: func(port int) error { return nil },
		},
		{
			name: "OKWithLpaTeams",
			id:   65,
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists and LPA team exists").
					UponReceiving("A request for an LPA team").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/team/65"),
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
								"id":          dsl.Like(65),
								"displayName": dsl.Like("Cool Team"),
								"members": dsl.EachLike(map[string]interface{}{
									"displayName": dsl.Like("Carline"),
									"email":       dsl.Like("carline@opgtest.com"),
								}, 1),
							},
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedResponse: Team{
				ID:          65,
				DisplayName: "Cool Team",
				Members: []TeamMember{
					{
						DisplayName: "Carline",
						Email:       "carline@opgtest.com",
					},
				},
				Type: "LPA",
			},
			expectedError: func(port int) error { return nil },
		},
		{
			name: "Unauthorized",
			id:   65,
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request for a team without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/team/65"),
						Headers: dsl.MapMatcher{
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusUnauthorized,
					})
			},
			expectedResponse: Team{},
			expectedError:    func(port int) error { return ErrUnauthorized },
		},
		{
			name: "DoesNotExist",
			id:   3589359,
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request for a team which doesn't exist").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/team/3589359"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusNotFound,
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedResponse: Team{},
			expectedError: func(port int) error {
				return StatusError{
					Code:   http.StatusNotFound,
					URL:    fmt.Sprintf("http://localhost:%d/api/team/3589359", port),
					Method: http.MethodGet,
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				team, err := client.Team(context.Background(), tc.cookies, tc.id)
				assert.Equal(t, tc.expectedResponse, team)
				assert.Equal(t, tc.expectedError(pact.Server.Port), err)
				return nil
			}))
		})
	}
}
