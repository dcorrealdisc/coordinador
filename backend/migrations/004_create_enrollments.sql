-- Migration: 004_create_enrollments
-- Description: Inscripciones de estudiantes y calificaciones
-- Author: Agente DBA
-- Date: 2024-02-13

BEGIN;

-- =============================================================================
-- INSCRIPCIONES Y CALIFICACIONES
-- =============================================================================

CREATE TABLE enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    scheduled_course_id UUID NOT NULL REFERENCES scheduled_courses(id) ON DELETE RESTRICT,
    enrolled_at TIMESTAMP DEFAULT NOW(),
    status VARCHAR(20) NOT NULL CHECK (status IN ('enrolled', 'completed', 'withdrawn', 'failed')) DEFAULT 'enrolled',
    final_grade NUMERIC(3,2),
    credits_earned INTEGER DEFAULT 0,
    graded_by UUID REFERENCES tutors(id) ON DELETE SET NULL,
    graded_at TIMESTAMP,
    notes TEXT,
    
    -- Auditoría completa
    created_at TIMESTAMP DEFAULT NOW(),
    created_by UUID REFERENCES system_users(id) ON DELETE SET NULL,
    updated_at TIMESTAMP DEFAULT NOW(),
    updated_by UUID REFERENCES system_users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP,
    deleted_by UUID REFERENCES system_users(id) ON DELETE SET NULL,
    
    CONSTRAINT uk_student_scheduled_course UNIQUE(student_id, scheduled_course_id),
    CONSTRAINT chk_grade_range CHECK (final_grade IS NULL OR (final_grade >= 0 AND final_grade <= 5)),
    CONSTRAINT chk_credits_earned CHECK (credits_earned >= 0)
);

COMMENT ON TABLE enrollments IS 'Inscripciones de estudiantes a cursos programados';
COMMENT ON COLUMN enrollments.status IS 'Estado: enrolled, completed, withdrawn, failed';
COMMENT ON COLUMN enrollments.final_grade IS 'Nota final (escala 0.00 a 5.00)';
COMMENT ON COLUMN enrollments.graded_by IS 'Tutor que registró la calificación (si aplica)';
COMMENT ON COLUMN enrollments.notes IS 'Notas adicionales sobre el desempeño o situaciones especiales';

CREATE INDEX idx_enrollments_student ON enrollments(student_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_enrollments_scheduled_course ON enrollments(scheduled_course_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_enrollments_status ON enrollments(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_enrollments_graded_by ON enrollments(graded_by) WHERE deleted_at IS NULL;
CREATE INDEX idx_enrollments_deleted ON enrollments(deleted_at);

-- Índice compuesto para queries comunes
CREATE INDEX idx_enrollments_student_status 
ON enrollments(student_id, status) 
WHERE deleted_at IS NULL;

COMMIT;
