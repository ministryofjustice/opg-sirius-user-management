package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type AddTeamClient interface {
	AddTeam(ctx context.Context, cookies []*http.Cookie, name, teamType, phone, email string) (int, error)
	TeamTypes(context.Context, []*http.Cookie) ([]sirius.RefDataTeamType, error)
}

type addTeamVars struct {
	Path      string
	SiriusURL string
	TeamTypes []sirius.RefDataTeamType
	Name      string
	Service   string
	TeamType  string
	Phone     string
	Email     string
	Success   bool
	Errors    sirius.ValidationErrors
}

func addTeam(client AddTeamClient, tmpl Template, siriusURL string) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		switch r.Method {
		case http.MethodGet:
			teamTypes, err := client.TeamTypes(r.Context(), r.Cookies())
			if err != nil {
				return err
			}

			vars := addTeamVars{
				Path:      r.URL.Path,
				SiriusURL: siriusURL,
				TeamTypes: teamTypes,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			var (
				name     = r.PostFormValue("name")
				service  = r.PostFormValue("service")
				teamType = r.PostFormValue("supervision-type")
				phone    = r.PostFormValue("phone")
				email    = r.PostFormValue("email")
			)

			if service == "lpa" {
				teamType = ""
			}

			id, err := client.AddTeam(r.Context(), r.Cookies(), name, teamType, phone, email)

			if verr, ok := err.(sirius.ValidationError); ok {
				teamTypes, err := client.TeamTypes(r.Context(), r.Cookies())
				if err != nil {
					return err
				}

				vars := addTeamVars{
					Path:      r.URL.Path,
					SiriusURL: siriusURL,
					TeamTypes: teamTypes,
					Name:      name,
					Service:   service,
					TeamType:  teamType,
					Phone:     phone,
					Email:     email,
					Errors:    verr.Errors,
				}

				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)
			} else if err != nil {
				return err
			}

			return RedirectError(fmt.Sprintf("/teams/%d", id))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
