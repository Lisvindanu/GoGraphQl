.PHONY: generate run test test-race lint docker-up docker-down tidy seed-wilayah seed-harilibur migrate

generate:
	go run github.com/99designs/gqlgen generate

run:
	go run ./cmd/server

test:
	go test ./...

test-race:
	go test -race ./...

lint:
	go vet ./...

tidy:
	go mod tidy

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-build:
	docker compose build

migrate:
	psql $$DATABASE_URL -f internal/database/migrations/001_wilayah.sql
	psql $$DATABASE_URL -f internal/database/migrations/002_hari_libur.sql

seed-wilayah:
	go run ./scripts/seed_wilayah/

seed-harilibur:
	go run ./scripts/seed_harilibur/
