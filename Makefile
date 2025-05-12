.PHONY: up-all down-all

# Поднять collector и scoring-manager
up-all:
	docker-compose -f collector/docker-compose.yml up -d
	docker-compose -f scoring-manager/docker-compose.yml up -d

# Остановить всё
down-all:
	docker-compose -f collector/docker-compose.yml down -v
	docker-compose -f scoring-manager/docker-compose.yml down -v

# mockgen -source="scoring-manager/internal/client/repository/repository.go" -destination="scoring-manager/internal/pkg/mocks/repository.go" -package=mocks