package server

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type EditUserClient interface {
	User(context.Context, []*http.Cookie, int) (sirius.AuthUser, error)
	EditUser(context.Context, []*http.Cookie, sirius.AuthUser) error
}

type editUserVars struct {
	Path      string
	SiriusURL string
	User      sirius.AuthUser
	Errors    sirius.ValidationErrors
}

func editUser(client EditUserClient, tmpl Template, siriusURL string) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/edit-user/"))
		if err != nil {
			return StatusError(http.StatusNotFound)
		}

		vars := editUserVars{
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
			user := sirius.AuthUser{
				ID:           id,
				Firstname:    r.PostFormValue("firstname"),
				Surname:      r.PostFormValue("surname"),
				Email:        r.PostFormValue("email"),
				Organisation: r.PostFormValue("organisation"),
				Roles:        r.PostForm["roles"],
				Suspended:    r.PostFormValue("suspended") == "Yes",
				Locked:       r.PostFormValue("locked") == "Yes",
			}
			err := client.EditUser(r.Context(), r.Cookies(), user)

			if _, ok := err.(sirius.ClientError); ok {
				vars.User = user
				vars.Errors = sirius.ValidationErrors{
					"email": {
						"": err.Error(),
					},
				}

				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			if err != nil {
				return err
			}

			return RedirectError("/users")

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
