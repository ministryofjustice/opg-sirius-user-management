package sirius

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type editLayPercentageReviewCycleRequest struct {
	ReviewCycle string   `json:"reviewCycle"`
	LayPercentage string   `json:"layPercentage"`
}

func (c *Client) EditLayPercentageReviewCycle(ctx Context, reviewCycle string, layPercentage string) (error) {
	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(editLayPercentageReviewCycleRequest{
        ReviewCycle:        reviewCycle,
		LayPercentage:      layPercentage,
	})
	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/api/v1/random-review-settings", &body)
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
		    Detail           string           `json:"detail"`
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return ValidationError{
			    Message: v.Detail,
				Errors: v.ValidationErrors,
			}
		}

		return newStatusError(resp)
	}

	return nil
}

