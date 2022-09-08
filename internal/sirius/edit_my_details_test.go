package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

type editMyDetailsBadRequestResponse struct {
	Status           int    `json:"status" pact:"example=400"`
	Detail           string `json:"detail" pact:"example=Payload failed validation"`
	ValidationErrors *struct {
		PhoneNumber *struct {
			StringLengthTooLong string `json:"stringLengthTooLong" pact:"example=The input is more than 255 characters long"`
		} `json:"phoneNumber"`
	} `json:"validation_errors"`
}

func TestEditMyDetails(t *testing.T) {
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
		phoneNumber   string
		setup         func()
		expectedError error
	}{
		{
			name:        "OK",
			phoneNumber: "01210930320",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request to change my phone number").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/api/v1/users/47/updateTelephoneNumber"),
						Body: map[string]string{
							"phoneNumber": "01210930320",
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
		},

		{
			name:        "BadRequest",
			phoneNumber: "85845984598649858684596849859549684568465894689498468495689645468384938743893892317571934751439574638753683761084565480713465618457365784613876481376457651471645463178546357843615971435645387364139756147361456145161587165477143576698764574569834659465974657946574569856896745229786",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("An invalid request to change my phone number").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/api/v1/users/47/updateTelephoneNumber"),
						Body: map[string]string{
							"phoneNumber": "85845984598649858684596849859549684568465894689498468495689645468384938743893892317571934751439574638753683761084565480713465618457365784613876481376457651471645463178546357843615971435645387364139756147361456145161587165477143576698764574569834659465974657946574569856896745229786",
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusBadRequest,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/problem+json")},
						Body:    dsl.Match(editMyDetailsBadRequestResponse{}),
					})
			},
			expectedError: &ValidationError{
				Message: "Payload failed validation",
				Errors: ValidationErrors{
					"phoneNumber": {
						"stringLengthTooLong": "The input is more than 255 characters long",
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

				err := client.EditMyDetails(Context{Context: context.Background()}, 47, tc.phoneNumber)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestEditMyDetailsStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.EditMyDetails(Context{Context: context.Background()}, 47, "")
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/users/47/updateTelephoneNumber",
		Method: http.MethodPut,
	}, err)
}
