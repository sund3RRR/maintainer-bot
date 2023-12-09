FROM golang:1.21.5-bookworm as builder

COPY . /build

WORKDIR /build

RUN set -eux; \
    go build -o /build/maintainer-bot ./cmd

FROM debian:12.2

COPY config/config.yml /app/config/

COPY --from=builder /build/maintainer-bot /app/maintainer-bot

WORKDIR /app

RUN set -eux; \
    apt update -y && apt install ca-certificates -y; \
    chmod +x maintainer-bot; \
    groupadd -r maintainerbot; \
    useradd -r -g maintainerbot -d /app -s /bin/bash maintainerbot; \
    chown -R maintainerbot:maintainerbot /app;

USER maintainerbot

CMD [ "./maintainer-bot" ]