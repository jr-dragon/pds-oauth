FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /usr/local/bin/pds-oauth ./cmd/pds-oauth

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

COPY --from=builder /usr/local/bin/pds-oauth /usr/local/bin/pds-oauth
COPY configs/ /etc/pds-oauth/

EXPOSE 8000

ENTRYPOINT ["pds-oauth"]
CMD ["-config", "/etc/pds-oauth/"]
