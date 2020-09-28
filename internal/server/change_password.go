package server

import (
	"context"
	"log"
	"net/http"
	"net/url"

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
		switch r.Method {
		case http.MethodGet:
			vars := changePasswordVars{
				Path:      r.URL.Path,
				SiriusURL: siriusURL,
				Prefix:    prefix,
				Error:     r.FormValue("error"),
			}

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
				return
			} else if _, ok := err.(sirius.ClientError); ok {
				http.Redirect(w, r, prefix+"/change-password?error="+url.QueryEscape(err.Error()), http.StatusFound)
				return
			} else if err != nil {
				logger.Println("changePassword:", err)
				http.Redirect(w, r, prefix+"/change-password", http.StatusFound)
				return
			}

			http.Redirect(w, r, prefix+"/my-details", http.StatusFound)

		default:
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	})
}
