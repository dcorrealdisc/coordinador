# Migraciones de Base de Datos

Migraciones SQL para el sistema Coordinador usando PostgreSQL 15+.

## 📋 Migraciones Disponibles

| # | Archivo | Descripción | Estado |
|---|---------|-------------|--------|
| 001 | `create_base_schema.sql` | Catálogos maestros y configuración | ✅ Listo |
| 002 | `create_academic_tables.sql` | Cursos y períodos académicos | ✅ Listo |
| 003 | `create_people_tables.sql` | Estudiantes, profesores, tutores | ✅ Listo |
| 004 | `create_enrollments.sql` | Inscripciones y calificaciones | ✅ Listo |
| 005 | `create_materialized_views.sql` | Vistas para reportes (CQRS) | ✅ Listo |
| 006 | `create_functions_triggers.sql` | Funciones y triggers automáticos | ✅ Listo |

## 🚀 Aplicar Migraciones

### Opción 1: Manualmente con psql (Desarrollo)

```bash
# 1. Asegurarse que PostgreSQL está corriendo
docker-compose up -d postgres

# 2. Aplicar migraciones en orden
psql -U coordinador -d coordinador_db -f backend/migrations/001_create_base_schema.sql
psql -U coordinador -d coordinador_db -f backend/migrations/002_create_academic_tables.sql
psql -U coordinador -d coordinador_db -f backend/migrations/003_create_people_tables.sql
psql -U coordinador -d coordinador_db -f backend/migrations/004_create_enrollments.sql
psql -U coordinador -d coordinador_db -f backend/migrations/005_create_materialized_views.sql
psql -U coordinador -d coordinador_db -f backend/migrations/006_create_functions_triggers.sql

# 3. Verificar tablas creadas
psql -U coordinador -d coordinador_db -c "\dt"
```

### Opción 2: Script de Aplicación Rápida

```bash
# Crear script helper
cat > backend/migrations/apply_all.sh << 'EOF'
#!/bin/bash
set -e

DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-coordinador}
DB_NAME=${DB_NAME:-coordinador_db}

echo "🔄 Aplicando migraciones a $DB_NAME..."

for migration in backend/migrations/*.sql; do
    echo "  ⚙️  Aplicando $(basename $migration)..."
    PGPASSWORD=${DB_PASSWORD} psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$migration"
done

echo "✅ Migraciones aplicadas exitosamente"
EOF

chmod +x backend/migrations/apply_all.sh

# Ejecutar
./backend/migrations/apply_all.sh
```

### Opción 3: Con golang-migrate (Recomendado para Producción)

```bash
# 1. Instalar golang-migrate
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# 2. Aplicar migraciones
migrate -path backend/migrations -database "postgresql://coordinador:coordinador_dev_2024@localhost:5432/coordinador_db?sslmode=disable" up

# 3. Ver versión actual
migrate -path backend/migrations -database "postgresql://..." version

# 4. Rollback (si es necesario)
migrate -path backend/migrations -database "postgresql://..." down 1
```

## 🔍 Verificar el Schema

```bash
# Conectar a PostgreSQL
psql -U coordinador -d coordinador_db

# Ver todas las tablas
\dt

# Ver schema de una tabla
\d students

# Ver índices
\di

# Ver vistas materializadas
\dm

# Ver funciones
\df

# Salir
\q
```

## 📊 Estructura del Schema

### Catálogos Maestros
- `countries` - Países (ISO codes)
- `cities` - Ciudades
- `universities` - Universidades consolidadas
- `companies` - Empresas consolidadas

### Sistema
- `system_users` - Usuarios administrativos (para auditoría)
- `program_configuration` - Configuración del programa

### Académico
- `courses` - Catálogo de cursos
- `course_prerequisites` - Prerrequisitos entre cursos
- `academic_periods` - Períodos académicos
- `scheduled_courses` - Cursos programados por período

### Personas
- `students` - Estudiantes (con auditoría completa)
- `student_universities` - Universidades de procedencia
- `professors` - Profesores
- `course_professor_assignments` - Asignación profesor-curso
- `tutors` - Tutores/monitores
- `tutor_course_interests` - Intereses de tutores
- `course_tutor_assignments` - Asignación tutor-curso

### Inscripciones
- `enrollments` - Inscripciones y calificaciones

### Vistas Materializadas (Reportes)
- `student_academic_progress` - Progreso por estudiante
- `course_period_statistics` - Estadísticas por curso
- `students_by_location` - Distribución geográfica
- `students_by_university` - Por universidad
- `students_by_company` - Por empleador
- `tutor_workload` - Carga de tutores
- `students_age_distribution` - Rangos de edad

