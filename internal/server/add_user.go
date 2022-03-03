package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/handler"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/template"
)

type AddUserClient interface {
	AddUser(ctx sirius.Context, email, firstname, surname, organisation string, roles []string) error
	Roles(sirius.Context) ([]string, error)
}

type addUserVars struct {
	Path      string
	XSRFToken string
	Roles     []string
	Success   bool
	Errors    sirius.ValidationErrors
}

func addUser(client AddUserClient, tmpl template.Template) handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		perm := myPermissions(r)
		if !perm.HasPermission("v1-users", http.MethodPost) {
			return handler.Status(http.StatusForbidden)
		}

		ctx := getContext(r)

		roles, err := client.Roles(ctx)
		if err != nil {
			return err
		}

		vars := addUserVars{
			Path:      r.URL.Path,
			XSRFToken: ctx.XSRFToken,
			Roles:     roles,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl(w, vars)

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
				return tmpl(w, vars)
			}

			if err != nil {
				return err
			}

			vars.Success = true
			return tmpl(w, vars)

		default:
			return handler.Status(http.StatusMethodNotAllowed)
		}
	}
}
