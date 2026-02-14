-- Migration: 009_refactor_student_name_and_job
-- Description: Separar full_name, catálogo de cargos y catálogo de profesiones
-- Author: Agente DBA
-- Date: 2026-02-13
--
-- Cambios:
--   1. full_name → first_names + last_names
--   2. Nueva tabla job_title_categories (catálogo para análisis)
--   3. students.job_title_category_id FK (opcional)
--   4. Nueva tabla professions (catálogo para análisis)
--   5. students.profession_id FK (opcional)

BEGIN;

-- =============================================================================
-- 1. SEPARAR NOMBRE COMPLETO
-- =============================================================================

ALTER TABLE students
    ADD COLUMN first_names VARCHAR(150),
    ADD COLUMN last_names VARCHAR(150);

-- Migrar datos existentes: intentar separar por el primer espacio
-- Si no hay espacio, todo va a first_names
UPDATE students
SET first_names = CASE
        WHEN position(' ' IN full_name) > 0
        THEN left(full_name, position(' ' IN full_name) - 1)
        ELSE full_name
    END,
    last_names = CASE
        WHEN position(' ' IN full_name) > 0
        THEN substring(full_name FROM position(' ' IN full_name) + 1)
        ELSE ''
    END
WHERE full_name IS NOT NULL;

-- Aplicar NOT NULL después de migrar
ALTER TABLE students
    ALTER COLUMN first_names SET NOT NULL,
    ALTER COLUMN last_names SET NOT NULL;

-- Eliminar columna vieja
ALTER TABLE students
    DROP COLUMN full_name;

COMMENT ON COLUMN students.first_names IS 'Nombres del estudiante';
COMMENT ON COLUMN students.last_names IS 'Apellidos del estudiante';

-- =============================================================================
-- 2. CATÁLOGO DE CATEGORÍAS DE CARGO
-- =============================================================================

CREATE TABLE job_title_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT uk_job_title_category_name UNIQUE(name)
);

COMMENT ON TABLE job_title_categories IS 'Catálogo de categorías de cargo para análisis de perfil laboral de estudiantes';

CREATE TRIGGER trigger_job_title_categories_updated_at
    BEFORE UPDATE ON job_title_categories
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Seed de categorías iniciales
INSERT INTO job_title_categories (name, description) VALUES
    ('Director/a', 'Dirección general, CEO, CTO, CFO, etc.'),
    ('Gerente', 'Gerencia de área o departamento'),
    ('Subgerente', 'Subgerencia o dirección adjunta'),
    ('Coordinador/a', 'Coordinación de equipos o proyectos'),
    ('Jefe de área', 'Jefatura de área o sección'),
    ('Líder técnico', 'Liderazgo técnico de equipos'),
    ('Consultor/a', 'Consultoría interna o externa'),
    ('Analista', 'Análisis de datos, negocio, sistemas, etc.'),
    ('Desarrollador/a', 'Desarrollo de software, ingeniería'),
    ('Diseñador/a', 'Diseño UX/UI, gráfico, industrial'),
    ('Investigador/a', 'Investigación académica o corporativa'),
    ('Docente', 'Docencia universitaria o educación'),
    ('Especialista', 'Especialista en un área técnica'),
    ('Profesional independiente', 'Freelance, emprendedor, consultor independiente'),
    ('Otro', 'Categoría no listada');

-- =============================================================================
-- 3. CATÁLOGO DE PROFESIONES
-- =============================================================================

CREATE TABLE professions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(150) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT uk_profession_name UNIQUE(name)
);

COMMENT ON TABLE professions IS 'Catálogo de profesiones/carreras de pregrado de los estudiantes';

CREATE TRIGGER trigger_professions_updated_at
    BEFORE UPDATE ON professions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Seed de profesiones comunes en maestrías de tecnología/gestión
INSERT INTO professions (name) VALUES
    ('Ingeniería de Sistemas'),
    ('Ingeniería de Software'),
    ('Ingeniería Industrial'),
    ('Ingeniería Electrónica'),
    ('Ingeniería de Telecomunicaciones'),
    ('Ingeniería Mecánica'),
    ('Ingeniería Civil'),
    ('Administración de Empresas'),
    ('Economía'),
    ('Contaduría Pública'),
    ('Matemáticas'),
    ('Estadística'),
    ('Física'),
    ('Ciencias de la Computación'),
    ('Diseño Industrial'),
    ('Comunicación Social'),
    ('Derecho'),
    ('Psicología'),
    ('Otra');

-- =============================================================================
-- 4. AGREGAR FKs A STUDENTS
-- =============================================================================

ALTER TABLE students
    ADD COLUMN job_title_category_id UUID REFERENCES job_title_categories(id) ON DELETE SET NULL,
    ADD COLUMN profession_id UUID REFERENCES professions(id) ON DELETE SET NULL;

CREATE INDEX idx_students_job_title_category ON students(job_title_category_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_students_profession ON students(profession_id) WHERE deleted_at IS NULL;

COMMENT ON COLUMN students.job_title_category_id IS 'Categoría de cargo del estudiante en su empresa (para análisis)';
COMMENT ON COLUMN students.profession_id IS 'Profesión/carrera de pregrado del estudiante';

COMMIT;
