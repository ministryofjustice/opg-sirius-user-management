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
	Prefix    string
	Error     string
}

func changePassword(logger *log.Logger, client ChangePasswordClient, tmpl Template, prefix, siriusURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := changePasswordVars{
			Path:      r.URL.Path,
			SiriusURL: siriusURL,
			Prefix:    prefix,
		}

		switch r.Method {
		case http.MethodGet:
			if err := tmpl.ExecuteTemplate(w, "page", vars); err != nil {
				logger.Println("changePassword:", err)
			}

		case http.MethodPost:
			var (
				currentPassword = r.PostFormValue("currentpassword")
				password1       = r.PostFormValue("password1")
				password2       = r.PostFormValue("password2")
			)

			err := client.ChangePassword(r.Context(), r.Cookies(), currentPassword, password1, password2)

			if err == sirius.ErrUnauthorized {
				http.Redirect(w, r, siriusURL+"/auth", http.StatusFound)

			} else if err != nil {
				if _, ok := err.(sirius.ClientError); ok {
					vars.Error = err.Error()
				} else {
					logger.Println("changePassword:", err)
					vars.Error = "Could not connect to Sirius"
				}

				w.WriteHeader(http.StatusBadRequest)
				if err := tmpl.ExecuteTemplate(w, "page", vars); err != nil {
					logger.Println("changePassword:", err)
				}

			} else {
				http.Redirect(w, r, prefix+"/my-details", http.StatusFound)
			}

		default:
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	})
}
