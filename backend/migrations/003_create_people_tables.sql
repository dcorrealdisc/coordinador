-- Migration: 003_create_people_tables
-- Description: Estudiantes, profesores y tutores
-- Author: Agente DBA
-- Date: 2024-02-13

BEGIN;

-- =============================================================================
-- ESTUDIANTES
-- =============================================================================

CREATE TABLE students (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Información personal
    full_name VARCHAR(255) NOT NULL,
    document_id VARCHAR(50) UNIQUE,
    birth_date DATE NOT NULL,
    profile_photo_url TEXT,
    
    -- Procedencia
    city_origin_id UUID REFERENCES cities(id) ON DELETE SET NULL,
    country_origin_id UUID NOT NULL REFERENCES countries(id) ON DELETE RESTRICT,
    
    -- Contacto (arrays para múltiples valores)
    emails TEXT[] NOT NULL,
    phones TEXT[],
    
    -- Laboral
    company_id UUID REFERENCES companies(id) ON DELETE SET NULL,
    
    -- Estado académico
    status VARCHAR(20) NOT NULL CHECK (status IN ('active', 'graduated', 'withdrawn', 'suspended')),
    cohort VARCHAR(10) NOT NULL,
    enrollment_date DATE NOT NULL,
    graduation_date DATE,
    
    -- Auditoría completa
    created_at TIMESTAMP DEFAULT NOW(),
    created_by UUID REFERENCES system_users(id) ON DELETE SET NULL,
    updated_at TIMESTAMP DEFAULT NOW(),
    updated_by UUID REFERENCES system_users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP,
    deleted_by UUID REFERENCES system_users(id) ON DELETE SET NULL,
    
    CONSTRAINT chk_emails_not_empty CHECK (array_length(emails, 1) > 0),
    CONSTRAINT chk_birth_date_valid CHECK (birth_date < CURRENT_DATE),
    CONSTRAINT chk_age_minimum CHECK (EXTRACT(YEAR FROM AGE(birth_date)) >= 18)
);

COMMENT ON TABLE students IS 'Estudiantes de la maestría (activos, graduados e históricos)';
COMMENT ON COLUMN students.status IS 'Estado: active, graduated, withdrawn, suspended';
COMMENT ON COLUMN students.cohort IS 'Cohorte de ingreso (ej: 2024-1, 2024-2)';
COMMENT ON COLUMN students.emails IS 'Array de correos electrónicos';
COMMENT ON COLUMN students.phones IS 'Array de teléfonos';

