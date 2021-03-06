FROM node:14.13.1-alpine3.12 as asset-env

WORKDIR /app

RUN mkdir -p web/static

COPY web/assets web/assets
COPY webpack.config.js .
COPY package.json .
COPY yarn.lock .

RUN yarn && yarn build

FROM golang:1.14 as build-env

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/opg-sirius-user-management

FROM alpine:3.12

WORKDIR /go/bin

RUN apk --update --no-cache add \
    ca-certificates \
    && rm -rf /var/cache/apk/*
RUN apk --no-cache add tzdata

COPY --from=build-env /go/bin/opg-sirius-user-management opg-sirius-user-management
COPY --from=build-env /app/web/template web/template
COPY --from=asset-env /app/web/static web/static
ENTRYPOINT ["./opg-sirius-user-management"]
