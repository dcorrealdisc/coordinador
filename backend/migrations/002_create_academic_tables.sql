-- Migration: 002_create_academic_tables
-- Description: Cursos, períodos académicos y programación
-- Author: Agente DBA
-- Date: 2024-02-13

BEGIN;

-- =============================================================================
-- CURSOS
-- =============================================================================

CREATE TABLE courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    credits INTEGER NOT NULL CHECK (credits > 0),
    course_type VARCHAR(20) NOT NULL CHECK (course_type IN ('required', 'elective')),
    description TEXT,
    syllabus_url TEXT,
    is_active BOOLEAN DEFAULT true,
    
    -- Auditoría
    created_at TIMESTAMP DEFAULT NOW(),
    created_by UUID REFERENCES system_users(id) ON DELETE SET NULL,
    updated_at TIMESTAMP DEFAULT NOW(),
    updated_by UUID REFERENCES system_users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP,
    deleted_by UUID REFERENCES system_users(id) ON DELETE SET NULL
);

COMMENT ON TABLE courses IS 'Catálogo de cursos de la maestría';
COMMENT ON COLUMN courses.code IS 'Código único del curso (ej: MATE-101, DATA-201)';
COMMENT ON COLUMN courses.course_type IS 'Tipo: required (obligatorio) o elective (electivo)';

CREATE INDEX idx_courses_type ON courses(course_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_courses_active ON courses(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_courses_deleted ON courses(deleted_at);

-- Prerrequisitos (auto-referencial muchos a muchos)
CREATE TABLE course_prerequisites (
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    prerequisite_course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (course_id, prerequisite_course_id),
    CONSTRAINT chk_not_self_prerequisite CHECK (course_id != prerequisite_course_id)
);

COMMENT ON TABLE course_prerequisites IS 'Prerrequisitos entre cursos (N:M)';

CREATE INDEX idx_prerequisites_course ON course_prerequisites(course_id);
CREATE INDEX idx_prerequisites_prereq ON course_prerequisites(prerequisite_course_id);

-- =============================================================================
-- PERÍODOS ACADÉMICOS
-- =============================================================================

CREATE TABLE academic_periods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT false,
    
    -- Auditoría
    created_at TIMESTAMP DEFAULT NOW(),
    created_by UUID REFERENCES system_users(id) ON DELETE SET NULL,
    updated_at TIMESTAMP DEFAULT NOW(),
    updated_by UUID REFERENCES system_users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP,
    deleted_by UUID REFERENCES system_users(id) ON DELETE SET NULL,
    
    CONSTRAINT chk_period_dates CHECK (end_date > start_date)
);

COMMENT ON TABLE academic_periods IS 'Períodos académicos (semestres)';
COMMENT ON COLUMN academic_periods.name IS 'Nombre del período (ej: 2024-1, 2024-2)';
COMMENT ON COLUMN academic_periods.is_active IS 'Solo un período puede estar activo a la vez';

CREATE INDEX idx_periods_active ON academic_periods(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_periods_dates ON academic_periods(start_date, end_date);
CREATE INDEX idx_periods_deleted ON academic_periods(deleted_at);

-- =============================================================================
-- CURSOS PROGRAMADOS (oferta de cursos por período)
-- =============================================================================

CREATE TABLE scheduled_courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE RESTRICT,
    academic_period_id UUID NOT NULL REFERENCES academic_periods(id) ON DELETE CASCADE,
    max_students INTEGER,
    schedule VARCHAR(255),
    classroom VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    
    -- Auditoría
    created_at TIMESTAMP DEFAULT NOW(),
    created_by UUID REFERENCES system_users(id) ON DELETE SET NULL,
    updated_at TIMESTAMP DEFAULT NOW(),
    updated_by UUID REFERENCES system_users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP,
    deleted_by UUID REFERENCES system_users(id) ON DELETE SET NULL,
    
    CONSTRAINT uk_course_period UNIQUE(course_id, academic_period_id)
);

COMMENT ON TABLE scheduled_courses IS 'Cursos programados en períodos específicos';
COMMENT ON COLUMN scheduled_courses.schedule IS 'Horario (ej: Lunes y Miércoles 18:00-20:00)';

CREATE INDEX idx_scheduled_courses_course ON scheduled_courses(course_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_scheduled_courses_period ON scheduled_courses(academic_period_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_scheduled_courses_active ON scheduled_courses(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_scheduled_courses_deleted ON scheduled_courses(deleted_at);

COMMIT;
