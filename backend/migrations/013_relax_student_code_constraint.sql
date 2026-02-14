-- Migration 013: Relajar constraint de student_code a solo 9 dígitos
ALTER TABLE students DROP CONSTRAINT IF EXISTS chk_student_code_format;
ALTER TABLE students ADD CONSTRAINT chk_student_code_format CHECK (student_code ~ '^[0-9]{9}$');