CREATE INDEX idx_students_status ON students(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_students_cohort ON students(cohort) WHERE deleted_at IS NULL;
CREATE INDEX idx_students_city ON students(city_origin_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_students_country ON students(country_origin_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_students_company ON students(company_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_students_birth_date ON students(birth_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_students_deleted ON students(deleted_at);

-- Universidades de procedencia (muchos a muchos)
CREATE TABLE student_universities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    university_id UUID NOT NULL REFERENCES universities(id) ON DELETE RESTRICT,
    degree_obtained VARCHAR(255),
    graduation_year INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uk_student_university UNIQUE(student_id, university_id),
    CONSTRAINT chk_graduation_year CHECK (graduation_year IS NULL OR (graduation_year >= 1950 AND graduation_year <= EXTRACT(YEAR FROM CURRENT_DATE)))
);

COMMENT ON TABLE student_universities IS 'Universidades de procedencia de estudiantes (N:M)';
COMMENT ON COLUMN student_universities.degree_obtained IS 'Título obtenido (ej: Ingeniería de Sistemas)';

CREATE INDEX idx_student_universities_student ON student_universities(student_id);
CREATE INDEX idx_student_universities_university ON student_universities(university_id);

-- =============================================================================
-- PROFESORES
-- =============================================================================

CREATE TABLE professors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(50),
    birth_date DATE NOT NULL,
    profile_photo_url TEXT,
    specialization VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    
    -- Auditoría completa
    created_at TIMESTAMP DEFAULT NOW(),
    created_by UUID REFERENCES system_users(id) ON DELETE SET NULL,
    updated_at TIMESTAMP DEFAULT NOW(),
    updated_by UUID REFERENCES system_users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP,
    deleted_by UUID REFERENCES system_users(id) ON DELETE SET NULL,
    
    CONSTRAINT chk_professor_birth_date CHECK (birth_date < CURRENT_DATE)
);

COMMENT ON TABLE professors IS 'Profesores que dictan cursos';
COMMENT ON COLUMN professors.specialization IS 'Área de especialización del profesor';

CREATE INDEX idx_professors_active ON professors(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_professors_birth_date ON professors(birth_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_professors_deleted ON professors(deleted_at);

-- Asignación de profesores a cursos programados
CREATE TABLE course_professor_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scheduled_course_id UUID NOT NULL REFERENCES scheduled_courses(id) ON DELETE CASCADE,
    professor_id UUID NOT NULL REFERENCES professors(id) ON DELETE RESTRICT,
    assigned_at TIMESTAMP DEFAULT NOW(),
    assigned_by UUID REFERENCES system_users(id) ON DELETE SET NULL,
    CONSTRAINT uk_scheduled_course_professor UNIQUE(scheduled_course_id, professor_id)
);

COMMENT ON TABLE course_professor_assignments IS 'Asignación de profesores a cursos programados';

CREATE INDEX idx_professor_assignments_course ON course_professor_assignments(scheduled_course_id);
CREATE INDEX idx_professor_assignments_professor ON course_professor_assignments(professor_id);

-- =============================================================================
-- TUTORES/MONITORES
-- =============================================================================

CREATE TABLE tutors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name VARCHAR(255) NOT NULL,
    emails TEXT[] NOT NULL,
    phones TEXT[],
    birth_date DATE NOT NULL,
    profile_photo_url TEXT,
    current_employer VARCHAR(255),
    academic_background TEXT,
    status VARCHAR(20) NOT NULL CHECK (status IN ('active', 'inactive')) DEFAULT 'active',
    
    -- Auditoría completa
    created_at TIMESTAMP DEFAULT NOW(),
    created_by UUID REFERENCES system_users(id) ON DELETE SET NULL,
    updated_at TIMESTAMP DEFAULT NOW(),
    updated_by UUID REFERENCES system_users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP,
    deleted_by UUID REFERENCES system_users(id) ON DELETE SET NULL,
    
    CONSTRAINT chk_tutor_emails_not_empty CHECK (array_length(emails, 1) > 0),
    CONSTRAINT chk_tutor_birth_date CHECK (birth_date < CURRENT_DATE)
);

COMMENT ON TABLE tutors IS 'Pool de tutores/monitores disponibles';
COMMENT ON COLUMN tutors.academic_background IS 'Formación académica del tutor';

CREATE INDEX idx_tutors_status ON tutors(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_tutors_birth_date ON tutors(birth_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_tutors_deleted ON tutors(deleted_at);

-- Intereses de tutores (tutores expresan interés en cursos)
CREATE TABLE tutor_course_interests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tutor_id UUID NOT NULL REFERENCES tutors(id) ON DELETE CASCADE,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    interested_at TIMESTAMP DEFAULT NOW(),
    notes TEXT,
    CONSTRAINT uk_tutor_course_interest UNIQUE(tutor_id, course_id)
);

COMMENT ON TABLE tutor_course_interests IS 'Tutores expresan interés en cursos específicos';
COMMENT ON COLUMN tutor_course_interests.notes IS 'Razón del interés o experiencia relacionada';

CREATE INDEX idx_tutor_interests_tutor ON tutor_course_interests(tutor_id);
CREATE INDEX idx_tutor_interests_course ON tutor_course_interests(course_id);

-- Asignaciones reales de tutores a cursos programados
CREATE TABLE course_tutor_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scheduled_course_id UUID NOT NULL REFERENCES scheduled_courses(id) ON DELETE CASCADE,
    tutor_id UUID NOT NULL REFERENCES tutors(id) ON DELETE RESTRICT,
    assigned_at TIMESTAMP DEFAULT NOW(),
    assigned_by UUID REFERENCES professors(id) ON DELETE SET NULL,
    CONSTRAINT uk_scheduled_course_tutor UNIQUE(scheduled_course_id, tutor_id)
);

COMMENT ON TABLE course_tutor_assignments IS 'Asignaciones reales de tutores a cursos';
COMMENT ON COLUMN course_tutor_assignments.assigned_by IS 'Profesor que asignó al tutor';

CREATE INDEX idx_tutor_assignments_course ON course_tutor_assignments(scheduled_course_id);
CREATE INDEX idx_tutor_assignments_tutor ON course_tutor_assignments(tutor_id);
CREATE INDEX idx_tutor_assignments_assigner ON course_tutor_assignments(assigned_by);

COMMIT;
