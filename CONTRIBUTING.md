# Guía de Contribución - Coordinador

Bienvenido al proyecto Coordinador. Este documento describe cómo contribuir al proyecto siguiendo el enfoque de desarrollo basado en agentes especializados.

## 🎯 Filosofía del Proyecto

Este proyecto sigue un enfoque único de **desarrollo asistido por agentes de IA**, donde cada aspecto técnico es manejado por un agente especializado que mantiene contexto y coherencia en su dominio.

## 🤖 Desarrollo Basado en Agentes

### Agentes Disponibles

Cada agente tiene expertise específico y responsabilidades claras:

| Agente | Dominio | Guía |
|--------|---------|------|
| **Arquitecto** | Decisiones de diseño, ADRs, coherencia arquitectónica | [Ver guía](./docs/agents/agente-arquitecto.md) |
| **DBA** | Modelado de datos, optimizaciones PostgreSQL | _Próximamente_ |
| **Go/Backend** | Implementación backend, APIs, lógica de negocio | _Próximamente_ |
| **Svelte** | Interfaces de usuario, componentes | _Próximamente_ |
| **DevOps** | CI/CD, containerización, deployment | _Próximamente_ |

### Workflow con Agentes

```
1. Identificar tipo de tarea
   ↓
2. Consultar al agente correspondiente
   ↓
3. Agente revisa contexto (ADRs, código existente)
   ↓
4. Agente propone solución
   ↓
5. Revisar y aprobar
   ↓
6. Documentar decisión (si aplica)
   ↓
7. Implementar
```

## 📋 Proceso de Desarrollo

### 1. Antes de Empezar

- [ ] Lee el [README principal](./README.md)
- [ ] Revisa [ADR-001: Arquitectura General](./docs/adrs/001-arquitectura-general.md)
- [ ] Familiarízate con la estructura del proyecto
- [ ] Identifica qué agente necesitas consultar

### 2. Para Nuevas Features

1. **Consulta al Agente Arquitecto** primero
   - ¿Esta feature afecta la arquitectura?
   - ¿Necesitamos un nuevo módulo?
   - ¿Hay patrones existentes que debamos seguir?

2. **Define el diseño** con el agente correspondiente
   - DBA: Si involucra cambios en BD
   - Go/Backend: Si es lógica de negocio
   - Svelte: Si es UI

3. **Documenta** si la decisión es significativa
   - Crear ADR si es decisión arquitectónica
   - Actualizar documentación del módulo

4. **Implementa** siguiendo el diseño aprobado

5. **Prueba** con tests adecuados

### 3. Para Bug Fixes

1. Identifica el módulo afectado
2. Consulta al agente del módulo
3. Propón fix
4. Agrega test que reproduzca el bug
5. Implementa y verifica

### 4. Para Refactoring

1. **Consulta al Agente Arquitecto**
2. Justifica el refactoring (deuda técnica, performance, etc.)
3. Documenta en ADR si es significativo
4. Asegura que tests existentes pasen
5. Implementa incrementalmente

## 🏗️ Estándares de Código

### Backend (Go)

```go
// Estructura de archivos por módulo
module_name/
├── handler.go      // HTTP handlers
├── service.go      // Lógica de negocio
├── repository.go   // Acceso a datos
├── models.go       // Estructuras de dominio
├── dto.go          // Data Transfer Objects
└── *_test.go       // Tests
```

**Convenciones:**
- Nombres exportados: `PascalCase`
- Nombres privados: `camelCase`
- Errores: retornar siempre como último valor
- Tests: cobertura mínima 70%
- Comentarios: godoc style para exports
- Formato: `go fmt` antes de commit

### Frontend (Svelte)

```
component/
├── ComponentName.svelte
├── ComponentName.test.ts
└── index.ts
```

**Convenciones:**
- Componentes: `PascalCase.svelte`
- Stores: `camelCase.ts`
- Types: `types.ts` por módulo
- Estilos: TailwindCSS utility-first
- Tests: Vitest
- Formato: Prettier antes de commit

### SQL (PostgreSQL)

**Convenciones:**
- Tablas: `snake_case` plural
- Columnas: `snake_case`
- Índices: `idx_table_column`
- FK: `fk_table_referenced`
- Migraciones: timestamp-based

## 📝 Documentación

### Cuándo Crear un ADR

Crea un ADR (Architectural Decision Record) cuando:

- ✅ Cambias patrones arquitectónicos
- ✅ Introduces nueva tecnología/librería significativa
- ✅ Modificas estructura de módulos
- ✅ Cambias estrategia de datos (ej: agregar caché)
- ✅ Tomas decisión que afecta múltiples módulos

### Formato de ADR

