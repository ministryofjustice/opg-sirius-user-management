package sirius

import (
	"encoding/json"
	"net/http"
)

type RandomReviews struct {
	LayPercentage int `json:"layPercentage"`
	PaPercentage  int `json:"paPercentage"`
	ProPercentage int `json:"proPercentage"`
	ReviewCycle   int `json:"reviewCycle"`
}

func (c *Client) RandomReviews(ctx Context) (RandomReviews, error) {
	var data RandomReviews

	req, err := c.newRequest(ctx, http.MethodGet, "/supervision-api/v1/random-review-settings", nil)
	if err != nil {
		return data, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close() //nolint:errcheck // no need to check error when closing body

	if resp.StatusCode == http.StatusUnauthorized {
		return data, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return data, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	return data, err
}
