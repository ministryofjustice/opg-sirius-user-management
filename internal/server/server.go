package server

import (
	"html/template"
	"io"
	"log"
	"net/http"
)

type Client interface {
	MyDetailsClient
	EditMyDetailsClient
}

type Template interface {
	ExecuteTemplate(io.Writer, string, interface{}) error
}

func New(logger *log.Logger, client Client, templates map[string]*template.Template, siriusURL, webDir string) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", http.RedirectHandler("/my-details", http.StatusFound))
	mux.Handle("/my-details", myDetails(logger, client, templates["my-details.gotmpl"], siriusURL))
	mux.Handle("/my-details/edit", editMyDetails(logger, client, templates["edit-my-details.gotmpl"], siriusURL))

	static := http.FileServer(http.Dir(webDir + "/static"))
	mux.Handle("/assets/", static)
	mux.Handle("/javascript/", static)
	mux.Handle("/stylesheets/", static)

	return mux
}
