package server

import (
	"context"
	"log"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type MyDetailsClient interface {
	MyDetails(context.Context, []*http.Cookie) (sirius.MyDetails, error)
	HasPermission(context.Context, []*http.Cookie, string, string) (bool, error)
}

type myDetailsVars struct {
	Path      string
	SiriusURL string

	ID           int
	Firstname    string
	Surname      string
	Email        string
	PhoneNumber  string
	Organisation string
	Roles        []string
	Teams        []string

	CanEditPhoneNumber bool
}

func myDetails(logger *log.Logger, client MyDetailsClient, tmpl Template, siriusURL string) http.Handler {
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
			logger.Println("myDetails:", err)
			http.Error(w, "Could not connect to Sirius", http.StatusInternalServerError)
			return
		}

		CanEditPhoneNumber := true

		if ok, err := client.HasPermission(r.Context(), r.Cookies(), "user", "patch"); !ok {
			if err == sirius.ErrUnauthorized {
				http.Redirect(w, r, siriusURL+"/auth", http.StatusFound)
				return
			} else if err != nil {
				logger.Println("myDetails:", err)
				http.Error(w, "Could not connect to Sirius", http.StatusInternalServerError)
				return
			}

			CanEditPhoneNumber = false
		}

		vars := myDetailsVars{
			Path:               r.URL.Path,
			SiriusURL:          siriusURL,
			ID:                 myDetails.ID,
			Firstname:          myDetails.Firstname,
			Surname:            myDetails.Surname,
			Email:              myDetails.Email,
			PhoneNumber:        myDetails.PhoneNumber,
			CanEditPhoneNumber: CanEditPhoneNumber,
		}

		for _, role := range myDetails.Roles {
			if role == "OPG User" || role == "COP User" {
				vars.Organisation = role
			} else {
				vars.Roles = append(vars.Roles, role)
			}
		}

		for _, team := range myDetails.Teams {
			vars.Teams = append(vars.Teams, team.DisplayName)
		}

		if err := tmpl.ExecuteTemplate(w, "page", vars); err != nil {
			logger.Println("myDetails:", err)
		}
	})
}
