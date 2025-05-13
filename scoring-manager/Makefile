.PHONY: test integration-test up down create-topics kafka-console-producer kafka-console-consumer

# ENV
DB_URL=postgres://testuser:testpass@localhost:5434/test_db_scoring_manager?sslmode=disable
TOPICS=scoring_requests scoring_results
BROKER=localhost:9092
COLLECTOR_GRPC_ADDR=localhost:50051

# Unit tests
test:
	go test -v -count=1 ./...

# Integration tests (Postgres must be running)
integration-test:
	set TEST_DATABASE_URL=$(DB_URL)&& go test -tags=integration -count=1 -v ./internal/client/repository
	go test -tags=integration -count=1 -v ./internal/client/kafka
	set COLLECTOR_GRPC_ADDR=$(COLLECTOR_GRPC_ADDR)&& go test -tags=integration -count=1 -v ./internal/client/grpc

# Start Postgres, Kafka, Zookeeper
up:
	docker-compose up -d postgres zookeeper kafka

# Stop all containers and remove volumes
down:
	docker-compose down -v

# Create Kafka topics
create-topics:
	docker exec kafka kafka-topics --create --topic scoring_requests --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 || true
	docker exec kafka kafka-topics --create --topic scoring_results --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 || true

# Kafka CLI: produce to topic (example: make kafka-console-producer TOPIC=scoring_requests)
kafka-console-producer:
	docker exec -it kafka kafka-console-producer --broker-list $(BROKER) --topic $(TOPIC)

# Kafka CLI: consume from topic (example: make kafka-console-consumer TOPIC=scoring_results)
kafka-console-consumer:
	docker exec -it kafka kafka-console-consumer --bootstrap-server $(BROKER) --topic $(TOPIC) --from-beginning
