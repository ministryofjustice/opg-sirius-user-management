package sirius

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type addUserRequest struct {
	Firstname string   `json:"firstname"`
	Surname   string   `json:"surname"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
}

func (c *Client) AddUser(ctx Context, email, firstName, lastName, organisation string, roles []string) error {
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

	req, err := c.newRequest(ctx, http.MethodPost, "/api/v1/users", &body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck // no need to check error when closing body

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode != http.StatusCreated {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return ValidationError{
				Errors: v.ValidationErrors,
			}
		}

		return newStatusError(resp)
	}

	return nil
}
