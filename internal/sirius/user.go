package sirius

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type AuthUser struct {
	ID           int
	Firstname    string
	Surname      string
	Email        string
	Organisation string
	Roles        []string
	Locked       bool
	Suspended    bool
	Inactive     bool
}

type authUserResponse struct {
	ID        int      `json:"id"`
	Firstname string   `json:"firstname"`
	Surname   string   `json:"surname"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
	Locked    bool     `json:"locked"`
	Suspended bool     `json:"suspended"`
	Inactive  bool     `json:"inactive"`
}

func (c *Client) User(ctx context.Context, cookies []*http.Cookie, id int) (AuthUser, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/auth/user/%d", id), nil, cookies)
	if err != nil {
		return AuthUser{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return AuthUser{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return AuthUser{}, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return AuthUser{}, errors.New("returned non-2XX response: " + strconv.Itoa(resp.StatusCode))
	}

	var v authUserResponse
	err = json.NewDecoder(resp.Body).Decode(&v)

	user := AuthUser{
		ID:        v.ID,
		Firstname: v.Firstname,
		Surname:   v.Surname,
		Email:     v.Email,
		Locked:    v.Locked,
		Suspended: v.Suspended,
		Inactive:  v.Inactive,
	}

	for _, role := range v.Roles {
		if role == "OPG User" || role == "COP User" {
			user.Organisation = role
		} else {
			user.Roles = append(user.Roles, role)
		}
	}

	return user, err
}
