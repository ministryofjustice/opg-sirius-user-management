package sirius

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type EditRandomReview struct {
	LayPercentage string `json:"layPercentage"`
	PaPercentage  string `json:"paPercentage"`
	ProPercentage string `json:"proPercentage"`
	ReviewCycle   string `json:"reviewCycle"`
}

func (c *Client) EditRandomReviewSettings(ctx Context, reviewSettings EditRandomReview) error {
	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(reviewSettings)
	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPost, SupervisionAPIPath + "/v1/random-review-settings", &body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck // no need to check error when closing body

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		var v struct {
			Detail           string           `json:"detail"`
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return ValidationError{
				Message: v.Detail,
				Errors:  v.ValidationErrors,
			}
		}

		return newStatusError(resp)
	}

	return nil
}
