.PHONY: up-all down-all create-topics start up down logs restart

# Поднять collector и scoring-manager
up-all:
	docker-compose -f collector/docker-compose.yml up -d
	docker-compose -f scoring-manager/docker-compose.yml up -d


# Остановить всё
down-all:
	docker-compose -f collector/docker-compose.yml down -v
	docker-compose -f scoring-manager/docker-compose.yml down -v

create-topics:
	docker exec kafka kafka-topics --create --topic scoring_requests --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 || true
	docker exec kafka kafka-topics --create --topic scoring_results --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 || true

up:
	docker-compose up --build -d
	docker exec kafka bash -c 'while ! nc -z localhost 9092; do echo "Waiting for Kafka..."; sleep 1; done'
	docker exec kafka kafka-topics --create --topic scoring_requests --bootstrap-server kafka:9092 --replication-factor 1 --partitions 1 || echo "Failed to create scoring_requests"
	docker exec kafka kafka-topics --create --topic scoring_results --bootstrap-server kafka:9092 --replication-factor 1 --partitions 1 || echo "Failed to create scoring_results"
	timeout /t 10 >nul
	echo "we are ready"

delete-topics:
	docker exec kafka kafka-topics --bootstrap-server kafka:9092 --delete --topic scoring_result || echo "Failed to delete scoring_results"
	docker exec kafka kafka-topics --bootstrap-server kafka:9092 --delete --topic scoring_request || echo "Failed to delete scoring_request"

down:
	docker-compose down -v

logs:
	docker-compose logs -f --tail=100

restart:
	make down
	make up


# mockgen -source="scoring-manager/internal/client/repository/repository.go" -destination="scoring-manager/internal/pkg/mocks/repository.go" -package=mocks