# opg-sirius-user-management

User management frontend for Sirius: Managed by opg-org-infra &amp; Terraform

To run locally at http://localhost:8888/ against a Sirius running on
http://localhost:8080/ use:

```
docker-compose -f docker/docker-compose.yml up -d --build
```

## Testing

The pact tests will require `pact` to be somewhere on your `$PATH`, follow the
instructions on <https://github.com/pact-foundation/pact-ruby-standalone> to
install.

You can then run the Go tests with `go test ./...`. This will create a `./pacts`
directory containing the pact definitions for the service which are then used
for mocking in further tests. 

To run the Cypress tests locally:

- install with `npm i -g cypress@5.3.0`
- start the service `docker-compose -f docker/docker-compose.cypress.yml up -d`
- open Cypress `cypress open -P .`
