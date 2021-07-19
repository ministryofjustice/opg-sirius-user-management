package sirius

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
)

type editLayReviewCycleRequest struct {
	ReviewCycle int   `json:"reviewCycle"`
	LayPercentage int   `json:"layPercentage"`
}

func (c *Client) EditLayReviewCycle(ctx Context, reviewCycle string, layPercentage int) (error) {
	var body bytes.Buffer
	reviewCycleNumber, _ := strconv.Atoi(reviewCycle)

	err := json.NewEncoder(&body).Encode(editLayPercentageRequest{
        ReviewCycle: reviewCycleNumber,
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

