.PHONY: dep
dep:
	go mod tidy     # Оновлення та видалення невикористаних залежностей
	go mod download # Завантаження залежностей

.PHONY: lint
lint:
	golangci-lint run # Виконання аналізу коду

.PHONY: clean
clean:
	rm -rf ./build || true # Видалення папки зі збірками

.PHONY: up-dev
up-dev:
	docker-compose -f docker-compose.dev.yml up -d # Запуск докера

.PHONY: down-dev
down-dev:
	docker-compose -f docker-compose.dev.yml down # Зупинка докера

.PHONY: generate-mocks
generate-mocks:
	mockgen # Генерація моків для сервісу
	mockgen # Генерація моків для сховища

.PHONY: migrate-up
migrate-up:
	go run ./Mongo_storage/migrations/migrate.go -direction=up # Виконання міграцій вгору

.PHONY: migrate-down
migrate-down:
	go run ./Mongo_storage/migrations/migrate.go -direction=down # Виконання міграцій вниз
