package sirius

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type RefDataTeamType struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

func (c *Client) TeamTypes(ctx context.Context, cookies []*http.Cookie) ([]RefDataTeamType, error) {
	var v struct {
		Data []RefDataTeamType `json:"teamType"`
	}

	types := []RefDataTeamType{}

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/reference-data?filter=teamType", nil, cookies)
	if err != nil {
		return types, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return types, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return types, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return types, errors.New("returned non-2XX response: " + strconv.Itoa(resp.StatusCode))
	}

	err = json.NewDecoder(resp.Body).Decode(&v)

	if err != nil {
		return types, err
	}

	return v.Data, nil
}
