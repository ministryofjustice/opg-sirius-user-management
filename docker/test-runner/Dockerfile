FROM golang:1.27-bookworm@sha256:69a7b9788769bec032d238959b61854e9ae87f57be9029ec04e9885fabf99195

RUN go install gotest.tools/gotestsum@latest

RUN go install github.com/pact-foundation/pact-go/v2@v2.4.1
RUN pact-go -l DEBUG install
