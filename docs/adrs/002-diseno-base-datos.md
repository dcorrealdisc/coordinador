# ADR-002: Diseño de Base de Datos y Estrategia de Datos

**Estado**: Aceptado  
**Fecha**: 2026-02-13  
**Decisor**: Dario Correal  
**Agente**: DBA

## Contexto y Problema

El sistema Coordinador requiere una base de datos que soporte:

1. **Gestión académica completa**: Estudiantes, cursos, inscripciones, calificaciones
2. **Planificación**: Períodos académicos, programación de cursos, asignación de recursos
3. **Gestión de personal**: Profesores, tutores con sus intereses y asignaciones
4. **Reportes complejos**: Progreso estudiantil, estadísticas de cursos, analítica histórica
5. **Auditoría**: Rastreo completo de quién modificó qué y cuándo
6. **Datos consolidados**: Universidades, empresas, ubicaciones geográficas

### Requerimientos No Funcionales

- Escala: 100-500 estudiantes activos, <100 usuarios concurrentes
- Performance: Reportes complejos <3 segundos
- Integridad: ACID para datos transaccionales
- Auditoría: Trazabilidad completa de cambios
- Evolución: Diseño flexible para cambios futuros

## Decisión

### 1. PostgreSQL como RDBMS

**Decisión**: Usar PostgreSQL 15+ como base de datos principal.

**Justificación**:
- ✅ Relaciones complejas entre entidades académicas
- ✅ Soporte nativo para JSONB (flexibilidad futura)
- ✅ Arrays nativos (emails[], phones[])
- ✅ Vistas materializadas para reportes
- ✅ Triggers y funciones almacenadas
- ✅ Excelente performance con índices apropiados
- ✅ Open source, maduro, ampliamente adoptado
- ✅ Se alinea con ADR-001 (arquitectura general)

### 2. UUIDs como Identificadores Primarios

**Decisión**: Usar UUIDs (gen_random_uuid()) en lugar de SERIAL/BIGSERIAL.

**Justificación**:
- ✅ Globalmente únicos (facilita merge de datos)
- ✅ No predecibles (seguridad)
- ✅ Permiten generación distribuida
- ✅ Facilitan testing (IDs reproducibles)
- ⚠️ Ligeramente más lentos que INT (aceptable para nuestra escala)

**Alternativa rechazada**: SERIAL/BIGSERIAL
- Más rápidos pero predecibles y no globalmente únicos

### 3. Soft Delete para Datos de Negocio

**Decisión**: Eliminación lógica (soft delete) usando campo `deleted_at`.

**Justificación**:
- ✅ Auditoría completa (nada se pierde)
- ✅ Recuperación de datos posible
- ✅ Análisis histórico completo
- ✅ Cumplimiento regulatorio
- ⚠️ Queries deben filtrar `deleted_at IS NULL`

**Implementación**:
```sql
deleted_at TIMESTAMP,
deleted_by UUID REFERENCES system_users(id)

-- Índices excluyen eliminados
CREATE INDEX idx_table_field 
ON table(field) 
WHERE deleted_at IS NULL;
```

### 4. Auditoría Completa en Todas las Tablas Principales

**Decisión**: Campos de auditoría estándar en todas las tablas de negocio.

**Campos**:
```sql
created_at TIMESTAMP DEFAULT NOW(),
created_by UUID REFERENCES system_users(id),
updated_at TIMESTAMP DEFAULT NOW(),
updated_by UUID REFERENCES system_users(id),
deleted_at TIMESTAMP,
deleted_by UUID REFERENCES system_users(id)
```

**Justificación**:
- ✅ Trazabilidad completa (quién, qué, cuándo)
- ✅ Debugging facilitado
- ✅ Análisis de uso
- ✅ Cumplimiento normativo
- ✅ Triggers automatizan `updated_at`

### 5. Catálogos Maestros para Consolidación

**Decisión**: Tablas separadas para entidades consolidables.

**Implementación**:
- `countries` - Catálogo de países (ISO codes)
- `cities` - Catálogo de ciudades
- `universities` - Universidades consolidadas
- `companies` - Empresas consolidadas

**Justificación**:
- ✅ Datos limpios y consistentes
- ✅ Reportes precisos por país/ciudad/universidad
- ✅ Autocompletado en UIs
- ✅ Análisis demográfico confiable

