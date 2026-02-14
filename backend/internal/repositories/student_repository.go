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
	Status             *string
	Cohort             *string
	ResidenceCountryID *uuid.UUID
	Search             *string // ILIKE search on first_names || last_names
	Limit              int
	Offset             int
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
			id, first_names, last_names, document_id, birth_date, profile_photo_url,
			gender, nationality_country_id, residence_country_id, residence_city_id,
			emails, phones, company_id, job_title_category_id, profession_id,
			student_code, status, cohort, enrollment_date
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
		)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		student.ID,
		student.FirstNames,
		student.LastNames,
		student.DocumentID,
		student.BirthDate,
		student.ProfilePhotoURL,
		student.Gender,
		student.NationalityCountryID,
		student.ResidenceCountryID,
		student.ResidenceCityID,
		student.Emails,
		student.Phones,
		student.CompanyID,
		student.JobTitleCategoryID,
		student.ProfessionID,
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
			id, first_names, last_names, document_id, birth_date, profile_photo_url,
			gender, nationality_country_id, residence_country_id, residence_city_id,
			emails, phones, company_id, job_title_category_id, profession_id,
			student_code, status, cohort, enrollment_date, graduation_date,
			created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM students
		WHERE id = $1 AND deleted_at IS NULL
	`

	student := &models.Student{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&student.ID,
		&student.FirstNames,
		&student.LastNames,
		&student.DocumentID,
		&student.BirthDate,
		&student.ProfilePhotoURL,
		&student.Gender,
		&student.NationalityCountryID,
		&student.ResidenceCountryID,
		&student.ResidenceCityID,
		&student.Emails,
		&student.Phones,
		&student.CompanyID,
		&student.JobTitleCategoryID,
		&student.ProfessionID,
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
			id, first_names, last_names, document_id, birth_date, profile_photo_url,
			gender, nationality_country_id, residence_country_id, residence_city_id,
			emails, phones, company_id, job_title_category_id, profession_id,
			student_code, status, cohort, enrollment_date, graduation_date,
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

	if filters.ResidenceCountryID != nil {
		query += fmt.Sprintf(" AND residence_country_id = $%d", argCount)
		args = append(args, *filters.ResidenceCountryID)
		argCount++
	}

	if filters.Search != nil {
		query += fmt.Sprintf(" AND (first_names ILIKE $%d OR last_names ILIKE $%d)", argCount, argCount)
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
			&student.FirstNames,
			&student.LastNames,
			&student.DocumentID,
			&student.BirthDate,
			&student.ProfilePhotoURL,
			&student.Gender,
			&student.NationalityCountryID,
			&student.ResidenceCountryID,
			&student.ResidenceCityID,
			&student.Emails,
			&student.Phones,
			&student.CompanyID,
			&student.JobTitleCategoryID,
			&student.ProfessionID,
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
			first_names = $2,
			last_names = $3,
			document_id = $4,
			profile_photo_url = $5,
			gender = $6,
			nationality_country_id = $7,
			residence_country_id = $8,
			residence_city_id = $9,
			emails = $10,
			phones = $11,
			company_id = $12,
			job_title_category_id = $13,
			profession_id = $14,
			student_code = $15,
			status = $16,
			updated_by = $17
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		student.ID,
		student.FirstNames,
		student.LastNames,
		student.DocumentID,
		student.ProfilePhotoURL,
		student.Gender,
		student.NationalityCountryID,
		student.ResidenceCountryID,
		student.ResidenceCityID,
		student.Emails,
		student.Phones,
		student.CompanyID,
		student.JobTitleCategoryID,
		student.ProfessionID,
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

	if filters.ResidenceCountryID != nil {
		query += fmt.Sprintf(" AND residence_country_id = $%d", argCount)
		args = append(args, *filters.ResidenceCountryID)
		argCount++
	}

	if filters.Search != nil {
		query += fmt.Sprintf(" AND (first_names ILIKE $%d OR last_names ILIKE $%d)", argCount, argCount)
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
