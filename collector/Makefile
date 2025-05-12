.PHONY: test integration-test up down

# ENV
DB_URL==postgres://testuser:testpass@localhost:5432/test_db_collector?sslmodedisable

# Unit tests
test:
	go test -v ./...

# Integration tests (Postgres must be running)
integration-test:
	set TEST_DATABASE_URL=postgres://testuser:testpass@localhost:5433/test_db_collector?sslmode=disable&& go test -tags=integration -v ./internal/clients/repository

# Start Postgres
up:
	docker-compose up -d postgres

# Stop all containers and remove volumes
down:
	docker-compose down -v
