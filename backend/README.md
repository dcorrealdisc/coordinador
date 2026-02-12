# Backend - Coordinador API

API REST en Go/Fiber para el sistema de gestión académica Coordinador.

## 🏗️ Arquitectura

### Monolito Modular

```
backend/
├── cmd/
│   └── api/              # Entry point de la aplicación
├── internal/             # Código privado del backend
│   ├── students/        # Módulo de gestión de estudiantes
│   ├── courses/         # Módulo de catálogo de cursos
│   ├── planning/        # Módulo de planificación académica
│   ├── reports/         # Módulo de reportes y analítica
│   ├── tutors/          # Módulo de tutores/monitores
│   ├── auth/            # Módulo de autenticación
│   └── shared/          # Código compartido entre módulos
├── migrations/          # Migraciones de base de datos
└── pkg/                 # Código público reutilizable
```

### Patrón CQRS Light

- **Write Path**: Endpoints transaccionales (POST, PUT, DELETE)
- **Read Path**: Endpoints de consulta (GET) optimizados con vistas materializadas

## 🚀 Inicio Rápido

### Prerrequisitos

- Go 1.21 o superior
- PostgreSQL 15 o superior

### Instalación

```bash
cd backend

# Descargar dependencias
go mod download

# Ejecutar la aplicación
go run cmd/api/main.go
```

La API estará disponible en `http://localhost:8080`

### Desarrollo

```bash
# Ejecutar con hot reload (usando air)
air

# Ejecutar tests
go test ./...

# Ejecutar tests con cobertura
go test -cover ./...

# Linter
golangci-lint run
```

## 📡 API Endpoints

### Health Check

```
GET /health
```

Respuesta:
```json
{
  "status": "ok",
  "service": "coordinador-api",
  "version": "0.1.0"
}
```

### API v1

Base URL: `/api/v1`

#### Estudiantes
- `GET /api/v1/students` - Listar estudiantes
- `GET /api/v1/students/:id` - Obtener estudiante
- `POST /api/v1/students` - Crear estudiante
- `PUT /api/v1/students/:id` - Actualizar estudiante
- `DELETE /api/v1/students/:id` - Eliminar estudiante

#### Cursos
- `GET /api/v1/courses` - Listar cursos
- `GET /api/v1/courses/:id` - Obtener curso
- `POST /api/v1/courses` - Crear curso
- `PUT /api/v1/courses/:id` - Actualizar curso

#### Planificación
- `GET /api/v1/planning/pensum` - Obtener pensum
- `GET /api/v1/planning/periods` - Listar períodos académicos
- `POST /api/v1/planning/periods` - Crear período

#### Reportes
- `GET /api/v1/reports/students/:id/transcript` - Hoja de vida académica
- `GET /api/v1/reports/courses/:id/enrollment` - Estudiantes por curso
- `GET /api/v1/reports/analytics/performance` - Análisis de desempeño

#### Tutores
- `GET /api/v1/tutors` - Listar tutores disponibles
- `POST /api/v1/tutors/assign` - Asignar tutor a curso

_Ver documentación completa de API (próximamente)_

## 🗄️ Base de Datos

### Configuración

Variables de entorno requeridas:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=coordinador
DB_PASSWORD=secure_password
DB_NAME=coordinador_db
DB_SSLMODE=disable
```

### Migraciones

```bash
# Ejecutar migraciones
go run cmd/migrate/main.go up

# Revertir última migración
go run cmd/migrate/main.go down

# Crear nueva migración
go run cmd/migrate/main.go create nombre_migracion
```

## 📦 Estructura de Módulos

Cada módulo en `internal/` sigue la misma estructura:

```
students/
├── handler.go       # HTTP handlers
├── service.go       # Lógica de negocio
├── repository.go    # Acceso a datos
├── models.go        # Estructuras de datos
└── dto.go           # Data Transfer Objects
```

### Dependencias entre capas

```
Handler → Service → Repository → DB
   ↓         ↓
  DTO     Models
```

## 🧪 Testing

### Estructura de Tests

```
students/
├── handler_test.go
├── service_test.go
└── repository_test.go
```

### Convenciones

- Tests unitarios: `*_test.go` en cada paquete
- Tests de integración: `integration_test.go`
- Mocks: `mocks/` usando mockery

### Ejecutar Tests

```bash
# Todos los tests
go test ./...

# Tests de un módulo específico
go test ./internal/students/...

# Con cobertura
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🔧 Configuración

### Variables de Entorno

Crear archivo `.env` en la raíz del backend:

```env
# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=coordinador
DB_PASSWORD=your_password
DB_NAME=coordinador_db
DB_SSLMODE=disable

# JWT
JWT_SECRET=your_jwt_secret_here
JWT_EXPIRATION=24h

# CORS
CORS_ORIGINS=http://localhost:3000,http://localhost:3001
```

## 📊 Monitoreo y Logging

### Logging

- Usar `log/slog` para logging estructurado
- Niveles: DEBUG, INFO, WARN, ERROR
- Formato JSON en producción

### Métricas

- Endpoint de métricas: `/metrics` (Prometheus format)
- Health checks: `/health`

## 🚀 Deployment

### Build

```bash
# Build para producción
go build -o bin/api cmd/api/main.go

# Ejecutar binario
./bin/api
```

### Docker

```bash
# Build imagen
docker build -t coordinador-api .

# Ejecutar container
docker run -p 8080:8080 coordinador-api
```

## 📚 Recursos

- [Fiber Documentation](https://docs.gofiber.io/)
- [Effective Go](https://go.dev/doc/effective_go)
- [ADR-001: Arquitectura General](../docs/adrs/001-arquitectura-general.md)
- [Agente Go/Backend](../docs/agents/agente-go-backend.md) _(próximamente)_

## 🤝 Contribución

Ver guía del [Agente Go/Backend](../docs/agents/agente-go-backend.md) para convenciones y best practices.

## 📄 Licencia

[Definir licencia]
