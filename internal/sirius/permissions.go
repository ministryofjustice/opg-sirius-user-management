package sirius

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type PermissionGroup struct {
	Permissions []string `json:"permissions"`
}

type PermissionSet map[string]PermissionGroup

func (ps PermissionSet) HasPermission(group string, method string) bool {
	for _, b := range ps[group].Permissions {
		if strings.EqualFold(b, method) {
			return true
		}
	}
	return false
}

type myPermissions struct {
	Data PermissionSet `json:"data"`
}

func (c *Client) HasPermission(ctx context.Context, cookies []*http.Cookie, group string, method string) (bool, error) {
	var v myPermissions

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
		return false, errors.New("returned non-2XX response: " + strconv.Itoa(resp.StatusCode))
	}

	err = json.NewDecoder(resp.Body).Decode(&v)

	if err != nil {
		return false, err
	}

	return v.Data.HasPermission(group, method), nil
}
