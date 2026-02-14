-- Migration 011: Hacer birth_date y emails opcionales

ALTER TABLE students ALTER COLUMN birth_date DROP NOT NULL;
ALTER TABLE students ALTER COLUMN emails DROP NOT NULL;
