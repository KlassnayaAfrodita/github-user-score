.PHONY: test integration-test up down create-topics kafka-console-producer kafka-console-consumer

SCORING_MANAGER_ADDRESS=localhost:50052

# Unit tests
test:
	go test -v -count=1 ./...

# Integration tests (Postgres must be running)
integration-test:
	set SCORING_MANAGER_ADDRESS=$(SCORING_MANAGER_ADDRESS)&& go test -tags=integration -count=1 -v ./internal/client/grpc
