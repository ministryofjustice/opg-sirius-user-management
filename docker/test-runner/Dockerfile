FROM golang:1.26-bookworm@sha256:a688600ca24f8a4d3ca77f95b0dd40704a9fc787c826660eb7ba0b641b8b175d

RUN go install gotest.tools/gotestsum@latest

RUN go install github.com/pact-foundation/pact-go/v2@v2.4.1
RUN pact-go -l DEBUG install