Ver [template en Agente Arquitecto](./docs/agents/agente-arquitecto.md#template-de-adr)

### Documentación de Código

- Backend: Comentarios godoc para exports
- Frontend: JSDoc para funciones públicas
- APIs: OpenAPI/Swagger (próximamente)

## 🧪 Testing

### Niveles de Testing

1. **Unit Tests**: Cada función/método
2. **Integration Tests**: Interacción entre módulos
3. **E2E Tests**: Flujos completos (próximamente)

### Cobertura Mínima

- Backend: 70%
- Frontend: 60% (componentes críticos)

### Ejecutar Tests

```bash
# Todos
make test

# Backend solo
make test-backend

# Con cobertura
make test-backend-coverage

# Admin web
make test-admin

# Portal web
make test-portal
```

## 🔀 Git Workflow

### Branches

```
main              # Producción
├── develop       # Desarrollo
    ├── feature/* # Nuevas features
    ├── fix/*     # Bug fixes
    └── refactor/* # Refactorings
```

### Commits

Seguimos [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: agregar endpoint de reportes de estudiantes
fix: corregir cálculo de promedio ponderado
docs: actualizar ADR-001 con decisión de caché
refactor: extraer lógica de validación a servicio
test: agregar tests para módulo de cursos
chore: actualizar dependencias de Go
```

**Tipos:**
- `feat`: Nueva funcionalidad
- `fix`: Corrección de bug
- `docs`: Solo documentación
- `refactor`: Refactoring (sin cambio funcional)
- `test`: Agregar o modificar tests
- `chore`: Mantenimiento (deps, configs)
- `perf`: Mejoras de performance

### Pull Requests

1. **Título**: Descriptivo y claro
2. **Descripción**: 
   - ¿Qué problema resuelve?
   - ¿Cómo lo resuelve?
   - ¿Qué agente consultaste?
3. **Checklist**:
   - [ ] Tests agregados/actualizados
   - [ ] Documentación actualizada
   - [ ] ADR creado si aplica
   - [ ] Código formateado
   - [ ] Linter pasa

## 🔧 Setup de Desarrollo

### Prerrequisitos

- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Docker & Docker Compose (opcional)
- Make (opcional, pero recomendado)

### Primera vez

```bash
# Clonar repo
git clone <url>
cd coordinador

# Instalar dependencias
make install

# Levantar PostgreSQL
make db-up

# Iniciar servicios (en terminales separadas)
make dev-backend
make dev-admin
make dev-portal
```

### Variables de Entorno

Crear archivos `.env`:

**Backend** (`backend/.env`):
```env
PORT=8080
ENV=development
DB_HOST=localhost
DB_PORT=5432
DB_USER=coordinador
DB_PASSWORD=coordinador_dev_2024
DB_NAME=coordinador_db
DB_SSLMODE=disable
```

**Admin Web** (`admin-web/.env.local`):
```env
VITE_API_URL=http://localhost:8080/api/v1
```

**Portal Web** (`portal-web/.env.local`):
```env
VITE_API_URL=http://localhost:8080/api/v1
```

## 🚨 Señales de Alerta

No hagas commit si:

- ❌ Tests fallan
- ❌ Linter tiene errores
- ❌ Código no está formateado
- ❌ Falta documentación obligatoria
- ❌ Rompe la arquitectura sin ADR que lo justifique

## 💡 Tips

### Para Consultar Agentes

Cuando consultes a un agente (vía Claude AI):

1. **Proporciona contexto**:
   ```
   "Asume el rol del Agente Arquitecto.
   Quiero agregar caché para reportes.
   Contexto: [pegar info relevante]
   ADRs relacionados: ADR-001"
   ```

2. **Sé específico**:
   - ❌ "¿Cómo hago reportes?"
   - ✅ "Necesito cachear el reporte de desempeño estudiantil que se ejecuta cada hora. ¿Redis o vistas materializadas?"

3. **Referencia decisiones previas**:
   - Menciona ADRs
   - Señala código existente
   - Indica qué has intentado

### Para Resolver Conflictos

1. Revisa ADRs relacionados
2. Consulta al Agente Arquitecto
3. Si hay ambigüedad, documenta ambas opciones y sus trade-offs
4. Discute con el equipo/proyecto owner

## 📚 Recursos

- [README del Proyecto](./README.md)
- [Guía del Agente Arquitecto](./docs/agents/agente-arquitecto.md)
- [Índice de ADRs](./docs/adrs/README.md)
- [Go Best Practices](https://go.dev/doc/effective_go)
- [Svelte Documentation](https://svelte.dev/docs)

## 📞 Preguntas

Si tienes dudas:

1. Revisa documentación existente
2. Consulta al agente correspondiente
3. Revisa ADRs relacionados
4. Contacta al maintainer del proyecto

---

**¡Gracias por contribuir a Coordinador!** 🎓
