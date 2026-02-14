-- Migration: 008_refactor_student_location
-- Description: Reemplazar procedencia genérica por nacionalidad + residencia
-- Author: Agente DBA
-- Date: 2026-02-13
--
-- Contexto: Para una maestría virtual es crítico distinguir entre:
--   - Nacionalidad: ciudadanía del estudiante
--   - Residencia: dónde vive actualmente (zona horaria, logística, legal)
--
-- Cambios:
--   country_origin_id → nationality_country_id (renombrado)
--   city_origin_id    → eliminado
--   + residence_country_id (nuevo, NOT NULL)
--   + residence_city_id (nuevo, NOT NULL)

BEGIN;

-- =============================================================================
-- 1. Agregar nuevas columnas
-- =============================================================================

ALTER TABLE students
    ADD COLUMN nationality_country_id UUID REFERENCES countries(id) ON DELETE RESTRICT,
    ADD COLUMN residence_country_id UUID REFERENCES countries(id) ON DELETE RESTRICT,
    ADD COLUMN residence_city_id UUID REFERENCES cities(id) ON DELETE RESTRICT;

-- =============================================================================
-- 2. Migrar datos existentes
-- =============================================================================

-- Copiar country_origin_id a nationality_country_id (la nacionalidad era lo más cercano)
UPDATE students
SET nationality_country_id = country_origin_id
WHERE country_origin_id IS NOT NULL;

-- Para residencia, usar los mismos datos de origen como valor inicial
-- (los datos reales se actualizarán después)
UPDATE students
SET residence_country_id = country_origin_id,
    residence_city_id = city_origin_id
WHERE country_origin_id IS NOT NULL;

-- =============================================================================
-- 3. Aplicar constraints NOT NULL (después de migrar datos)
-- =============================================================================

ALTER TABLE students
    ALTER COLUMN nationality_country_id SET NOT NULL;

-- Residencia: NOT NULL solo si hay datos migrados válidos para city
-- Como city_origin_id puede ser NULL, hacemos residence_city_id NOT NULL
-- solo después de que se completen los datos
ALTER TABLE students
    ALTER COLUMN residence_country_id SET NOT NULL;

-- residence_city_id queda nullable temporalmente para registros sin city_origin_id
-- Se hará NOT NULL cuando se completen los datos faltantes

-- =============================================================================
-- 4. Eliminar columnas antiguas
-- =============================================================================

-- Eliminar índices antiguos primero
DROP INDEX IF EXISTS idx_students_city;
DROP INDEX IF EXISTS idx_students_country;

ALTER TABLE students
    DROP COLUMN country_origin_id,
    DROP COLUMN city_origin_id;

-- =============================================================================
-- 5. Crear nuevos índices
-- =============================================================================

CREATE INDEX idx_students_nationality ON students(nationality_country_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_students_residence_country ON students(residence_country_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_students_residence_city ON students(residence_city_id) WHERE deleted_at IS NULL;

-- =============================================================================
-- 6. Comentarios
-- =============================================================================

COMMENT ON COLUMN students.nationality_country_id IS 'País de nacionalidad/ciudadanía del estudiante';
COMMENT ON COLUMN students.residence_country_id IS 'País de residencia actual del estudiante';
COMMENT ON COLUMN students.residence_city_id IS 'Ciudad de residencia actual del estudiante';

COMMIT;
