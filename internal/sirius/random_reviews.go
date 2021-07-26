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
	var d RandomReviews

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/random-review-settings", nil)
	if err != nil {
		return d, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return d, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return d, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return d, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&d)
	return d, err
}
