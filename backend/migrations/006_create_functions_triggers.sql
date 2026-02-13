-- Migration: 006_create_functions_triggers
-- Description: Funciones helper y triggers automáticos
-- Author: Agente DBA
-- Date: 2024-02-13

BEGIN;

-- =============================================================================
-- FUNCIÓN: Actualizar updated_at automáticamente
-- =============================================================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION update_updated_at_column() IS 'Actualiza automáticamente el campo updated_at al modificar un registro';

-- =============================================================================
-- TRIGGERS: Aplicar update_updated_at a todas las tablas relevantes
-- =============================================================================

-- Catálogos
CREATE TRIGGER trigger_countries_updated_at
    BEFORE UPDATE ON countries
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_cities_updated_at
    BEFORE UPDATE ON cities
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_universities_updated_at
    BEFORE UPDATE ON universities
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_companies_updated_at
    BEFORE UPDATE ON companies
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_system_users_updated_at
    BEFORE UPDATE ON system_users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_program_configuration_updated_at
    BEFORE UPDATE ON program_configuration
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Académico
CREATE TRIGGER trigger_courses_updated_at
    BEFORE UPDATE ON courses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_academic_periods_updated_at
    BEFORE UPDATE ON academic_periods
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_scheduled_courses_updated_at
    BEFORE UPDATE ON scheduled_courses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Personas
CREATE TRIGGER trigger_students_updated_at
    BEFORE UPDATE ON students
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_professors_updated_at
    BEFORE UPDATE ON professors
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_tutors_updated_at
    BEFORE UPDATE ON tutors
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Inscripciones
CREATE TRIGGER trigger_enrollments_updated_at
    BEFORE UPDATE ON enrollments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =============================================================================
-- FUNCIÓN: Refresh de todas las vistas materializadas
-- =============================================================================

CREATE OR REPLACE FUNCTION refresh_all_materialized_views()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY student_academic_progress;
    REFRESH MATERIALIZED VIEW CONCURRENTLY course_period_statistics;
    REFRESH MATERIALIZED VIEW CONCURRENTLY students_by_location;
    REFRESH MATERIALIZED VIEW CONCURRENTLY students_by_university;
    REFRESH MATERIALIZED VIEW CONCURRENTLY students_by_company;
    REFRESH MATERIALIZED VIEW CONCURRENTLY tutor_workload;
    REFRESH MATERIALIZED VIEW CONCURRENTLY students_age_distribution;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION refresh_all_materialized_views() IS 'Refresca todas las vistas materializadas para reportes';

-- =============================================================================
-- FUNCIÓN: Validar que solo un período esté activo
-- =============================================================================

CREATE OR REPLACE FUNCTION validate_single_active_period()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_active = true THEN
        -- Desactivar todos los otros períodos
        UPDATE academic_periods 
        SET is_active = false 
        WHERE id != NEW.id AND is_active = true AND deleted_at IS NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION validate_single_active_period() IS 'Asegura que solo un período académico esté activo a la vez';

CREATE TRIGGER trigger_single_active_period
    BEFORE INSERT OR UPDATE ON academic_periods
    FOR EACH ROW
    WHEN (NEW.is_active = true)
    EXECUTE FUNCTION validate_single_active_period();

-- =============================================================================
-- FUNCIÓN: Calcular créditos ganados al aprobar un curso
-- =============================================================================

CREATE OR REPLACE FUNCTION calculate_credits_earned()
RETURNS TRIGGER AS $$
DECLARE
    course_credits INTEGER;
    passing_grade NUMERIC;
BEGIN
    -- Obtener nota aprobatoria de configuración
    SELECT (value::numeric) INTO passing_grade 
    FROM program_configuration 
    WHERE key = 'passing_grade';
    
    -- Si el curso fue completado y aprobado
    IF NEW.status = 'completed' AND NEW.final_grade >= passing_grade THEN
        -- Obtener créditos del curso
        SELECT c.credits INTO course_credits
        FROM scheduled_courses sc
        JOIN courses c ON sc.course_id = c.id
        WHERE sc.id = NEW.scheduled_course_id;
        
        NEW.credits_earned = course_credits;
    ELSE
        NEW.credits_earned = 0;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION calculate_credits_earned() IS 'Calcula automáticamente los créditos ganados al completar un curso';

