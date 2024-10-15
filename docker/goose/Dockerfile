FROM docker.io/golang:1.22.6-alpine3.20 AS builder

RUN GOOS=linux go install github.com/pressly/goose/v3/cmd/goose@latest 

FROM docker.io/alpine:3.20

COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

