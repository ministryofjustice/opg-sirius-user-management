package server

import (
	"context"
	"log"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type ChangePasswordClient interface {
	ChangePassword(ctx context.Context, cookies []*http.Cookie, currentPassword, newPassword, newPasswordConfirm string) error
}

type changePasswordVars struct {
	Path      string
	SiriusURL string
	Error     string
}

func changePassword(logger *log.Logger, client ChangePasswordClient, tmpl Template, siriusURL string) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		vars := changePasswordVars{
			Path:      r.URL.Path,
			SiriusURL: siriusURL,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			var (
				currentPassword = r.PostFormValue("currentpassword")
				password1       = r.PostFormValue("password1")
				password2       = r.PostFormValue("password2")
			)

			err := client.ChangePassword(r.Context(), r.Cookies(), currentPassword, password1, password2)

			if err == sirius.ErrUnauthorized {
				return err

			} else if err != nil {
				if _, ok := err.(sirius.ClientError); ok {
					vars.Error = err.Error()
				} else {
					logger.Println("changePassword:", err)
					vars.Error = "Could not connect to Sirius"
				}

				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)

			} else {
				return RedirectError("/my-details")
			}

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
