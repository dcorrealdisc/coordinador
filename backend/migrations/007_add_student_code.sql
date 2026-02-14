-- Migration: 007_add_student_code
-- Description: Agrega código de estudiante universitario (asignado por la universidad)
-- Author: Agente DBA
-- Date: 2026-02-13

BEGIN;

-- =============================================================================
-- CÓDIGO DE ESTUDIANTE
-- Formato: YYYYS#### (9 caracteres)
--   YYYY = año (4 dígitos)
--   S    = semestre (1 o 2)
--   #### = secuencia (4 dígitos)
-- Ejemplo: 202620190 = año 2026, semestre 2, secuencia 0190
-- Es asignado por la universidad, NO autogenerado.
-- Es el vínculo principal entre estudiantes y cursos.
-- =============================================================================

ALTER TABLE students
    ADD COLUMN student_code VARCHAR(9) UNIQUE;

-- Validar formato: exactamente 9 dígitos, 5to dígito solo 1 o 2
ALTER TABLE students
    ADD CONSTRAINT chk_student_code_format
    CHECK (student_code ~ '^[0-9]{4}[12][0-9]{4}$');

-- Índice para búsquedas por código (ya tiene UNIQUE, pero explícito para queries con soft delete)
CREATE INDEX idx_students_student_code ON students(student_code) WHERE deleted_at IS NULL;

COMMENT ON COLUMN students.student_code IS 'Código universitario del estudiante. Formato YYYYS#### (año-semestre-secuencia). Asignado por la universidad.';

COMMIT;
