package sirius

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

func NewClient(httpClient *http.Client, baseURL string) (*Client, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		http:    httpClient,
		baseURL: parsed,
	}, nil
}

type Client struct {
	http    *http.Client
	baseURL *url.URL
}

func (c *Client) ChangePassword(ctx context.Context, existingPassword, password, confirmPassword string) error {
	form := url.Values{
		"existingPassword": {existingPassword},
		"password":         {password},
		"confirmPassword":  {confirmPassword},
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		c.url("/auth/change-password"),
		strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp.Header.Get("Content-Type") == "application/json" {
			var v struct {
				Errors string `json:"errors"`
			}
			json.NewDecoder(resp.Body).Decode(&v)
			return errors.New(v.Errors)
		}

		return errors.New("returned non-2XX response")
	}

	return nil
}

func (c *Client) url(path string) string {
	partial, _ := url.Parse(path)

	return c.baseURL.ResolveReference(partial).String()
}
