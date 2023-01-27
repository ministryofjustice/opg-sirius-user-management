package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type editUserRequest struct {
	ID        int      `json:"id"`
	Email     string   `json:"email,omitempty"`
	Firstname string   `json:"firstname"`
	Surname   string   `json:"surname"`
	Roles     []string `json:"roles"`
	Suspended bool     `json:"suspended"`
}

func (c *Client) EditUser(ctx Context, user AuthUser) error {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(editUserRequest{
		ID:        user.ID,
		Email:     user.Email,
		Firstname: user.Firstname,
		Surname:   user.Surname,
		Roles:     append(user.Roles, user.Organisation),
		Suspended: user.Suspended,
	})
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf("/auth/user/%d", user.ID)

	req, err := c.newRequest(ctx, http.MethodPut, requestURL, &body)
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

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return ClientError(v.Message)
		}

		return newStatusError(resp)
	}

	return nil
}
