.PHONY: help install dev build test clean docker-up docker-down

# Variables
BACKEND_DIR = backend
ADMIN_WEB_DIR = admin-web
PORTAL_WEB_DIR = portal-web

# Colores para output
GREEN = \033[0;32m
NC = \033[0m # No Color

help: ## Mostrar esta ayuda
	@echo "Comandos disponibles para Coordinador:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  ${GREEN}%-20s${NC} %s\n", $$1, $$2}'

## Instalación

install: install-backend install-admin install-portal ## Instalar todas las dependencias

install-backend: ## Instalar dependencias del backend
	@echo "${GREEN}Instalando dependencias del backend...${NC}"
	cd $(BACKEND_DIR) && go mod download

install-admin: ## Instalar dependencias del admin-web
	@echo "${GREEN}Instalando dependencias del admin-web...${NC}"
	cd $(ADMIN_WEB_DIR) && npm install

install-portal: ## Instalar dependencias del portal-web
	@echo "${GREEN}Instalando dependencias del portal-web...${NC}"
	cd $(PORTAL_WEB_DIR) && npm install

## Desarrollo

dev: ## Iniciar todos los servicios en modo desarrollo (requiere terminales separadas)
	@echo "${GREEN}Usa 'make dev-backend', 'make dev-admin' y 'make dev-portal' en terminales separadas${NC}"

dev-backend: ## Iniciar backend en modo desarrollo
	@echo "${GREEN}Iniciando backend...${NC}"
	cd $(BACKEND_DIR) && go run cmd/api/main.go

dev-admin: ## Iniciar admin-web en modo desarrollo
	@echo "${GREEN}Iniciando admin-web en http://localhost:3000...${NC}"
	cd $(ADMIN_WEB_DIR) && npm run dev

dev-portal: ## Iniciar portal-web en modo desarrollo
	@echo "${GREEN}Iniciando portal-web en http://localhost:3001...${NC}"
	cd $(PORTAL_WEB_DIR) && npm run dev

## Build

build: build-backend build-admin build-portal ## Build de todos los componentes

build-backend: ## Build del backend
	@echo "${GREEN}Building backend...${NC}"
	cd $(BACKEND_DIR) && go build -o bin/api cmd/api/main.go

build-admin: ## Build del admin-web
	@echo "${GREEN}Building admin-web...${NC}"
	cd $(ADMIN_WEB_DIR) && npm run build

build-portal: ## Build del portal-web
	@echo "${GREEN}Building portal-web...${NC}"
	cd $(PORTAL_WEB_DIR) && npm run build

## Testing

test: test-backend test-admin test-portal ## Ejecutar todos los tests

test-backend: ## Ejecutar tests del backend
	@echo "${GREEN}Ejecutando tests del backend...${NC}"
	cd $(BACKEND_DIR) && go test ./...

test-backend-coverage: ## Ejecutar tests del backend con cobertura
	@echo "${GREEN}Ejecutando tests del backend con cobertura...${NC}"
	cd $(BACKEND_DIR) && go test -coverprofile=coverage.out ./...
	cd $(BACKEND_DIR) && go tool cover -html=coverage.out -o coverage.html

test-admin: ## Ejecutar tests del admin-web
	@echo "${GREEN}Ejecutando tests del admin-web...${NC}"
	cd $(ADMIN_WEB_DIR) && npm run test

test-portal: ## Ejecutar tests del portal-web
	@echo "${GREEN}Ejecutando tests del portal-web...${NC}"
	cd $(PORTAL_WEB_DIR) && npm run test

## Linting

lint: lint-backend lint-admin lint-portal ## Ejecutar linters en todos los componentes

lint-backend: ## Ejecutar linter del backend
	@echo "${GREEN}Linting backend...${NC}"
	cd $(BACKEND_DIR) && golangci-lint run

lint-admin: ## Ejecutar linter del admin-web
	@echo "${GREEN}Linting admin-web...${NC}"
	cd $(ADMIN_WEB_DIR) && npm run lint

lint-portal: ## Ejecutar linter del portal-web
	@echo "${GREEN}Linting portal-web...${NC}"
	cd $(PORTAL_WEB_DIR) && npm run lint

## Docker

docker-up: ## Levantar servicios con Docker Compose
	@echo "${GREEN}Levantando servicios con Docker...${NC}"
	docker-compose up -d

docker-down: ## Detener servicios de Docker Compose
	@echo "${GREEN}Deteniendo servicios...${NC}"
	docker-compose down

docker-logs: ## Ver logs de Docker Compose
	docker-compose logs -f

docker-clean: ## Limpiar contenedores y volúmenes de Docker
	@echo "${GREEN}Limpiando Docker...${NC}"
	docker-compose down -v

## Database

db-up: ## Levantar solo PostgreSQL
	@echo "${GREEN}Levantando PostgreSQL...${NC}"
	docker-compose up -d postgres

db-down: ## Detener PostgreSQL
	docker-compose stop postgres

db-shell: ## Conectar a PostgreSQL shell
	docker-compose exec postgres psql -U coordinador -d coordinador_db

## Limpieza

clean: ## Limpiar archivos de build
	@echo "${GREEN}Limpiando archivos de build...${NC}"
	rm -rf $(BACKEND_DIR)/bin
	rm -rf $(BACKEND_DIR)/tmp
	rm -rf $(ADMIN_WEB_DIR)/build
	rm -rf $(ADMIN_WEB_DIR)/.svelte-kit
	rm -rf $(PORTAL_WEB_DIR)/build
	rm -rf $(PORTAL_WEB_DIR)/.svelte-kit

clean-deps: ## Limpiar dependencias (node_modules, vendor)
	@echo "${GREEN}Limpiando dependencias...${NC}"
	rm -rf $(ADMIN_WEB_DIR)/node_modules
	rm -rf $(PORTAL_WEB_DIR)/node_modules
	rm -rf $(BACKEND_DIR)/vendor

## Formato

fmt: fmt-backend fmt-admin fmt-portal ## Formatear código de todos los componentes

fmt-backend: ## Formatear código del backend
	@echo "${GREEN}Formateando backend...${NC}"
	cd $(BACKEND_DIR) && go fmt ./...

fmt-admin: ## Formatear código del admin-web
	@echo "${GREEN}Formateando admin-web...${NC}"
	cd $(ADMIN_WEB_DIR) && npm run format

fmt-portal: ## Formatear código del portal-web
	@echo "${GREEN}Formateando portal-web...${NC}"
	cd $(PORTAL_WEB_DIR) && npm run format
