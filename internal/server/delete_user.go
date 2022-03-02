package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/handler"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/template"
)

type DeleteUserClient interface {
	User(sirius.Context, int) (sirius.AuthUser, error)
	DeleteUser(sirius.Context, int) error
}

type deleteUserVars struct {
	Path           string
	XSRFToken      string
	User           sirius.AuthUser
	Errors         sirius.ValidationErrors
	SuccessMessage string
}

func deleteUser(client DeleteUserClient, tmpl template.Template) handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		perm := myPermissions(r)
		if !perm.HasPermission("v1-users", http.MethodDelete) {
			return handler.Status(http.StatusForbidden)
		}

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/delete-user/"))

		if err != nil {
			return handler.Status(http.StatusNotFound)
		}

		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			return handler.Status(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		user, err := client.User(ctx, id)
		if err != nil {
			return err
		}

		vars := deleteUserVars{
			Path:      r.URL.Path,
			XSRFToken: ctx.XSRFToken,
			User:      user,
		}

		if r.Method == http.MethodPost {
			err := client.DeleteUser(ctx, id)

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
				vars.SuccessMessage = fmt.Sprintf("User %s %s (%s) was deleted.", user.Firstname, user.Surname, user.Email)
			}
		}

		return tmpl(w, vars)
	}
}
