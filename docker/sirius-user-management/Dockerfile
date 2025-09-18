FROM node:22.12.0-alpine3.19 AS asset-env

WORKDIR /app

RUN mkdir -p web/static

COPY package.json .
COPY yarn.lock .
RUN yarn install --frozen-lockfile --ignore-scripts

COPY web/assets web/assets
RUN yarn build

FROM golang:1.25-alpine AS build-env

ARG TARGETARCH
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY main.go main.go
COPY internal internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -a -installsuffix cgo -o /go/bin/opg-sirius-user-management

FROM alpine:3

WORKDIR /go/bin

RUN apk --update --no-cache add \
    ca-certificates \
    tzdata \
    && rm -rf /var/cache/apk/*

RUN apk upgrade libcrypto3 libssl3 busybox

COPY --from=build-env /go/bin/opg-sirius-user-management opg-sirius-user-management
COPY --from=asset-env /app/web/static web/static
COPY web/template web/template

RUN addgroup -S app && \
    adduser -S -g app app && \
    chown -R app:app opg-sirius-user-management web/template web/static
USER app
ENTRYPOINT ["./opg-sirius-user-management"]
