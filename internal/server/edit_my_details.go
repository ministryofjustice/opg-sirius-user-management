package server

import (
	"context"
	"log"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type EditMyDetailsClient interface {
	MyDetails(context.Context, []*http.Cookie) (sirius.MyDetails, error)
	EditMyDetails(context.Context, []*http.Cookie, int, string) error
	HasPermission(context.Context, []*http.Cookie, string, string) (bool, error)
}

type editMyDetailsVars struct {
	Path      string
	SiriusURL string
	Errors    sirius.ValidationErrors

	PhoneNumber string
}

func editMyDetails(logger *log.Logger, client EditMyDetailsClient, tmpl Template, siriusURL string) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var validationErrors sirius.ValidationErrors

		switch r.Method {
		case http.MethodGet, http.MethodPost:
			break
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}

		if ok, err := client.HasPermission(r.Context(), r.Cookies(), "user", "patch"); !ok {
			if err != nil {
				return err
			}

			return RedirectError("/my-details")
		}

		myDetails, err := client.MyDetails(r.Context(), r.Cookies())
		if err != nil {
			return err
		}

		phoneNumber := myDetails.PhoneNumber

		if r.Method == http.MethodPost {
			var err error

			phoneNumber = r.FormValue("phonenumber")
			err = client.EditMyDetails(r.Context(), r.Cookies(), myDetails.ID, phoneNumber)

			if e, ok := err.(*sirius.ValidationError); ok {
				validationErrors = e.Errors
				w.WriteHeader(http.StatusBadRequest)
			} else if err != nil {
				return err
			} else {
				return RedirectError("/my-details")
			}
		}

		vars := editMyDetailsVars{
			Path:        r.URL.Path,
			SiriusURL:   siriusURL,
			Errors:      validationErrors,
			PhoneNumber: phoneNumber,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
