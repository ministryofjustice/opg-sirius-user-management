package sirius

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) DeleteUser(ctx context.Context, cookies []*http.Cookie, userID int) error {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/auth/user/%d", userID), nil, cookies)
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return newStatusError(resp)
	}

	return nil
}
