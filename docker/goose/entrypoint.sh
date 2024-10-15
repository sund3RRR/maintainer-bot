#!/bin/sh
set -e

for var in $(env | grep '^POSTGRESQL_URL_' | awk -F= '{print $1}'); do
    URL=$(eval echo \$$var)
    if [ -n "$URL" ]; then
        echo "Running migrations on PostgreSQL $URL"
        goose -dir /migrations/postgres postgres "$URL" up
    fi
done

if [ "$DO_NOT_EXIT" = "1" ]; then
    tail -f /dev/null
fi

exit 0
