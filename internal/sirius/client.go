package sirius

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const ErrUnauthorized ClientError = "unauthorized"

type ClientError string

func (e ClientError) Error() string {
	return string(e)
}

type ValidationErrors map[string]map[string]string

type ValidationError struct {
	Message string
	Errors  ValidationErrors
}

func (ve ValidationError) Error() string {
	return ve.Message
}

type StatusError struct {
	Code   int    `json:"code"`
	URL    string `json:"url"`
	Method string `json:"method"`
}

func newStatusError(resp *http.Response) StatusError {
	return StatusError{
		Code:   resp.StatusCode,
		URL:    resp.Request.URL.String(),
		Method: resp.Request.Method,
	}
}

func (e StatusError) Error() string {
	return fmt.Sprintf("%s %s returned %d", e.Method, e.URL, e.Code)
}

func (e StatusError) Title() string {
	return "unexpected response from Sirius"
}

func (e StatusError) Data() interface{} {
	return e
}

type Context struct {
	Context   context.Context
	Cookies   []*http.Cookie
	XSRFToken string
}

func NewClient(httpClient HTTPClient, baseURL string) (*Client, error) {
	return &Client{
		http:    httpClient,
		baseURL: baseURL,
	}, nil
}

type Client struct {
	http    HTTPClient
	baseURL string
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (c *Client) newRequest(ctx Context, method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx.Context, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	for _, c := range ctx.Cookies {
		req.AddCookie(c)
	}

	req.Header.Add("OPG-Bypass-Membrane", "1")
	req.Header.Add("X-XSRF-TOKEN", ctx.XSRFToken)

	return req, err
}

func getContext(r *http.Request) Context {
	token := ""

	if r.Method == http.MethodGet {
		if cookie, err := r.Cookie("XSRF-TOKEN"); err == nil {
			token, _ = url.QueryUnescape(cookie.Value)
		}
	} else {
		token = r.FormValue("xsrfToken")
	}

	return Context{
		Context:   r.Context(),
		Cookies:   r.Cookies(),
		XSRFToken: token,
	}
}
