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

type UnlockUserClient interface {
	User(sirius.Context, int) (sirius.AuthUser, error)
	EditUser(sirius.Context, sirius.AuthUser) error
}

type unlockUserVars struct {
	Path      string
	XSRFToken string
	User      sirius.AuthUser
	Errors    sirius.ValidationErrors
}

func unlockUser(client UnlockUserClient, tmpl template.Template) handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		perm := myPermissions(r)
		if !perm.HasPermission("v1-users", http.MethodPut) {
			return handler.Status(http.StatusForbidden)
		}

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/unlock-user/"))
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

		vars := unlockUserVars{
			Path:      r.URL.Path,
			XSRFToken: ctx.XSRFToken,
			User:      user,
		}

		if r.Method == http.MethodPost {
			user.Locked = false
			err := client.EditUser(ctx, user)

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
				return handler.Redirect(fmt.Sprintf("/edit-user/%d", user.ID))
			}
		}

		return tmpl(w, vars)
	}
}
