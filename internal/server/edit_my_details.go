package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/handler"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/template"
)

type EditMyDetailsClient interface {
	MyDetails(sirius.Context) (sirius.MyDetails, error)
	EditMyDetails(sirius.Context, int, string) error
}

type editMyDetailsVars struct {
	Path        string
	XSRFToken   string
	Success     bool
	Errors      sirius.ValidationErrors
	PhoneNumber string
}

func editMyDetails(client EditMyDetailsClient, tmpl template.Template) handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		perm := myPermissions(r)
		if !perm.HasPermission("v1-users-updatetelephonenumber", http.MethodPut) {
			return handler.Status(http.StatusForbidden)
		}

		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			return handler.Status(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		myDetails, err := client.MyDetails(ctx)
		if err != nil {
			return err
		}

		vars := editMyDetailsVars{
			Path:        r.URL.Path,
			XSRFToken:   ctx.XSRFToken,
			PhoneNumber: myDetails.PhoneNumber,
		}

		if r.Method == http.MethodPost {
			vars.PhoneNumber = r.FormValue("phonenumber")
			err := client.EditMyDetails(ctx, myDetails.ID, vars.PhoneNumber)

			if e, ok := err.(*sirius.ValidationError); ok {
				vars.Errors = e.Errors
				w.WriteHeader(http.StatusBadRequest)
			} else if err != nil {
				return err
			} else {
				vars.Success = true
			}
		}

		return tmpl(w, vars)
	}
}
