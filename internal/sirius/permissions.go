package sirius

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type PermissionGroup struct {
	Permissions []string `json:"permissions"`
}

type PermissionSet map[string]PermissionGroup

type MyPermissions struct {
	Data PermissionSet `json:"data"`
}

func (c *Client) MyPermissions(ctx context.Context, cookies []*http.Cookie) (PermissionSet, error) {
	var v MyPermissions

	req, err := c.newRequest(ctx, http.MethodGet, "/api/permission", nil, cookies)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("returned non-2XX response: " + strconv.Itoa(resp.StatusCode))
	}

	err = json.NewDecoder(resp.Body).Decode(&v)
	return v.Data, err
}
