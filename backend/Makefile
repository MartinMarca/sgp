.PHONY: run build test clean install migrate

# Variables
BINARY_NAME=sgp-server
MAIN_PATH=./cmd/server/main.go

# Comandos de desarrollo
run:
	@echo "Ejecutando servidor en modo desarrollo..."
	go run $(MAIN_PATH)

build:
	@echo "Compilando binario..."
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

test:
	@echo "Ejecutando tests..."
	go test -v ./...

clean:
	@echo "Limpiando archivos generados..."
	rm -rf bin/
	go clean

install:
	@echo "Descargando dependencias..."
	go mod download
	go mod tidy

migrate:
	@echo "Ejecutando migraciones..."
	go run $(MAIN_PATH) migrate

# Hot reload (requiere air: go install github.com/cosmtrek/air@latest)
dev:
	@echo "Iniciando con hot-reload..."
	air

# Generar documentación
docs:
	@echo "Generando documentación..."
	godoc -http=:6060
