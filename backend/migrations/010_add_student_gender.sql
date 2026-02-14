-- Migration 010: Agregar género del estudiante
-- Valores permitidos: 'M' (masculino) y 'F' (femenino), nullable

ALTER TABLE students
ADD COLUMN gender CHAR(1) CHECK (gender IN ('M', 'F'));

COMMENT ON COLUMN students.gender IS 'Género del estudiante: M = Masculino, F = Femenino';
