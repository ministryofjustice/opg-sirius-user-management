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

func TestEditMyDetails(t *testing.T) {
	pact, err := newPact()
	assert.NoError(t, err)

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
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/v1/users/47/updateTelephoneNumber"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]string{
							"phoneNumber": "01210930320",
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
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
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/v1/users/47/updateTelephoneNumber"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]string{
							"phoneNumber": "85845984598649858684596849859549684568465894689498468495689645468384938743893892317571934751439574638753683761084565480713465618457365784613876481376457651471645463178546357843615971435645387364139756147361456145161587165477143576698764574569834659465974657946574569856896745229786",
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusBadRequest,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/problem+json")},
						Body: matchers.Like(map[string]interface{}{
							"detail": matchers.Like("Payload failed validation"),
							"validation_errors": matchers.Like(map[string]interface{}{
								"phoneNumber": matchers.Like(map[string]interface{}{
									"stringLengthTooLong": matchers.Like("The input is more than 255 characters long"),
								}),
							}),
						}),
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

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

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
		URL:    s.URL + "/v1/users/47/updateTelephoneNumber",
		Method: http.MethodPut,
	}, err)
}
