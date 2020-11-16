package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type AllowRolesClient interface {
	MyDetails(sirius.Context) (sirius.MyDetails, error)
}

func allowRoles(client AllowRolesClient, allowedRoles ...string) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			myDetails, err := client.MyDetails(getContext(r))
			if err != nil {
				return err
			}

			for _, role := range myDetails.Roles {
				for _, allowed := range allowedRoles {
					if role == allowed {
						return next(w, r)
					}
				}
			}

			return StatusError(http.StatusForbidden)
		}
	}
}
