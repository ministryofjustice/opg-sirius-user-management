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

func myDetails(logger *log.Logger, client MyDetailsClient, tmpl Template, siriusURL string) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		myDetails, err := client.MyDetails(r.Context(), r.Cookies())
		if err != nil {
			return err
		}

		canEditPhoneNumber, err := client.HasPermission(r.Context(), r.Cookies(), "user", "patch")
		if err != nil {
			return err
		}

		vars := myDetailsVars{
			Path:               r.URL.Path,
			SiriusURL:          siriusURL,
			ID:                 myDetails.ID,
			Firstname:          myDetails.Firstname,
			Surname:            myDetails.Surname,
			Email:              myDetails.Email,
			PhoneNumber:        myDetails.PhoneNumber,
			CanEditPhoneNumber: canEditPhoneNumber,
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

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
