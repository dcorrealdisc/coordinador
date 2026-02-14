package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CatalogRepository provides find-or-create access to catalog tables.
type CatalogRepository interface {
	FindCountryByName(ctx context.Context, name string) (uuid.UUID, error)
	CreateCountry(ctx context.Context, name string) (uuid.UUID, error)

	FindCityByName(ctx context.Context, name string, countryID uuid.UUID) (uuid.UUID, error)
	CreateCity(ctx context.Context, name string, countryID uuid.UUID) (uuid.UUID, error)

	FindProfessionByName(ctx context.Context, name string) (uuid.UUID, error)
	CreateProfession(ctx context.Context, name string) (uuid.UUID, error)

	FindJobTitleCategoryByName(ctx context.Context, name string) (uuid.UUID, error)
	CreateJobTitleCategory(ctx context.Context, name string) (uuid.UUID, error)

	FindUniversityByName(ctx context.Context, name string, countryID uuid.UUID) (uuid.UUID, error)
	CreateUniversity(ctx context.Context, name string, cityID *uuid.UUID, countryID uuid.UUID) (uuid.UUID, error)

	CreateStudentUniversity(ctx context.Context, studentID, universityID uuid.UUID) error
}

type catalogRepository struct {
	db *pgxpool.Pool
}

// NewCatalogRepository creates a new CatalogRepository.
func NewCatalogRepository(db *pgxpool.Pool) CatalogRepository {
	return &catalogRepository{db: db}
}

func (r *catalogRepository) FindCountryByName(ctx context.Context, name string) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRow(ctx,
		"SELECT id FROM countries WHERE LOWER(unaccent(TRIM(name))) = LOWER(unaccent(TRIM($1)))", name,
	).Scan(&id)
	if err == pgx.ErrNoRows {
		return uuid.Nil, nil
	}
	return id, err
}

func (r *catalogRepository) CreateCountry(ctx context.Context, name string) (uuid.UUID, error) {
	// Generate a 3-letter code from the country name (uppercase, first 3 chars)
	code := strings.ToUpper(strings.TrimSpace(name))
	if len(code) > 3 {
		code = code[:3]
	}

	var id uuid.UUID
	err := r.db.QueryRow(ctx,
		"INSERT INTO countries (code, name) VALUES ($1, $2) RETURNING id", code, name,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create country %q: %w", name, err)
	}
	return id, nil
}

func (r *catalogRepository) FindCityByName(ctx context.Context, name string, countryID uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRow(ctx,
		"SELECT id FROM cities WHERE LOWER(unaccent(TRIM(name))) = LOWER(unaccent(TRIM($1))) AND country_id = $2", name, countryID,
	).Scan(&id)
	if err == pgx.ErrNoRows {
		return uuid.Nil, nil
	}
	return id, err
}

func (r *catalogRepository) CreateCity(ctx context.Context, name string, countryID uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRow(ctx,
		"INSERT INTO cities (name, country_id) VALUES ($1, $2) RETURNING id", name, countryID,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create city %q: %w", name, err)
	}
	return id, nil
}

func (r *catalogRepository) FindProfessionByName(ctx context.Context, name string) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRow(ctx,
		"SELECT id FROM professions WHERE LOWER(unaccent(TRIM(name))) = LOWER(unaccent(TRIM($1)))", name,
	).Scan(&id)
	if err == pgx.ErrNoRows {
		return uuid.Nil, nil
	}
	return id, err
}

func (r *catalogRepository) CreateProfession(ctx context.Context, name string) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRow(ctx,
		"INSERT INTO professions (name) VALUES ($1) RETURNING id", name,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create profession %q: %w", name, err)
	}
	return id, nil
}

func (r *catalogRepository) FindJobTitleCategoryByName(ctx context.Context, name string) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRow(ctx,
		"SELECT id FROM job_title_categories WHERE LOWER(unaccent(TRIM(name))) = LOWER(unaccent(TRIM($1)))", name,
	).Scan(&id)
	if err == pgx.ErrNoRows {
		return uuid.Nil, nil
	}
	return id, err
}

func (r *catalogRepository) CreateJobTitleCategory(ctx context.Context, name string) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRow(ctx,
		"INSERT INTO job_title_categories (name) VALUES ($1) RETURNING id", name,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create job_title_category %q: %w", name, err)
	}
	return id, nil
}

func (r *catalogRepository) FindUniversityByName(ctx context.Context, name string, countryID uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRow(ctx,
		"SELECT id FROM universities WHERE LOWER(unaccent(TRIM(name))) = LOWER(unaccent(TRIM($1))) AND country_id = $2", name, countryID,
	).Scan(&id)
	if err == pgx.ErrNoRows {
		return uuid.Nil, nil
	}
	return id, err
}

func (r *catalogRepository) CreateUniversity(ctx context.Context, name string, cityID *uuid.UUID, countryID uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRow(ctx,
		"INSERT INTO universities (name, city_id, country_id) VALUES ($1, $2, $3) RETURNING id",
		name, cityID, countryID,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create university %q: %w", name, err)
	}
	return id, nil
}

func (r *catalogRepository) CreateStudentUniversity(ctx context.Context, studentID, universityID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		"INSERT INTO student_universities (student_id, university_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		studentID, universityID,
	)
	if err != nil {
		return fmt.Errorf("failed to create student_university: %w", err)
	}
	return nil
}
