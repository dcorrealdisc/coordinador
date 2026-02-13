-- Migration: 005_create_materialized_views
-- Description: Vistas materializadas para reportes (CQRS Read Path)
-- Author: Agente DBA
-- Date: 2024-02-13

BEGIN;

-- =============================================================================
-- VISTA: Progreso Académico de Estudiantes
-- =============================================================================

CREATE MATERIALIZED VIEW student_academic_progress AS
SELECT 
    s.id as student_id,
    s.full_name,
    s.status,
    s.cohort,
    s.enrollment_date,
    co.name as country_origin,
    ci.name as city_origin,
    comp.name as current_employer,
    
    -- Estadísticas de cursos
    COUNT(DISTINCT e.id) FILTER (WHERE e.status = 'completed' AND e.final_grade >= 3.0) as courses_completed,
    COUNT(DISTINCT e.id) FILTER (WHERE e.status = 'completed' AND e.final_grade < 3.0) as courses_failed,
    COUNT(DISTINCT e.id) FILTER (WHERE e.status = 'enrolled') as courses_in_progress,
    
    -- Créditos
    COALESCE(SUM(c.credits) FILTER (WHERE e.status = 'completed' AND e.final_grade >= 3.0), 0) as total_credits_earned,
    COALESCE(SUM(c.credits) FILTER (WHERE e.status = 'completed' AND e.final_grade >= 3.0 AND c.course_type = 'required'), 0) as required_credits_earned,
    COALESCE(SUM(c.credits) FILTER (WHERE e.status = 'completed' AND e.final_grade >= 3.0 AND c.course_type = 'elective'), 0) as elective_credits_earned,
    48 - COALESCE(SUM(c.credits) FILTER (WHERE e.status = 'completed' AND e.final_grade >= 3.0), 0) as credits_remaining,
    
    -- Porcentaje de avance
    ROUND((COALESCE(SUM(c.credits) FILTER (WHERE e.status = 'completed' AND e.final_grade >= 3.0), 0)::numeric / 48) * 100, 2) as completion_percentage,
    
    -- Promedio (GPA)
    COALESCE(ROUND(AVG(e.final_grade) FILTER (WHERE e.status = 'completed'), 2), 0) as gpa,
    
    -- Tiempo en el programa
    EXTRACT(YEAR FROM AGE(COALESCE(s.graduation_date, CURRENT_DATE), s.enrollment_date)) * 12 + 
    EXTRACT(MONTH FROM AGE(COALESCE(s.graduation_date, CURRENT_DATE), s.enrollment_date)) as months_in_program
    
FROM students s
LEFT JOIN countries co ON s.country_origin_id = co.id
LEFT JOIN cities ci ON s.city_origin_id = ci.id
LEFT JOIN companies comp ON s.company_id = comp.id
LEFT JOIN enrollments e ON s.id = e.student_id AND e.deleted_at IS NULL
LEFT JOIN scheduled_courses sc ON e.scheduled_course_id = sc.id AND sc.deleted_at IS NULL
LEFT JOIN courses c ON sc.course_id = c.id AND c.deleted_at IS NULL
WHERE s.deleted_at IS NULL
GROUP BY s.id, s.full_name, s.status, s.cohort, s.enrollment_date, s.graduation_date, co.name, ci.name, comp.name;

CREATE UNIQUE INDEX idx_student_progress_id ON student_academic_progress(student_id);
CREATE INDEX idx_student_progress_status ON student_academic_progress(status);
CREATE INDEX idx_student_progress_cohort ON student_academic_progress(cohort);
CREATE INDEX idx_student_progress_completion ON student_academic_progress(completion_percentage);

COMMENT ON MATERIALIZED VIEW student_academic_progress IS 'Progreso académico consolidado por estudiante';

-- =============================================================================
-- VISTA: Estadísticas de Cursos por Período
-- =============================================================================

