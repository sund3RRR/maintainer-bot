FROM golang:1.21.5-bookworm as builder

COPY . /build

WORKDIR /build

RUN set -eux; \
    go build -ldflags="-s -w" -o /build/maintainer-bot cmd/main.go


FROM debian:12.2

COPY --from=builder /build/maintainer-bot /app/maintainer-bot

WORKDIR /app

# Change Timezone to Europe/Moscow
RUN apt-get update && apt-get install tzdata -y
ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Install ca-certificates for TelegramBotAPI
RUN apt install ca-certificates -y

# Set up user for maintainer bot
RUN set -eux; \
    chmod +x maintainer-bot; \
    groupadd -r maintainerbot; \
    useradd -r -g maintainerbot -d /app -s /bin/bash maintainerbot; \
    chown -R maintainerbot:maintainerbot /app;

USER maintainerbot

CMD [ "./maintainer-bot" ]