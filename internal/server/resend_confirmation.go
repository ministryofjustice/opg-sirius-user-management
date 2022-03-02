package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/handler"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/template"
)

type ResendConfirmationClient interface {
	ResendConfirmation(sirius.Context, string) error
}

type resendConfirmationVars struct {
	Path  string
	ID    string
	Email string
}

func resendConfirmation(client ResendConfirmationClient, tmpl template.Template) handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		perm := myPermissions(r)
		if !perm.HasPermission("v1-users", http.MethodPut) {
			return handler.Status(http.StatusForbidden)
		}

		switch r.Method {
		case http.MethodGet:
			return handler.Redirect("/users")

		case http.MethodPost:
			vars := resendConfirmationVars{
				Path:  r.URL.Path,
				ID:    r.PostFormValue("id"),
				Email: r.PostFormValue("email"),
			}

			err := client.ResendConfirmation(getContext(r), vars.Email)
			if err != nil {
				return err
			}

			return tmpl(w, vars)

		default:
			return handler.Status(http.StatusMethodNotAllowed)
		}
	}
}
