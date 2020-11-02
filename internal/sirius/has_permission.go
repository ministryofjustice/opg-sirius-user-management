package sirius

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type PermissionGroup struct {
	Permissions []string `json:"permissions"`
}

type permissionSet map[string]PermissionGroup

func (ps permissionSet) HasPermission(group string, method string) bool {
	for _, b := range ps[group].Permissions {
		if strings.EqualFold(b, method) {
			return true
		}
	}
	return false
}

type myPermissions struct {
	Data permissionSet `json:"data"`
}

func (c *Client) HasPermission(ctx context.Context, cookies []*http.Cookie, group string, method string) (bool, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/permission", nil, cookies)
	if err != nil {
		return false, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return false, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return false, newStatusError(resp)
	}

	var v myPermissions
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return false, err
	}

	return v.Data.HasPermission(group, method), nil
}
