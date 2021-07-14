package sirius

import (
	"encoding/json"
	"net/http"
)

type RandomReviews struct {
	LayPercentage int   `json:"layPercentage"`
	ReviewCycle   int   `json:"reviewCycle"`
}

func (c *Client) RandomReviews(ctx Context) (RandomReviews, error) {
	var v RandomReviews

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/random-review-settings", nil)
	if err != nil {
		return v, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&v)
	return v, err
}
