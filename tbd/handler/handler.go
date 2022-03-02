package handler

import (
	"fmt"
	"net/http"
)

type Logger interface {
	Request(*http.Request, error)
}

type Handler func(w http.ResponseWriter, r *http.Request) error

type Middleware func(Handler) Handler

type ErrorFunc func(w http.ResponseWriter, r *http.Request, statusCode int, err error)

func New(redirectBaseURL, authURL string, errorFunc ErrorFunc, middleware ...Middleware) func(Handler) http.Handler {
	return func(next Handler) http.Handler {
		f := next
		for _, m := range middleware {
			f = m(f)
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := f(w, r); err != nil {
				statusCode := http.StatusInternalServerError

				switch v := err.(type) {
				case RedirectError:
					http.Redirect(w, r, redirectBaseURL+v.RedirectTo(), http.StatusFound)
					return
				case StatusError:
					if v.StatusCode() == http.StatusUnauthorized {
						http.Redirect(w, r, authURL, http.StatusFound)
						return
					}

					if v.StatusCode() == http.StatusForbidden || v.StatusCode() == http.StatusNotFound {
						statusCode = v.StatusCode()
					}
				}

				w.WriteHeader(statusCode)
				errorFunc(w, r, statusCode, err)
			}
		})
	}
}

type RedirectError interface {
	RedirectTo() string
}

type Redirect string

func (e Redirect) Error() string {
	return "redirect to " + string(e)
}

func (e Redirect) RedirectTo() string {
	return string(e)
}

type StatusError interface {
	StatusCode() int
}

type Status int

func (e Status) Error() string {
	code := e.StatusCode()

	return fmt.Sprintf("%d %s", code, http.StatusText(code))
}

func (e Status) StatusCode() int {
	return int(e)
}