CREATE MATERIALIZED VIEW course_period_statistics AS
SELECT 
    sc.id as scheduled_course_id,
    c.id as course_id,
    c.code,
    c.name as course_name,
    c.course_type,
    c.credits,
    ap.id as period_id,
    ap.name as period,
    
    -- Inscripciones
    COUNT(DISTINCT e.student_id) as enrolled_students,
    sc.max_students,
    CASE 
        WHEN sc.max_students IS NOT NULL AND sc.max_students > 0 
        THEN ROUND((COUNT(DISTINCT e.student_id)::numeric / sc.max_students) * 100, 2)
        ELSE NULL
    END as capacity_percentage,
    
    -- Resultados
    AVG(e.final_grade) FILTER (WHERE e.status = 'completed') as average_grade,
    COUNT(*) FILTER (WHERE e.status = 'completed' AND e.final_grade >= 3.0) as students_passed,
    COUNT(*) FILTER (WHERE e.status = 'completed' AND e.final_grade < 3.0) as students_failed,
    COUNT(*) FILTER (WHERE e.status = 'withdrawn') as students_withdrawn,
    
    -- Tasas de aprobación
    CASE 
        WHEN COUNT(*) FILTER (WHERE e.status = 'completed') > 0
        THEN ROUND((COUNT(*) FILTER (WHERE e.status = 'completed' AND e.final_grade >= 3.0)::numeric / 
                    COUNT(*) FILTER (WHERE e.status = 'completed')) * 100, 2)
        ELSE NULL
    END as pass_rate
    
FROM scheduled_courses sc
JOIN courses c ON sc.course_id = c.id AND c.deleted_at IS NULL
JOIN academic_periods ap ON sc.academic_period_id = ap.id AND ap.deleted_at IS NULL
LEFT JOIN enrollments e ON sc.id = e.scheduled_course_id AND e.deleted_at IS NULL
WHERE sc.deleted_at IS NULL
GROUP BY sc.id, c.id, c.code, c.name, c.course_type, c.credits, ap.id, ap.name, sc.max_students;

CREATE UNIQUE INDEX idx_course_stats_id ON course_period_statistics(scheduled_course_id);
CREATE INDEX idx_course_stats_course ON course_period_statistics(course_id);
CREATE INDEX idx_course_stats_period ON course_period_statistics(period_id);
CREATE INDEX idx_course_stats_type ON course_period_statistics(course_type);

COMMENT ON MATERIALIZED VIEW course_period_statistics IS 'Estadísticas de cursos por período académico';

-- =============================================================================
-- VISTA: Estudiantes por Ubicación (País y Ciudad)
-- =============================================================================

CREATE MATERIALIZED VIEW students_by_location AS
SELECT 
    co.id as country_id,
    co.name as country_name,
    ci.id as city_id,
    ci.name as city_name,
    
    COUNT(DISTINCT s.id) as total_students,
    COUNT(DISTINCT s.id) FILTER (WHERE s.status = 'active') as active_students,
    COUNT(DISTINCT s.id) FILTER (WHERE s.status = 'graduated') as graduated_students,
    COUNT(DISTINCT s.id) FILTER (WHERE s.status = 'withdrawn') as withdrawn_students,
    COUNT(DISTINCT s.id) FILTER (WHERE s.status = 'suspended') as suspended_students
    
FROM countries co
LEFT JOIN cities ci ON co.id = ci.country_id
LEFT JOIN students s ON s.country_origin_id = co.id 
    AND (s.city_origin_id = ci.id OR (s.city_origin_id IS NULL AND ci.id IS NULL))
    AND s.deleted_at IS NULL
GROUP BY co.id, co.name, ci.id, ci.name;

CREATE INDEX idx_students_location_country ON students_by_location(country_id);
CREATE INDEX idx_students_location_city ON students_by_location(city_id);

COMMENT ON MATERIALIZED VIEW students_by_location IS 'Distribución de estudiantes por país y ciudad';

-- =============================================================================
-- VISTA: Estudiantes por Universidad de Procedencia
-- =============================================================================

CREATE MATERIALIZED VIEW students_by_university AS
SELECT 
    u.id as university_id,
    u.name as university_name,
    co.name as country,
    
    COUNT(DISTINCT su.student_id) as student_count,
    COUNT(DISTINCT su.student_id) FILTER (WHERE s.status = 'active') as active_students,
    COUNT(DISTINCT su.student_id) FILTER (WHERE s.status = 'graduated') as graduated_students,
    
    -- Promedio GPA de estudiantes de esta universidad
    ROUND(AVG((
        SELECT AVG(e.final_grade) 
        FROM enrollments e 
        WHERE e.student_id = su.student_id 
          AND e.status = 'completed'
          AND e.deleted_at IS NULL
    )), 2) as average_gpa
    
