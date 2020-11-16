package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type MyDetailsClient interface {
	MyDetails(sirius.Context) (sirius.MyDetails, error)
	HasPermission(sirius.Context, string, string) (bool, error)
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

func myDetails(client MyDetailsClient, tmpl Template, siriusURL string) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		myDetails, err := client.MyDetails(ctx)
		if err != nil {
			return err
		}

		canEditPhoneNumber, err := client.HasPermission(ctx, "user", "patch")
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
