.PHONY: run build migrate swagger test

run:
	go run .

build:
	go build -o bin/server .

swagger:
	swag init -g main.go -o docs

test:
	go test ./...

migrate:
	@export $$(grep -v '^#' .env | xargs) && \
	psql "$$DATABASE_URL" -f migrations/001_create_users.sql && \
	psql "$$DATABASE_URL" -f migrations/002_create_media_cache.sql && \
	psql "$$DATABASE_URL" -f migrations/003_create_user_library.sql && \
	psql "$$DATABASE_URL" -f migrations/004_triggers.sql && \
	echo "✓ migrations applied"
