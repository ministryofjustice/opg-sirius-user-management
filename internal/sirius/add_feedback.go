package sirius

import (
	"bytes"
	"encoding/json"
	"github.com/ministryofjustice/opg-sirius-user-management/internal/model"
	"net/http"
)

func (c *Client) AddFeedback(ctx Context, form model.FeedbackForm) error {
	var body bytes.Buffer
	var err error
	MaximumFormStringLength := 900

	if len(form.Message) == 0 {
		return ValidationError{
			Message: "isEmpty",
		}
	}

	if len(form.Message) > MaximumFormStringLength {
		return ValidationError{
			Message: "stringLengthTooLong",
		}
	}

	err = json.NewEncoder(&body).Encode(form)
	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/api/supervision-feedback", &body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return ValidationError{
				Errors: v.ValidationErrors,
			}
		}

		return newStatusError(resp)
	}

	return nil
}
