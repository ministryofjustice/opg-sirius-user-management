package sirius

import (
	"encoding/json"
	"net/http"
	"sort"
)

func (c *Client) Roles(ctx Context) ([]string, error) {
	var v []string

	req, err := c.newRequest(ctx, http.MethodGet, SupervisionAPIPath + "/v1/roles", nil)
	if err != nil {
		return v, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v, err
	}
	defer resp.Body.Close() //nolint:errcheck // no need to check error when closing body

	if resp.StatusCode == http.StatusUnauthorized {
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&v)

	var roles []string
	for _, role := range v {
		if role != "COP User" && role != "OPG User" {
			roles = append(roles, role)
		}
	}

	sort.Strings(roles)

	return roles, err
}
