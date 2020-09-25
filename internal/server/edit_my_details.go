package server

import (
	"log"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type editMyDetailsVars struct {
	Path             string
	SiriusURL        string
	ValidationErrors map[string]map[string]string

	PhoneNumber string
}

func editMyDetails(logger *log.Logger, client MyDetailsClient, tmpl Template, siriusURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var phoneNumber string
		var validationErrors map[string]map[string]string

		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}

		if r.Method == http.MethodPost {
			var err error

			phoneNumber = r.FormValue("phonenumber")
			validationErrors, err = client.EditMyDetails(r.Context(), r.Cookies(), 32, r.FormValue("phonenumber"))

			if err == sirius.ErrUnauthorized {
				http.Redirect(w, r, siriusURL+"/auth", http.StatusFound)
				return
			} else if err != nil {
				logger.Println("editMyDetails:", err)
				http.Error(w, "Could not connect to Sirius", http.StatusInternalServerError)
				return
			}

			if validationErrors == nil {
				http.Redirect(w, r, "/my-details", http.StatusFound)
				return
			}
		} else {
			myDetails, err := client.MyDetails(r.Context(), r.Cookies())

			if err == sirius.ErrUnauthorized {
				http.Redirect(w, r, siriusURL+"/auth", http.StatusFound)
				return
			} else if err != nil {
				logger.Println("editMyDetails:", err)
				http.Error(w, "Could not connect to Sirius", http.StatusInternalServerError)
				return
			}

			phoneNumber = myDetails.PhoneNumber
		}

		vars := editMyDetailsVars{
			Path:             r.URL.Path,
			SiriusURL:        siriusURL,
			ValidationErrors: validationErrors,
			PhoneNumber:      phoneNumber,
		}

		if err := tmpl.ExecuteTemplate(w, "page", vars); err != nil {
			logger.Println("editMyDetails:", err)
		}
	})
}
