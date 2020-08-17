package server

import (
	"context"
	"log"
	"net/http"
)

type ChangePasswordClient interface {
	ChangePassword(context.Context, string, string, string) error
}

func changePassword(client ChangePasswordClient, templates Templates) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			err := client.ChangePassword(r.Context(),
				r.PostFormValue("old-password"),
				r.PostFormValue("new-password"),
				r.PostFormValue("new-password-confirm"))

			if err != nil {
				log.Println("change password request:", err)
				http.Redirect(w, r, "/change-password?error=1", http.StatusFound)
				return
			}

			http.Redirect(w, r, "/my-details", http.StatusFound)

		case http.MethodGet:
			hasError := r.FormValue("error") == "1"

			if err := templates.ExecuteTemplate(w, "change-password.gotmpl", changePasswordVars{
				HasError: hasError,
			}); err != nil {
				log.Println("change-password.gotmpl:", err)
			}

		default:
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	})
}

type changePasswordVars struct {
	HasError bool
}
