FROM node:24.22.3-alpine3.24@sha256:116b800c3754b440ee0aa3ed2f20f5bf3d99a554bfd9adeea09ea8017f6e85f5 AS asset-env

WORKDIR /app

RUN mkdir -p web/static

COPY package.json .
COPY package-lock.json .
RUN npm ci --ignore-scripts

COPY web/assets web/assets
RUN npm run build

FROM golang:1.27.4@sha256:205c7636d74e0ffbe3d007162df7d74f859f5a8c8322dfb9f907977686c126cf AS build-env

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid 65532 \
    app

ARG TARGETARCH
WORKDIR /app

RUN update-ca-certificates

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY main.go main.go
COPY internal internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -a -installsuffix cgo -o /go/bin/opg-sirius-user-management


FROM build-env AS healthcheck-build
WORKDIR /app

COPY healthcheck healthcheck

WORKDIR /app/healthcheck
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/healthcheck

FROM scratch

WORKDIR /go/bin

COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-env /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build-env /etc/passwd /etc/passwd
COPY --from=build-env /etc/group /etc/group

COPY --from=build-env /go/bin/opg-sirius-user-management opg-sirius-user-management
COPY --from=healthcheck-build /go/bin/healthcheck healthcheck
COPY --from=asset-env /app/web/static web/static
COPY web/template web/template

USER app

HEALTHCHECK --interval=5s --timeout=5s --start-period=5s --retries=3 CMD [ "/go/bin/healthcheck" ]

ENTRYPOINT ["./opg-sirius-user-management"]
