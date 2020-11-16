package server

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type DeleteUserClient interface {
	User(context.Context, []*http.Cookie, int) (sirius.AuthUser, error)
	DeleteUser(context.Context, []*http.Cookie, int) error
}

type deleteUserVars struct {
	Path      string
	SiriusURL string
	User      sirius.AuthUser
	Errors    sirius.ValidationErrors
}

func deleteUser(client DeleteUserClient, tmpl Template, siriusURL string) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/delete-user/"))
		if err != nil {
			return StatusError(http.StatusNotFound)
		}

		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}

		user, err := client.User(r.Context(), r.Cookies(), id)
		if err != nil {
			return err
		}

		vars := deleteUserVars{
			Path:      r.URL.Path,
			SiriusURL: siriusURL,
			User:      user,
		}

		if r.Method == http.MethodPost {
			err := client.DeleteUser(r.Context(), r.Cookies(), id)

			if _, ok := err.(sirius.ClientError); ok {
				vars.Errors = sirius.ValidationErrors{
					"": {
						"": err.Error(),
					},
				}

				w.WriteHeader(http.StatusBadRequest)
			} else if err != nil {
				return err
			} else {
				return RedirectError("/users")
			}
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
