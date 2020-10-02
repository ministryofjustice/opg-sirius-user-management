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
	DisplayName  string
	Firstname    string
	Surname      string
	Email        string
	Organisation string
	Roles        []string
	Locked       bool
	Suspended    bool
}

type authUserResponse struct {
	ID          int      `json:"id"`
	DisplayName string   `json:"displayName"`
	Firstname   string   `json:"firstname"`
	Surname     string   `json:"surname"`
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
	Locked      bool     `json:"locked"`
	Suspended   bool     `json:"suspended"`
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
		ID:          v.ID,
		DisplayName: v.DisplayName,
		Firstname:   v.Firstname,
		Surname:     v.Surname,
		Email:       v.Email,
		Locked:      v.Locked,
		Suspended:   v.Suspended,
	}

	if len(v.Roles) > 0 {
		user.Organisation = v.Roles[0]
		user.Roles = v.Roles[1:]
	}

	return user, err
}
