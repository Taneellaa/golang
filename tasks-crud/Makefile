.PHONY: run build swag clean

setup:
	@echo "📦 Установка зависимостей..."
	go mod download
	@echo "✅ Зависимости установлены"

swag:
	@echo "📖 Генерация Swagger документации..."
	swag init -g cmd/api/main.go
	@echo "✅ Документация сгенерирована"

run:
	@echo "🚀 Запуск сервера в режиме разработки..."
	go run cmd/api/main.go

build:
	@echo "🔨 Сборка приложения..."
	go build -o bin/todo-api cmd/api/main.go
	@echo "✅ Собрано: bin/todo-api"

clean:
	@echo "🧹 Очистка..."
	rm -rf bin/ docs/
	@echo "✅ Очищено"

dev: swag run

check:
	go fmt ./...
	go vet ./...