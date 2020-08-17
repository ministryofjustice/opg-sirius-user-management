package server

import (
	"io"
	"net/http"
)

type Client interface {
	ChangePasswordClient
}

type Templates interface {
	ExecuteTemplate(io.Writer, string, interface{}) error
}

func New(webDir string, client Client, templates Templates) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", http.RedirectHandler("/change-password", http.StatusFound))
	mux.Handle("/change-password", changePassword(client, templates))

	static := http.FileServer(http.Dir(webDir + "/static"))
	mux.Handle("/public/", static)
	mux.Handle("/govuk/", static)

	return mux
}
