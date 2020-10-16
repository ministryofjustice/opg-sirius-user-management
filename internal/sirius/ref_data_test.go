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
		cookies          []*http.Cookie
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
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.Like(map[string]interface{}{
							"teamType": []map[string]interface{}{
								{
									"handle": dsl.String("ALLOCATIONS"),
									"label":  dsl.String("Allocations"),
								},
								{
									"handle": dsl.String("COMPLAINTS"),
									"label":  dsl.String("Complaints"),
								},
								{
									"handle": dsl.String("INVESTIGATIONS"),
									"label":  dsl.String("Investigations"),
								},
							},
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedResponse: []RefDataTeamType{
				{
					Handle: "ALLOCATIONS",
					Label:  "Allocations",
				},
				{
					Handle: "COMPLAINTS",
					Label:  "Complaints",
				},
				{
					Handle: "INVESTIGATIONS",
					Label:  "Investigations",
				},
			},
		},

		{
			name: "Unauthorized",
			setup: func() {
				pact.
					AddInteraction().
					Given("Some team types").
					UponReceiving("A request for team types without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/reference-data"),
						Query: dsl.MapMatcher{
							"filter": dsl.String("teamType"),
						},
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

				types, err := client.TeamTypes(context.Background(), tc.cookies)

				assert.Equal(t, tc.expectedResponse, types)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}
