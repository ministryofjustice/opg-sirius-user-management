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
}

func deleteUser(client DeleteUserClient, tmpl Template, siriusURL string) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/delete-user/"))
		if err != nil {
			return StatusError(http.StatusNotFound)
		}

		vars := deleteUserVars{
			Path:      r.URL.Path,
			SiriusURL: siriusURL,
		}

		switch r.Method {
		case http.MethodGet:
			user, err := client.User(r.Context(), r.Cookies(), id)
			if err != nil {
				return err
			}
			vars.User = user

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			err := client.DeleteUser(r.Context(), r.Cookies(), id)
			if err != nil {
				return err
			}

			return RedirectError("/users")

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
