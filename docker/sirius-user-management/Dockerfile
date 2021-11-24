FROM node:16.13.0-alpine3.12 as asset-env

WORKDIR /app

RUN mkdir -p web/static

COPY web/assets web/assets
COPY webpack.config.js .
COPY package.json .
COPY yarn.lock .

RUN yarn && yarn build

FROM golang:1.17.3 as build-env

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/opg-sirius-user-management

FROM alpine:3.14.3

WORKDIR /go/bin

RUN apk --update --no-cache add \
    ca-certificates \
    tzdata \
    && rm -rf /var/cache/apk/*

COPY --from=build-env /go/bin/opg-sirius-user-management opg-sirius-user-management
COPY --from=build-env /app/web/template web/template
COPY --from=asset-env /app/web/static web/static

RUN addgroup -S app && \
    adduser -S -g app app && \
    chown -R app:app opg-sirius-user-management web/template web/static
USER app
ENTRYPOINT ["./opg-sirius-user-management"]
