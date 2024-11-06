APP_NAME := simple-rest
CMD_PATH := ./cmd/app/main.go
MIGRATIONS_PATH := ./migrations
DB_URL ?= "postgresql://postgres:12345678@localhost:5432/library?sslmode=disable"

# Компиляция
.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	go build -o $(APP_NAME) $(CMD_PATH)

# Запуск
.PHONY: run
run: build
	@echo "Running $(APP_NAME)..."
	./$(APP_NAME)

#Создание миграций
.PHONY: migrate-create
migrate-create:
ifndef name
	$(error Migration name not provided. Usage: make migrate-create name=<migration_name>)
endif
	@echo "Creating new migration $(name)..."
ifeq ($(OS),Windows_NT)
	powershell -Command "migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(name)"
else
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(name)
endif

# Запуск миграций
.PHONY: migrate
migrate:
	@echo "Applying migrations..."
	migrate -path $(MIGRATIONS_PATH) -database $(DB_URL) up

# Откат миграций
.PHONY: migrate-down
migrate-down:
	@echo "Reverting migrations..."
	migrate -path $(MIGRATIONS_PATH) -database $(DB_URL) down

# Очистка артефактов сборки
.PHONY: clean
clean:
	@echo "Clean..."
ifeq ($(OS),Windows_NT)
	powershell -Command "if (Test-Path $(APP_NAME)) { Remove-Item -Force $(APP_NAME) }"
else
	rm -f $(APP_NAME)
endif

# Сборка docker-compose
.PHONY: docker-up
docker-up:
	@echo "Starting Docker containers..."
	docker-compose up --build -d

# Остановка docker-compose
.PHONY: docker-down
docker-down:
	@echo "Stopping and removing Docker containers..."
	docker-compose down


# Тесты
.PHONY: test
test:
	@echo "Running tests..."
	go test ./...
