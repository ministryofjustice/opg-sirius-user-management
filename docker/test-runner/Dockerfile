FROM golang:1.26-bookworm

RUN go install gotest.tools/gotestsum@latest

RUN go install github.com/pact-foundation/pact-go/v2@v2.4.1
RUN pact-go -l DEBUG install
