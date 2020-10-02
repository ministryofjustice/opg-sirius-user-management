package sirius

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type editUserRequest struct {
	ID        int      `json:"id"`
	Firstname string   `json:"firstname"`
	Surname   string   `json:"surname"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
	Locked    bool     `json:"locked"`
	Suspended bool     `json:"suspended"`
}

func (c *Client) EditUser(ctx context.Context, cookies []*http.Cookie, user AuthUser) error {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(editUserRequest{
		ID:        user.ID,
		Firstname: user.Firstname,
		Surname:   user.Surname,
		Email:     user.Email,
		Roles:     append(user.Roles, user.Organisation),
		Locked:    user.Locked,
		Suspended: user.Suspended,
	})
	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/auth/user/%d", user.ID), &body, cookies)
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
			Message string `json:"message"`
		}

		err = json.NewDecoder(resp.Body).Decode(&v)
		if err == nil {
			return ClientError(v.Message)
		}

		if err == io.EOF {
			return fmt.Errorf("returned non-200 response: %d", resp.StatusCode)
		}

		return err
	}

	return nil
}
