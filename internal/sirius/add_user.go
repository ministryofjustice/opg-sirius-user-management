package sirius

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type addUserRequest struct {
	Firstname string   `json:"firstname"`
	Surname   string   `json:"surname"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
}

func (c *Client) AddUser(ctx context.Context, cookies []*http.Cookie, email, firstName, lastName, organisation string, roles []string) error {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(addUserRequest{
		Firstname: firstName,
		Surname:   lastName,
		Email:     email,
		Roles:     append([]string{organisation}, roles...),
	})
	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/auth/user", &body, cookies)
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

	if resp.StatusCode != http.StatusCreated {
		var v struct {
			ErrorMessages ValidationErrors `json:"errorMessages"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return ValidationError{Errors: v.ErrorMessages}
		}

		return fmt.Errorf("returned non-201 response: %d", resp.StatusCode)
	}

	return nil
}
