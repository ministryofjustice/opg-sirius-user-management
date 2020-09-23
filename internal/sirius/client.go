package sirius

import (
	"context"
	"io"
	"net/http"
)

const ErrUnauthorized ClientError = "unauthorized"

type ClientError string

func (e ClientError) Error() string {
	return string(e)
}

func NewClient(httpClient *http.Client, baseURL string) (*Client, error) {
	return &Client{
		http:    httpClient,
		baseURL: baseURL,
	}, nil
}

type Client struct {
	http    *http.Client
	baseURL string
}

func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader, cookies []*http.Cookie) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}
	var xsrfToken string
	for _, c := range cookies {
		req.AddCookie(c)
		if c.Name == "XSRF-TOKEN" {
			xsrfToken = c.Value
		}
	}
	req.Header.Add("OPG-Bypass-Membrane", "1")
	req.Header.Add("X-XSRF-TOKEN", xsrfToken)

	return req, err
}
