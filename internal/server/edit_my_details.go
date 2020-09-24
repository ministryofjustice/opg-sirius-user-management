package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type editMyDetailsVars struct {
	Path      string
	SiriusURL string

	PhoneNumber string
}

func editMyDetails(logger *log.Logger, client MyDetailsClient, tmpl *template.Template, siriusURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "", http.StatusMethodNotAllowed)
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

		vars := editMyDetailsVars{
			Path:        r.URL.Path,
			SiriusURL:   siriusURL,
			PhoneNumber: myDetails.PhoneNumber,
		}

		if err := tmpl.ExecuteTemplate(w, "page", vars); err != nil {
			logger.Println("editMyDetails:", err)
		}
	})
}
