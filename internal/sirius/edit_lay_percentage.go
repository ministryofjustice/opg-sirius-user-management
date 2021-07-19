package sirius

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type editLayPercentageRequest struct {
	LayPercentage string   `json:"layPercentage"`
}

func (c *Client) EditLayPercentage(ctx Context, layPercentage string) (error) {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(editLayPercentageRequest{
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
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return &ValidationError{
				Errors: v.ValidationErrors,
			}
		}

		return newStatusError(resp)
	}

	return nil
}

