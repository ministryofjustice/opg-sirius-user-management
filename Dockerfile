FROM golang:1.14 as build-env

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/opg-sirius-user-management

FROM alpine:3.10

RUN apk --update --no-cache add \
    ca-certificates \
    && rm -rf /var/cache/apk/*
RUN apk --no-cache add tzdata

COPY --from=build-env /go/bin/opg-sirius-user-management /go/bin/opg-sirius-user-management
ENTRYPOINT ["/go/bin/opg-sirius-user-management"]