CREATE TRIGGER trigger_calculate_credits
    BEFORE INSERT OR UPDATE ON enrollments
    FOR EACH ROW
    WHEN (NEW.status = 'completed' OR NEW.final_grade IS NOT NULL)
    EXECUTE FUNCTION calculate_credits_earned();

-- =============================================================================
-- FUNCIÓN: Validar límite de tutores por curso
-- =============================================================================

CREATE OR REPLACE FUNCTION validate_tutor_limit()
RETURNS TRIGGER AS $$
DECLARE
    tutor_count INTEGER;
    max_tutors INTEGER;
BEGIN
    -- Obtener límite de configuración
    SELECT (value::integer) INTO max_tutors 
    FROM program_configuration 
    WHERE key = 'max_courses_per_tutor';
    
    -- Contar asignaciones actuales del tutor en el período del curso
    SELECT COUNT(*) INTO tutor_count
    FROM course_tutor_assignments cta
    JOIN scheduled_courses sc ON cta.scheduled_course_id = sc.id
    WHERE cta.tutor_id = NEW.tutor_id
      AND sc.academic_period_id = (
          SELECT academic_period_id 
          FROM scheduled_courses 
          WHERE id = NEW.scheduled_course_id
      );
    
    -- Validar límite
    IF tutor_count >= max_tutors THEN
        RAISE EXCEPTION 'El tutor ya tiene el máximo de cursos asignados (%) para este período', max_tutors;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION validate_tutor_limit() IS 'Valida que un tutor no exceda el límite de cursos por período';

CREATE TRIGGER trigger_validate_tutor_limit
    BEFORE INSERT ON course_tutor_assignments
    FOR EACH ROW
    EXECUTE FUNCTION validate_tutor_limit();

-- =============================================================================
-- FUNCIÓN: Obtener cursos pendientes de un estudiante
-- =============================================================================

CREATE OR REPLACE FUNCTION get_pending_courses(student_uuid UUID)
RETURNS TABLE (
    course_id UUID,
    course_code VARCHAR,
    course_name VARCHAR,
    credits INTEGER,
    course_type VARCHAR,
    has_prerequisites BOOLEAN,
    prerequisites_met BOOLEAN
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        c.id,
        c.code,
        c.name,
        c.credits,
        c.course_type,
        EXISTS(SELECT 1 FROM course_prerequisites WHERE course_id = c.id) as has_prerequisites,
        NOT EXISTS(
            SELECT 1 
            FROM course_prerequisites cp
            WHERE cp.course_id = c.id
              AND cp.prerequisite_course_id NOT IN (
                  SELECT sc.course_id
                  FROM enrollments e
                  JOIN scheduled_courses sc ON e.scheduled_course_id = sc.id
                  WHERE e.student_id = student_uuid
                    AND e.status = 'completed'
                    AND e.final_grade >= 3.0
                    AND e.deleted_at IS NULL
              )
        ) as prerequisites_met
    FROM courses c
    WHERE c.id NOT IN (
        SELECT sc.course_id
        FROM enrollments e
        JOIN scheduled_courses sc ON e.scheduled_course_id = sc.id
        WHERE e.student_id = student_uuid
          AND e.status = 'completed'
          AND e.final_grade >= 3.0
          AND e.deleted_at IS NULL
    )
    AND c.is_active = true
    AND c.deleted_at IS NULL
    ORDER BY c.course_type, c.code;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_pending_courses(UUID) IS 'Obtiene los cursos pendientes de un estudiante y si cumple prerrequisitos';

-- =============================================================================
-- FUNCIÓN: Calcular GPA de un estudiante
-- =============================================================================

CREATE OR REPLACE FUNCTION calculate_student_gpa(student_uuid UUID)
RETURNS NUMERIC AS $$
DECLARE
    student_gpa NUMERIC;
BEGIN
    SELECT ROUND(AVG(final_grade), 2) INTO student_gpa
    FROM enrollments
    WHERE student_id = student_uuid
      AND status = 'completed'
      AND final_grade IS NOT NULL
      AND deleted_at IS NULL;
    
    RETURN COALESCE(student_gpa, 0);
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION calculate_student_gpa(UUID) IS 'Calcula el GPA (promedio) de un estudiante';

COMMIT;
