package server

import (
	"context"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type EditMyDetailsClient interface {
	MyDetails(context.Context, []*http.Cookie) (sirius.MyDetails, error)
	EditMyDetails(context.Context, []*http.Cookie, int, string) error
	HasPermission(context.Context, []*http.Cookie, string, string) (bool, error)
}

type editMyDetailsVars struct {
	Path        string
	SiriusURL   string
	Success     bool
	Errors      sirius.ValidationErrors
	PhoneNumber string
}

func editMyDetails(client EditMyDetailsClient, tmpl Template, siriusURL string) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}

		if ok, err := client.HasPermission(r.Context(), r.Cookies(), "user", "patch"); !ok {
			if err != nil {
				return err
			}

			return StatusError(http.StatusForbidden)
		}

		myDetails, err := client.MyDetails(r.Context(), r.Cookies())
		if err != nil {
			return err
		}

		vars := editMyDetailsVars{
			Path:        r.URL.Path,
			SiriusURL:   siriusURL,
			PhoneNumber: myDetails.PhoneNumber,
		}

		if r.Method == http.MethodPost {
			vars.PhoneNumber = r.FormValue("phonenumber")
			err := client.EditMyDetails(r.Context(), r.Cookies(), myDetails.ID, vars.PhoneNumber)

			if e, ok := err.(*sirius.ValidationError); ok {
				vars.Errors = e.Errors
				w.WriteHeader(http.StatusBadRequest)
			} else if err != nil {
				return err
			} else {
				vars.Success = true
			}
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