**Alternativa rechazada**: Texto libre
- Genera duplicados ("Universidad de los Andes" vs "Uniandes")
- Reportes imprecisos

### 6. Arrays PostgreSQL para Listas Simples

**Decisión**: Usar arrays (TEXT[]) para emails y teléfonos.

**Justificación**:
- ✅ Más simple que tabla de unión
- ✅ Queries eficientes
- ✅ Suficiente para listas sin relaciones complejas

**Cuándo usar arrays**:
```sql
emails TEXT[],  -- ✅ Lista simple, sin metadata
phones TEXT[]   -- ✅ Lista simple, sin metadata
```

**Cuándo NO usar arrays** (usar tabla de unión):
- Relaciones complejas (ej: student_universities con degree_obtained)
- Queries frecuentes sobre elementos individuales
- Constraints sobre elementos individuales

### 7. CQRS Light con Vistas Materializadas

**Decisión**: Vistas materializadas para reportes complejos (Read Path).

**Vistas creadas**:
1. `student_academic_progress` - Progreso por estudiante
2. `course_period_statistics` - Estadísticas por curso/período
3. `students_by_location` - Distribución geográfica
4. `students_by_university` - Por universidad de procedencia
5. `students_by_company` - Por empleador
6. `tutor_workload` - Carga de tutores
7. `students_age_distribution` - Rangos de edad

**Justificación**:
- ✅ Queries complejos pre-computados
- ✅ Response time <1s para dashboards
- ✅ Se alinea con ADR-001 (CQRS Light)
- ✅ Refresh programado (cada hora aceptable)

**Estrategia de refresh**:
```sql
-- Manual (desarrollo)
REFRESH MATERIALIZED VIEW CONCURRENTLY view_name;

-- Programado (producción con pg_cron)
SELECT cron.schedule('job', '0 * * * *', 'REFRESH ...');

-- On-demand (post operaciones críticas)
CALL refresh_all_materialized_views();
```

### 8. Constraints en Base de Datos

**Decisión**: Validación de integridad en la base de datos, no solo en aplicación.

**Implementación**:
```sql
-- NOT NULL para campos obligatorios
full_name VARCHAR(255) NOT NULL

-- CHECK constraints para validación
status CHECK (status IN ('active', 'graduated', 'withdrawn'))
final_grade CHECK (final_grade >= 0 AND final_grade <= 5)
birth_date CHECK (birth_date < CURRENT_DATE)

-- UNIQUE constraints
UNIQUE(student_id, scheduled_course_id)

-- Foreign keys con políticas claras
REFERENCES table(id) ON DELETE CASCADE  -- Dependientes
REFERENCES table(id) ON DELETE RESTRICT -- Independientes
REFERENCES table(id) ON DELETE SET NULL -- Opcionales
```

**Justificación**:
- ✅ Base de datos como última línea de defensa
- ✅ Integridad garantizada independiente de la aplicación
- ✅ Previene estados inválidos
- ✅ Self-documenting (schema describe reglas)

### 9. Índices Estratégicos

**Decisión**: Índices basados en patrones de acceso reales.

**Estrategias**:
```sql
-- Índices en FKs (siempre)
CREATE INDEX idx_students_country ON students(country_origin_id);

-- Índices parciales (solo registros activos)
CREATE INDEX idx_students_status 
ON students(status) 
WHERE deleted_at IS NULL;

-- Índices compuestos (queries con múltiples filtros)
CREATE INDEX idx_enrollments_student_status 
ON enrollments(student_id, status);

-- Índices únicos en vistas materializadas
CREATE UNIQUE INDEX idx_view_pk ON view(id);
```

**Principio**: No sobre-indexar, solo donde mejora performance medible.

### 10. Migraciones Versionadas

**Decisión**: golang-migrate para migraciones incrementales.

**Estructura**:
```
migrations/
├── 001_create_base_schema.sql
├── 002_create_academic_tables.sql
├── 003_create_people_tables.sql
├── 004_create_enrollments.sql
├── 005_create_materialized_views.sql
└── 006_create_functions_triggers.sql
```

