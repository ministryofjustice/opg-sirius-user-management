package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/handler"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/template"
)

type AddTeamClient interface {
	AddTeam(ctx sirius.Context, name, teamType, phone, email string) (int, error)
	TeamTypes(sirius.Context) ([]sirius.RefDataTeamType, error)
}

type addTeamVars struct {
	Path      string
	XSRFToken string
	TeamTypes []sirius.RefDataTeamType
	Name      string
	Service   string
	TeamType  string
	Phone     string
	Email     string
	Success   bool
	Errors    sirius.ValidationErrors
}

func addTeam(client AddTeamClient, tmpl template.Template) handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		perm := myPermissions(r)
		if !perm.HasPermission("v1-teams", http.MethodPost) {
			return handler.Status(http.StatusForbidden)
		}

		ctx := getContext(r)

		switch r.Method {
		case http.MethodGet:
			teamTypes, err := client.TeamTypes(ctx)
			if err != nil {
				return err
			}

			vars := addTeamVars{
				Path:      r.URL.Path,
				XSRFToken: ctx.XSRFToken,
				TeamTypes: teamTypes,
			}

			return tmpl(w, vars)

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

			id, err := client.AddTeam(ctx, name, teamType, phone, email)

			if verr, ok := err.(sirius.ValidationError); ok {
				teamTypes, err := client.TeamTypes(ctx)
				if err != nil {
					return err
				}

				vars := addTeamVars{
					Path:      r.URL.Path,
					XSRFToken: ctx.XSRFToken,
					TeamTypes: teamTypes,
					Name:      name,
					Service:   service,
					TeamType:  teamType,
					Phone:     phone,
					Email:     email,
					Errors:    verr.Errors,
				}

				w.WriteHeader(http.StatusBadRequest)
				return tmpl(w, vars)
			} else if err != nil {
				return err
			}

			return handler.Redirect(fmt.Sprintf("/teams/%d", id))

		default:
			return handler.Status(http.StatusMethodNotAllowed)
		}
	}
}
