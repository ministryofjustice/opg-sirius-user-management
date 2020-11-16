package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type AddUserClient interface {
	AddUser(ctx sirius.Context, email, firstname, surname, organisation string, roles []string) error
}

type addUserVars struct {
	Path      string
	SiriusURL string
	XSRFToken string
	Success   bool
	Errors    sirius.ValidationErrors
}

func addUser(client AddUserClient, tmpl Template, siriusURL string) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		vars := addUserVars{
			Path:      r.URL.Path,
			SiriusURL: siriusURL,
			XSRFToken: ctx.XSRFToken,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			var (
				email        = r.PostFormValue("email")
				firstname    = r.PostFormValue("firstname")
				surname      = r.PostFormValue("surname")
				organisation = r.PostFormValue("organisation")
				roles        = r.PostForm["roles"]
			)

			err := client.AddUser(ctx, email, firstname, surname, organisation, roles)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = verr.Errors

				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			if err != nil {
				return err
			}

			vars.Success = true
			return tmpl.ExecuteTemplate(w, "page", vars)

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
