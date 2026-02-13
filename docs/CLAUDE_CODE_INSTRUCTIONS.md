# 🚀 Instrucciones para Claude Code

## 📋 Contexto Actual del Proyecto

Estás trabajando en **Coordinador**, un sistema de gestión académica para maestrías.

**Estado actual**:
- ✅ Arquitectura definida (ADR-001: Monolito modular + CQRS Light)
- ✅ Base de datos diseñada (ADR-002: PostgreSQL con 18 tablas + 7 vistas)
- ✅ Migraciones aplicadas (schema completo en PostgreSQL)
- ✅ 3 agentes activos: Arquitecto, DBA, Go/Backend

**Próximo paso**: Implementar el backend en Go/Fiber

---

## 🎯 Tu Rol: Agente Go/Backend

**Lee primero**: `/docs/agents/agente-go-backend.md`

### Responsabilidades:
1. Implementar APIs REST con Go/Fiber
2. Crear models, repositories, services, handlers
3. Escribir tests (>80% coverage)
4. Mantener código limpio y documentado

---

## 📊 Arquitectura del Backend

```
Layered Architecture:
HTTP Request → Handler → Service → Repository → PostgreSQL
```

**Estructura de directorios**:
```
backend/
├── cmd/api/main.go              # Entry point (YA EXISTE)
├── internal/
│   ├── models/                  # Structs del dominio (CREAR)
│   ├── repositories/            # Acceso a datos (CREAR)
│   ├── services/                # Lógica de negocio (CREAR)
│   ├── handlers/                # HTTP handlers (CREAR)
│   ├── database/                # Conexión DB (CREAR)
│   └── shared/                  # Utilidades (CREAR)
├── migrations/                  # SQL (YA EXISTE)
├── go.mod                       # Dependencias (YA EXISTE)
└── go.sum                       # Checksums (YA EXISTE)
```

---

## 🎯 Primera Tarea: Implementar Módulo de Estudiantes

### Orden de implementación:

1. **Setup de Database Connection**
   ```bash
   # Crear: internal/database/postgres.go
   # Implementar conexión con pgx/v5
   ```

2. **Modelo de Estudiante**
   ```bash
   # Crear: internal/models/student.go
   # Ver ejemplo completo en agente-go-backend.md
   ```

3. **Repository de Estudiante**
   ```bash
   # Crear: internal/repositories/student_repository.go
   # CRUD completo + filtros
   ```

4. **Service de Estudiante**
   ```bash
   # Crear: internal/services/student_service.go
   # Validaciones y lógica de negocio
   ```

5. **Handler de Estudiante**
   ```bash
   # Crear: internal/handlers/student_handler.go
   # Endpoints REST
   ```

6. **Shared Utilities**
   ```bash
   # Crear: internal/shared/response.go
   # Helpers para responses HTTP
   ```

7. **Actualizar main.go**
   ```bash
   # Conectar DB, registrar rutas
   ```

---

## 🧪 Testing

Después de cada capa:

```bash
# Unit tests (services)
go test ./internal/services/... -v

# Integration tests (repositories)
go test ./internal/repositories/... -v

# Coverage
go test ./... -cover
```

**Objetivo**: >80% coverage

---

## 📚 Referencias Importantes

1. **Decisiones arquitectónicas**:
   - `/docs/adrs/001-arquitectura-general.md`
   - `/docs/adrs/002-diseno-base-datos.md`

2. **Schema de base de datos**:
   - `/backend/migrations/003_create_people_tables.sql` (tabla students)

3. **Guía completa del agente**:
   - `/docs/agents/agente-go-backend.md`

---

## ⚙️ Comandos Útiles

```bash
# Instalar dependencias
cd backend && go mod download

# Ejecutar servidor
go run cmd/api/main.go

# Tests
go test ./... -v

# Crear archivo
# (Claude Code puede hacerlo directamente)

# Ver tablas en PostgreSQL
docker exec -i coordinador_db psql -U coordinador -d coordinador_db -c "\dt"

# Aplicar migraciones (si es necesario)
docker exec -i coordinador_db psql -U coordinador -d coordinador_db < migrations/001_create_base_schema.sql
```

---

## 🎯 Endpoints a Implementar (Estudiantes)

```
POST   /api/v1/students          # Crear estudiante
GET    /api/v1/students          # Listar estudiantes (con filtros y paginación)
GET    /api/v1/students/:id      # Obtener estudiante por ID
PUT    /api/v1/students/:id      # Actualizar estudiante
DELETE /api/v1/students/:id      # Eliminar estudiante (soft delete)
```

**Query params para GET /api/v1/students**:
- `status` (active, graduated, withdrawn, suspended)
- `cohort` (2024-1, 2024-2, etc.)
- `search` (búsqueda por nombre)
- `country_id` (UUID del país)
- `limit` (default: 20)
- `offset` (default: 0)

---

## 📋 Checklist de Implementación

- [ ] Database connection configurada
- [ ] Model Student creado
- [ ] Repository Student con CRUD
- [ ] Service Student con validaciones
- [ ] Handler Student con endpoints
- [ ] Shared utilities (response helpers)
- [ ] Rutas registradas en main.go
- [ ] Unit tests para service
- [ ] Integration tests para repository
- [ ] API probada con curl/Postman
- [ ] Código documentado
- [ ] Error handling implementado

---

## 🚨 Importante

1. **Siempre usar context.Context** como primer parámetro
2. **Filtrar soft-deleted**: `WHERE deleted_at IS NULL`
3. **Usar placeholders SQL**: `$1, $2` (prevenir SQL injection)
4. **Validar en múltiples capas**: DTO validation + business rules
5. **Wrap errors**: `fmt.Errorf("context: %w", err)`
6. **Seguir convenciones Go**: nombres, estructura, estilo

---

## 💡 Tips para Claude Code

- Lee `/docs/agents/agente-go-backend.md` para ejemplos completos
- Sigue el patrón: Model → Repository → Service → Handler
- Usa los ejemplos del Student como template para otros módulos
- Pregunta cuando necesites aclaración sobre decisiones arquitectónicas
- Documenta código complejo
- Escribe tests mientras implementas (no al final)

---

## 🎓 Filosofía del Código

- **Simplicidad**: Código simple es código mantenible
- **Testeable**: Diseña para testing desde el inicio
- **Explícito**: Errores explícitos, no ocultos
- **Consistente**: Sigue patrones establecidos
- **Documentado**: Código complejo merece comentarios

---

**¡Listo para comenzar!** 🚀

Empieza con:
```
"Asume el rol de Agente Go/Backend. 
Lee /docs/agents/agente-go-backend.md.
Vamos a implementar el módulo de estudiantes."
```
