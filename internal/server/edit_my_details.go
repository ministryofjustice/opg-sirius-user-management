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

func editMyDetails(logger *log.Logger, client EditMyDetailsClient, tmpl Template, prefix, siriusURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var validationErrors sirius.ValidationErrors

		switch r.Method {
		case http.MethodGet, http.MethodPost:
			break
		default:
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}

		if ok, err := client.HasPermission(r.Context(), r.Cookies(), "user", "patch"); !ok {
			if err == sirius.ErrUnauthorized {
				http.Redirect(w, r, siriusURL+"/auth", http.StatusFound)
				return
			} else if err != nil {
				logger.Println("myDetails:", err)
				http.Error(w, "Could not connect to Sirius", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/my-details", http.StatusFound)
			return
		}

		myDetails, err := client.MyDetails(r.Context(), r.Cookies())

		if err == sirius.ErrUnauthorized {
			http.Redirect(w, r, siriusURL+"/auth", http.StatusFound)
			return
		} else if err != nil {
			logger.Println("editMyDetails:", err)
			http.Error(w, "Could not connect to Sirius", http.StatusInternalServerError)
			return
		}

		phoneNumber := myDetails.PhoneNumber

		if r.Method == http.MethodPost {
			var err error

			phoneNumber = r.FormValue("phonenumber")
			err = client.EditMyDetails(r.Context(), r.Cookies(), myDetails.ID, phoneNumber)

			if err == sirius.ErrUnauthorized {
				http.Redirect(w, r, siriusURL+"/auth", http.StatusFound)
				return
			} else if e, ok := err.(*sirius.ValidationError); ok {
				validationErrors = e.Errors
				w.WriteHeader(http.StatusBadRequest)
			} else if err != nil {
				logger.Println("editMyDetails:", err)
				http.Error(w, "Could not connect to Sirius", http.StatusInternalServerError)
				return
			} else {
				http.Redirect(w, r, prefix+"/my-details", http.StatusFound)
				return
			}
		}

		vars := editMyDetailsVars{
			Path:        r.URL.Path,
			SiriusURL:   siriusURL,
			Errors:      validationErrors,
			PhoneNumber: phoneNumber,
		}

		if err := tmpl.ExecuteTemplate(w, "page", vars); err != nil {
			logger.Println("editMyDetails:", err)
		}
	})
}
