.PHONY: run build lint createNewMigration migrateStatus migrateUp migrateDown mock test
DB_URL ?= "postgresql://user:admin@127.0.0.1:5432/postgres?sslmode=disable"

run:
	go run cmd/main.go

build:
	cd cmd/ && go build -o ../../maintainer-bot

lint:
	golangci-lint run

createNewMigration:
	goose -dir migrations/ create migration sql

migrateStatus:
	goose -dir migrations/ postgres $(DB_URL) status

migrateUp:
	goose -dir migrations/ postgres $(DB_URL) up

migrateDown:
	@if [ -z "${VERSION}" ]; then \
		echo "VERSION is not set"; \
		exit 1; \
	fi
	goose -dir migrations/ postgres $(DB_URL) down-to $(VERSION)

mock:
	@if [ -z "${SOURCE_DIR}" ]; then \
		echo "SOURCE_DIR is not set"; \
		exit 1; \
	fi
	go run github.com/vektra/mockery/v2@latest --dir=$(SOURCE_DIR) --output=$(SOURCE_DIR)/mocks --with-expecter=true --all

test:
	go test -v -cover ./...
