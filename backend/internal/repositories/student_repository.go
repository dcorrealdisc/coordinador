package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dcorreal/coordinador/internal/models"
)

// StudentFilters holds the query filters for listing students.
type StudentFilters struct {
	Status          *string
	Cohort          *string
	CountryOriginID *uuid.UUID
	Search          *string // ILIKE search on full_name
	Limit           int
	Offset          int
}

// StudentRepository defines the data access interface for students.
type StudentRepository interface {
	Create(ctx context.Context, student *models.Student) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Student, error)
	List(ctx context.Context, filters StudentFilters) ([]*models.Student, error)
	Update(ctx context.Context, student *models.Student) error
	Delete(ctx context.Context, id uuid.UUID, deletedBy *uuid.UUID) error
	Count(ctx context.Context, filters StudentFilters) (int, error)
}

type studentRepository struct {
	db *pgxpool.Pool
}

// NewStudentRepository creates a new StudentRepository backed by pgxpool.
func NewStudentRepository(db *pgxpool.Pool) StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) Create(ctx context.Context, student *models.Student) error {
	query := `
		INSERT INTO students (
			id, full_name, document_id, birth_date, profile_photo_url,
			city_origin_id, country_origin_id, emails, phones,
			company_id, student_code, status, cohort, enrollment_date
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		student.ID,
		student.FullName,
		student.DocumentID,
		student.BirthDate,
		student.ProfilePhotoURL,
		student.CityOriginID,
		student.CountryOriginID,
		student.Emails,
		student.Phones,
		student.CompanyID,
		student.StudentCode,
		student.Status,
		student.Cohort,
		student.EnrollmentDate,
	).Scan(&student.CreatedAt, &student.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create student: %w", err)
	}

	return nil
}

func (r *studentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Student, error) {
	query := `
		SELECT
			id, full_name, document_id, birth_date, profile_photo_url,
			city_origin_id, country_origin_id, emails, phones,
			company_id, student_code, status, cohort, enrollment_date, graduation_date,
			created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM students
		WHERE id = $1 AND deleted_at IS NULL
	`

	student := &models.Student{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&student.ID,
		&student.FullName,
		&student.DocumentID,
		&student.BirthDate,
		&student.ProfilePhotoURL,
		&student.CityOriginID,
		&student.CountryOriginID,
		&student.Emails,
		&student.Phones,
		&student.CompanyID,
		&student.StudentCode,
		&student.Status,
		&student.Cohort,
		&student.EnrollmentDate,
		&student.GraduationDate,
		&student.CreatedAt,
		&student.CreatedBy,
		&student.UpdatedAt,
		&student.UpdatedBy,
		&student.DeletedAt,
		&student.DeletedBy,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("student not found")
		}
		return nil, fmt.Errorf("failed to get student: %w", err)
	}

	return student, nil
}

func (r *studentRepository) List(ctx context.Context, filters StudentFilters) ([]*models.Student, error) {
	query := `
		SELECT
			id, full_name, document_id, birth_date, profile_photo_url,
			city_origin_id, country_origin_id, emails, phones,
			company_id, student_code, status, cohort, enrollment_date, graduation_date,
			created_at, created_by, updated_at, updated_by
		FROM students
		WHERE deleted_at IS NULL
	`

	args := []interface{}{}
	argCount := 1

	if filters.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
		argCount++
	}

	if filters.Cohort != nil {
		query += fmt.Sprintf(" AND cohort = $%d", argCount)
		args = append(args, *filters.Cohort)
		argCount++
	}

	if filters.CountryOriginID != nil {
		query += fmt.Sprintf(" AND country_origin_id = $%d", argCount)
		args = append(args, *filters.CountryOriginID)
		argCount++
	}

	if filters.Search != nil {
		query += fmt.Sprintf(" AND full_name ILIKE $%d", argCount)
		args = append(args, "%"+*filters.Search+"%")
		argCount++
	}

	query += " ORDER BY enrollment_date DESC"

	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.Limit)
		argCount++
	}

	if filters.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, filters.Offset)
		argCount++
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list students: %w", err)
	}
	defer rows.Close()

	students := []*models.Student{}
	for rows.Next() {
		student := &models.Student{}
		err := rows.Scan(
			&student.ID,
			&student.FullName,
			&student.DocumentID,
			&student.BirthDate,
			&student.ProfilePhotoURL,
			&student.CityOriginID,
			&student.CountryOriginID,
			&student.Emails,
			&student.Phones,
			&student.CompanyID,
			&student.StudentCode,
			&student.Status,
			&student.Cohort,
			&student.EnrollmentDate,
			&student.GraduationDate,
			&student.CreatedAt,
			&student.CreatedBy,
			&student.UpdatedAt,
			&student.UpdatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan student row: %w", err)
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating student rows: %w", err)
	}

	return students, nil
}

func (r *studentRepository) Update(ctx context.Context, student *models.Student) error {
	query := `
		UPDATE students
		SET
			full_name = $2,
			document_id = $3,
			profile_photo_url = $4,
			emails = $5,
			phones = $6,
			company_id = $7,
			student_code = $8,
			status = $9,
			updated_by = $10
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		student.ID,
		student.FullName,
		student.DocumentID,
		student.ProfilePhotoURL,
		student.Emails,
		student.Phones,
		student.CompanyID,
		student.StudentCode,
		student.Status,
		student.UpdatedBy,
	).Scan(&student.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("student not found")
		}
		return fmt.Errorf("failed to update student: %w", err)
	}

	return nil
}

func (r *studentRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy *uuid.UUID) error {
	query := `
		UPDATE students
		SET deleted_at = NOW(), deleted_by = $2
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(ctx, query, id, deletedBy)
	if err != nil {
		return fmt.Errorf("failed to delete student: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("student not found")
	}

	return nil
}

func (r *studentRepository) Count(ctx context.Context, filters StudentFilters) (int, error) {
	query := "SELECT COUNT(*) FROM students WHERE deleted_at IS NULL"

	args := []interface{}{}
	argCount := 1

	if filters.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
		argCount++
	}

	if filters.Cohort != nil {
		query += fmt.Sprintf(" AND cohort = $%d", argCount)
		args = append(args, *filters.Cohort)
		argCount++
	}

	if filters.CountryOriginID != nil {
		query += fmt.Sprintf(" AND country_origin_id = $%d", argCount)
		args = append(args, *filters.CountryOriginID)
		argCount++
	}

	if filters.Search != nil {
		query += fmt.Sprintf(" AND full_name ILIKE $%d", argCount)
		args = append(args, "%"+*filters.Search+"%")
		argCount++
	}

	var count int
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count students: %w", err)
	}

	return count, nil
}