**Justificación**:
- ✅ Versionado claro y ordenado
- ✅ Reproducible en cualquier ambiente
- ✅ Rollback posible
- ✅ Sincronización entre desarrolladores
- ✅ Historial de cambios documentado

## Consecuencias

### Positivas

- ✅ Diseño normalizado (3NF) previene inconsistencias
- ✅ Vistas materializadas optimizan reportes sin duplicar código
- ✅ Auditoría completa facilita debugging y cumplimiento
- ✅ Soft delete preserva historial completo
- ✅ Catálogos consolidados generan reportes precisos
- ✅ UUIDs facilitan testing y evolución futura
- ✅ Constraints en BD garantizan integridad
- ✅ Migraciones permiten evolución controlada

### Negativas

- ⚠️ Queries deben recordar filtrar `deleted_at IS NULL`
- ⚠️ Vistas materializadas no son tiempo real (refresh cada hora)
- ⚠️ UUIDs ligeramente menos performantes que INTs
- ⚠️ Campos de auditoría agregan overhead en todas las tablas

### Riesgos y Mitigaciones

**Riesgo**: Vistas materializadas desactualizadas
- **Mitigación**: Refresh automático cada hora + refresh manual post-operaciones críticas
- **Monitoreo**: Timestamp de último refresh visible

**Riesgo**: Over-indexing degrada performance de writes
- **Mitigación**: Revisar pg_stat_user_indexes periódicamente
- **Principio**: Solo indexar campos realmente consultados

**Riesgo**: Soft delete causa confusión (registros "fantasma")
- **Mitigación**: Índices parciales excluyen eliminados
- **Convención**: Siempre filtrar `WHERE deleted_at IS NULL` en queries

**Riesgo**: Falta de constraints permite datos inválidos
- **Mitigación**: Constraints completos en schema + tests de integridad
- **Validación dual**: Frontend valida UX + BD valida integridad

## Alternativas Consideradas

### NoSQL (MongoDB, DynamoDB)

**Pros**: Flexibilidad de schema, escalamiento horizontal
**Contras**: 
- Relaciones complejas son difíciles
- No ACID multi-documento
- Reportes complejos requieren procesamiento adicional

**Decisión**: Rechazado - Datos altamente relacionales, ACID necesario

### MySQL/MariaDB

**Pros**: Ampliamente usado, performance excelente
**Contras**:
- Vistas materializadas no nativas
- Arrays no soportados
- JSONB inferior

**Decisión**: Rechazado - PostgreSQL superior para nuestro caso de uso

### Hard Delete

**Pros**: Datos más limpios, queries más simples
**Contras**: 
- Pérdida permanente de datos
- No auditoría
- No análisis histórico

**Decisión**: Rechazado - Auditoría y análisis histórico son críticos

### SERIAL en lugar de UUID

**Pros**: Más rápido, menos espacio
**Contras**:
- Predecibles (seguridad)
- No globalmente únicos
- Problemas al merge de datos

**Decisión**: Rechazado - Beneficios de UUIDs superan costo de performance

## Métricas de Éxito

- **Performance**: Reportes complejos <3s (p95)
- **Integridad**: 0 estados inválidos en producción
- **Auditoría**: 100% de cambios rastreables
- **Disponibilidad**: Migraciones zero-downtime
- **Mantenibilidad**: Nuevos campos agregables sin downtime

## Referencias

- [PostgreSQL Documentation](https://www.postgresql.org/docs/15/)
- [UUID Performance in PostgreSQL](https://www.2ndquadrant.com/en/blog/sequential-uuid-generators/)
- [Materialized Views Best Practices](https://www.postgresql.org/docs/current/rules-materializedviews.html)
- [Database Normalization](https://en.wikipedia.org/wiki/Database_normalization)
- [golang-migrate](https://github.com/golang-migrate/migrate)

## Notas

Este diseño está optimizado para:
- Escala actual (100-500 estudiantes)
- Patrones de uso (70% read, 30% write)
- Necesidad de reportes complejos

Si la escala crece significativamente (>5000 estudiantes), considerar:
- Particionamiento de tablas grandes
- Read replicas para reportes
- Caché distribuido (Redis) para queries frecuentes
- Separación física de write/read databases

El diseño modular facilita estas evoluciones sin reescritura completa.

---

**Última actualización**: 2026-02-13