## 🔄 Refresh de Vistas Materializadas

```sql
-- Manual - Una vista específica
REFRESH MATERIALIZED VIEW CONCURRENTLY student_academic_progress;

-- Manual - Todas las vistas
SELECT refresh_all_materialized_views();

-- Programado con pg_cron (configurar en producción)
SELECT cron.schedule(
    'refresh-reports',
    '0 * * * *',  -- Cada hora
    'SELECT refresh_all_materialized_views()'
);
```

## 🧪 Testing del Schema

```sql
-- Test 1: Verificar que todas las tablas existen
SELECT count(*) FROM information_schema.tables 
WHERE table_schema = 'public' AND table_type = 'BASE TABLE';
-- Esperado: 15 tablas

-- Test 2: Verificar vistas materializadas
SELECT count(*) FROM pg_matviews;
-- Esperado: 7 vistas

-- Test 3: Verificar triggers
SELECT count(*) FROM information_schema.triggers 
WHERE trigger_schema = 'public';
-- Esperado: 15+ triggers

-- Test 4: Verificar funciones
SELECT count(*) FROM pg_proc 
WHERE pronamespace = 'public'::regnamespace;
-- Esperado: 6+ funciones

-- Test 5: Insertar datos de prueba
BEGIN;

-- Insertar estudiante de prueba
INSERT INTO students (full_name, document_id, birth_date, emails, country_origin_id, status, cohort, enrollment_date)
SELECT 'Juan Test', '123456', '1995-01-01', ARRAY['juan@test.com'], id, 'active', '2024-1', '2024-01-15'
FROM countries WHERE code = 'COL' LIMIT 1;

-- Verificar
SELECT * FROM students WHERE document_id = '123456';

ROLLBACK;  -- No guardar datos de prueba
```

## 📝 Agregar Nueva Migración

```bash
# 1. Crear archivo con número secuencial
touch backend/migrations/007_nueva_feature.sql

# 2. Seguir template
cat > backend/migrations/007_nueva_feature.sql << 'EOF'
-- Migration: 007_nueva_feature
-- Description: Descripción de la migración
-- Author: Agente DBA / Tu Nombre
-- Date: YYYY-MM-DD

BEGIN;

-- Código SQL aquí
ALTER TABLE students ADD COLUMN nueva_columna TEXT;

COMMIT;

-- Rollback (comentado)
-- BEGIN;
-- ALTER TABLE students DROP COLUMN nueva_columna;
-- COMMIT;
EOF

# 3. Aplicar
psql -U coordinador -d coordinador_db -f backend/migrations/007_nueva_feature.sql
```

## 🚨 Troubleshooting

### Error: "relation already exists"

```bash
# Las migraciones ya fueron aplicadas
# Verificar estado:
psql -U coordinador -d coordinador_db -c "\dt"

# Si necesitas empezar de cero (¡CUIDADO! Borra todo):
psql -U coordinador -d coordinador_db -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
```

### Error: "database does not exist"

```bash
# Crear base de datos
createdb -U coordinador coordinador_db

# O con psql:
psql -U coordinador -c "CREATE DATABASE coordinador_db;"
```

### Error: "permission denied"

```bash
# Asegurar permisos correctos
psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE coordinador_db TO coordinador;"
```

### Verificar Password de PostgreSQL

```bash
# Si usas Docker Compose
docker-compose exec postgres psql -U coordinador -d coordinador_db

# Password está en docker-compose.yml:
# POSTGRES_PASSWORD: coordinador_dev_2024
```

## 📚 Recursos

- [PostgreSQL Docs](https://www.postgresql.org/docs/15/)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [ADR-002: Diseño de BD](../../docs/adrs/002-diseno-base-datos.md)
- [Guía Agente DBA](../../docs/agents/agente-dba.md)

## 🎯 Checklist Post-Migración

Después de aplicar migraciones en un ambiente nuevo:

- [ ] Todas las tablas creadas (`\dt` muestra 15 tablas)
- [ ] Todas las vistas materializadas creadas (`\dm` muestra 7 vistas)
- [ ] Triggers aplicados correctamente
- [ ] Datos iniciales cargados (países, configuración)
- [ ] Vistas materializadas tienen datos (`SELECT count(*) FROM student_academic_progress`)
- [ ] Tests básicos pasan
- [ ] Performance aceptable (verificar EXPLAIN ANALYZE en queries clave)

---

**Última actualización**: 2026-02-13  
**Versión del schema**: 006
