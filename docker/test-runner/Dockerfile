FROM golang:1.27-bookworm@sha256:484ef6066fa69acb059fdfeda7ba2b8f7391f2ef6abc6f9b8411e669ebd56466

RUN go install gotest.tools/gotestsum@latest

RUN go install github.com/pact-foundation/pact-go/v2@v2.4.1
RUN pact-go -l DEBUG install