FROM universities u
JOIN student_universities su ON u.id = su.university_id
JOIN students s ON su.student_id = s.id AND s.deleted_at IS NULL
JOIN countries co ON u.country_id = co.id
WHERE u.deleted_at IS NULL
GROUP BY u.id, u.name, co.name;

CREATE UNIQUE INDEX idx_students_by_uni_id ON students_by_university(university_id);

COMMENT ON MATERIALIZED VIEW students_by_university IS 'Estudiantes agrupados por universidad de procedencia';

-- =============================================================================
-- VISTA: Estudiantes por Empresa
-- =============================================================================

CREATE MATERIALIZED VIEW students_by_company AS
SELECT 
    c.id as company_id,
    c.name as company_name,
    
    COUNT(s.id) as student_count,
    COUNT(s.id) FILTER (WHERE s.status = 'active') as active_students,
    COUNT(s.id) FILTER (WHERE s.status = 'graduated') as graduated_students
    
FROM companies c
JOIN students s ON c.id = s.company_id AND s.deleted_at IS NULL
GROUP BY c.id, c.name;

CREATE UNIQUE INDEX idx_students_by_company_id ON students_by_company(company_id);

COMMENT ON MATERIALIZED VIEW students_by_company IS 'Estudiantes agrupados por empresa empleadora';

-- =============================================================================
-- VISTA: Carga de Trabajo de Tutores
-- =============================================================================

CREATE MATERIALIZED VIEW tutor_workload AS
SELECT 
    t.id as tutor_id,
    t.full_name as tutor_name,
    t.status,
    ap.name as period,
    
    COUNT(DISTINCT cta.scheduled_course_id) as courses_assigned,
    (SELECT value::integer FROM program_configuration WHERE key = 'max_courses_per_tutor') as max_allowed,
    
    -- Disponibilidad
    CASE 
        WHEN COUNT(DISTINCT cta.scheduled_course_id) >= (SELECT value::integer FROM program_configuration WHERE key = 'max_courses_per_tutor')
        THEN false
        ELSE true
    END as is_available
    
FROM tutors t
LEFT JOIN course_tutor_assignments cta ON t.id = cta.tutor_id
LEFT JOIN scheduled_courses sc ON cta.scheduled_course_id = sc.id AND sc.deleted_at IS NULL
LEFT JOIN academic_periods ap ON sc.academic_period_id = ap.id AND ap.deleted_at IS NULL
WHERE t.deleted_at IS NULL
GROUP BY t.id, t.full_name, t.status, ap.name;

CREATE INDEX idx_tutor_workload_tutor ON tutor_workload(tutor_id);
CREATE INDEX idx_tutor_workload_period ON tutor_workload(period);
CREATE INDEX idx_tutor_workload_available ON tutor_workload(is_available);

COMMENT ON MATERIALIZED VIEW tutor_workload IS 'Carga de trabajo de tutores por período';

-- =============================================================================
-- VISTA: Distribución de Edades
-- =============================================================================

CREATE MATERIALIZED VIEW students_age_distribution AS
SELECT 
    CASE 
        WHEN age < 25 THEN '18-24'
        WHEN age >= 25 AND age < 30 THEN '25-29'
        WHEN age >= 30 AND age < 35 THEN '30-34'
        WHEN age >= 35 AND age < 40 THEN '35-39'
        ELSE '40+'
    END as age_range,
    COUNT(*) as student_count,
    ROUND(AVG(age), 1) as average_age_in_range,
    COUNT(*) FILTER (WHERE status = 'active') as active_count,
    COUNT(*) FILTER (WHERE status = 'graduated') as graduated_count
FROM (
    SELECT 
        id,
        status,
        EXTRACT(YEAR FROM AGE(birth_date))::integer as age
    FROM students
    WHERE deleted_at IS NULL
) subquery
GROUP BY age_range
ORDER BY 
    CASE age_range
        WHEN '18-24' THEN 1
        WHEN '25-29' THEN 2
        WHEN '30-34' THEN 3
        WHEN '35-39' THEN 4
        ELSE 5
    END;

CREATE INDEX idx_age_distribution_range ON students_age_distribution(age_range);

COMMENT ON MATERIALIZED VIEW students_age_distribution IS 'Distribución de estudiantes por rango de edad';

COMMIT;
