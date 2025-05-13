.PHONY: up-all down-all create-topics

# Поднять collector и scoring-manager
up-all:
	docker-compose -f collector/docker-compose.yml up -d
	docker-compose -f scoring-manager/docker-compose.yml up -d
	docker exec kafka bash -c 'while ! nc -z localhost 9092; do echo "Waiting for Kafka..."; sleep 1; done'
	docker exec kafka kafka-topics --create --topic scoring_requests --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 || echo "Failed to create scoring_requests"
	docker exec kafka kafka-topics --create --topic scoring_results --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 || echo "Failed to create scoring_results"

# Остановить всё
down-all:
	docker-compose -f collector/docker-compose.yml down -v
	docker-compose -f scoring-manager/docker-compose.yml down -v

create-topics:
	docker exec kafka kafka-topics --create --topic scoring_requests --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 || true
	docker exec kafka kafka-topics --create --topic scoring_results --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 || true

# mockgen -source="scoring-manager/internal/client/repository/repository.go" -destination="scoring-manager/internal/pkg/mocks/repository.go" -package=mocks