-- Migration: 001_create_base_schema
-- Description: Catálogos maestros y tablas de configuración base
-- Author: Agente DBA
-- Date: 2024-02-13

BEGIN;

-- =============================================================================
-- CATÁLOGOS MAESTROS (para consolidación de datos)
-- =============================================================================

-- Países
CREATE TABLE countries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(3) NOT NULL UNIQUE,  -- ISO 3166-1 alpha-3
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

COMMENT ON TABLE countries IS 'Catálogo de países para consolidar información geográfica';
COMMENT ON COLUMN countries.code IS 'Código ISO 3166-1 alpha-3 (ej: COL, USA, MEX)';

-- Ciudades
CREATE TABLE cities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    country_id UUID NOT NULL REFERENCES countries(id) ON DELETE RESTRICT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uk_city_country UNIQUE(name, country_id)
);

COMMENT ON TABLE cities IS 'Catálogo de ciudades vinculadas a países';

CREATE INDEX idx_cities_country ON cities(country_id);

-- Universidades (consolidadas)
CREATE TABLE universities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    country_id UUID NOT NULL REFERENCES countries(id) ON DELETE RESTRICT,
    city_id UUID REFERENCES cities(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uk_universities_name_country UNIQUE(name, country_id)
);

COMMENT ON TABLE universities IS 'Catálogo consolidado de universidades de procedencia';

CREATE INDEX idx_universities_country ON universities(country_id);
CREATE INDEX idx_universities_city ON universities(city_id);

-- Empresas/Empleadores (consolidadas)
CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

COMMENT ON TABLE companies IS 'Catálogo consolidado de empresas donde trabajan estudiantes';

-- =============================================================================
-- USUARIOS DEL SISTEMA (para auditoría)
-- =============================================================================

CREATE TABLE system_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    full_name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'coordinator', 'staff')),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

COMMENT ON TABLE system_users IS 'Usuarios administrativos del sistema para auditoría';

CREATE INDEX idx_system_users_active ON system_users(is_active);
CREATE INDEX idx_system_users_role ON system_users(role);

-- =============================================================================
-- CONFIGURACIÓN DEL PROGRAMA
-- =============================================================================

CREATE TABLE program_configuration (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key VARCHAR(100) NOT NULL UNIQUE,
    value JSONB NOT NULL,
    description TEXT,
    updated_at TIMESTAMP DEFAULT NOW()
);

COMMENT ON TABLE program_configuration IS 'Configuración general del programa de maestría';

-- Insertar configuraciones iniciales
INSERT INTO program_configuration (key, value, description) VALUES
    ('total_credits_required', '48', 'Total de créditos para graduarse'),
    ('required_credits', '36', 'Créditos obligatorios requeridos'),
    ('elective_credits', '12', 'Créditos electivos requeridos'),
    ('max_courses_per_tutor', '2', 'Máximo de cursos que puede tener un tutor por período'),
    ('passing_grade', '3.0', 'Nota mínima aprobatoria (escala 0-5)'),
    ('min_student_age', '18', 'Edad mínima para estudiantes');

-- =============================================================================
-- DATOS INICIALES - Países de Latinoamérica
-- =============================================================================

INSERT INTO countries (code, name) VALUES 
    ('COL', 'Colombia'),
    ('MEX', 'México'),
    ('ARG', 'Argentina'),
    ('BRA', 'Brasil'),
    ('CHL', 'Chile'),
    ('PER', 'Perú'),
    ('VEN', 'Venezuela'),
    ('ECU', 'Ecuador'),
    ('BOL', 'Bolivia'),
    ('PRY', 'Paraguay'),
    ('URY', 'Uruguay'),
    ('USA', 'Estados Unidos'),
    ('CAN', 'Canadá'),
    ('ESP', 'España');

-- Ciudades principales de Colombia (ejemplo)
INSERT INTO cities (name, country_id) 
SELECT 'Bogotá', id FROM countries WHERE code = 'COL'
UNION ALL SELECT 'Medellín', id FROM countries WHERE code = 'COL'
UNION ALL SELECT 'Cali', id FROM countries WHERE code = 'COL'
UNION ALL SELECT 'Barranquilla', id FROM countries WHERE code = 'COL'
UNION ALL SELECT 'Cartagena', id FROM countries WHERE code = 'COL'
UNION ALL SELECT 'Bucaramanga', id FROM countries WHERE code = 'COL'
UNION ALL SELECT 'Manizales', id FROM countries WHERE code = 'COL'
UNION ALL SELECT 'Pereira', id FROM countries WHERE code = 'COL';

COMMIT;
