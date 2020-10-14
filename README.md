# opg-sirius-user-management

User management frontend for Sirius: Managed by opg-org-infra &amp; Terraform

## Testing

To run the Go tests use `go test ./...`, this will create a `./pacts` directory
containing the pact definition for the service which is then used for mocking in
further tests.

To run the Cypress tests locally:

- install with `npm i -g cypress@5.3.0`
- start the service `docker-compose -f docker/docker-compose.cypress.yml up -d`
- open Cypress `cypress open -P .`
