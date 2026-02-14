package models

import (
	"time"

	"github.com/google/uuid"
)

// StudentStatus represents the academic status of a student.
type StudentStatus string

const (
	StudentStatusActive    StudentStatus = "active"
	StudentStatusGraduated StudentStatus = "graduated"
	StudentStatusWithdrawn StudentStatus = "withdrawn"
	StudentStatusSuspended StudentStatus = "suspended"
)

// Student maps to the students table.
type Student struct {
	ID              uuid.UUID     `json:"id" db:"id"`
	FirstNames      string        `json:"first_names" db:"first_names"`
	LastNames       string        `json:"last_names" db:"last_names"`
	DocumentID      *string       `json:"document_id,omitempty" db:"document_id"`
	BirthDate       time.Time     `json:"birth_date" db:"birth_date"`
	ProfilePhotoURL *string       `json:"profile_photo_url,omitempty" db:"profile_photo_url"`
	Gender          *string       `json:"gender,omitempty" db:"gender"`

	// Ubicación
	NationalityCountryID uuid.UUID  `json:"nationality_country_id" db:"nationality_country_id"`
	ResidenceCountryID   uuid.UUID  `json:"residence_country_id" db:"residence_country_id"`
	ResidenceCityID      *uuid.UUID `json:"residence_city_id,omitempty" db:"residence_city_id"`

	// Contacto (PostgreSQL arrays)
	Emails []string `json:"emails" db:"emails"`
	Phones []string `json:"phones" db:"phones"`

	// Laboral
	CompanyID          *uuid.UUID `json:"company_id,omitempty" db:"company_id"`
	JobTitleCategoryID *uuid.UUID `json:"job_title_category_id,omitempty" db:"job_title_category_id"`
	ProfessionID       *uuid.UUID `json:"profession_id,omitempty" db:"profession_id"`

	// Estado academico
	StudentCode    *string       `json:"student_code,omitempty" db:"student_code"`
	Status         StudentStatus `json:"status" db:"status"`
	Cohort         string        `json:"cohort" db:"cohort"`
	EnrollmentDate time.Time     `json:"enrollment_date" db:"enrollment_date"`
	GraduationDate *time.Time    `json:"graduation_date,omitempty" db:"graduation_date"`

	// Auditoria
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	CreatedBy *uuid.UUID `json:"created_by,omitempty" db:"created_by"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	DeletedBy *uuid.UUID `json:"deleted_by,omitempty" db:"deleted_by"`
}

// CreateStudentRequest is the DTO for creating a student.
type CreateStudentRequest struct {
	FirstNames           string   `json:"first_names" validate:"required,min=2,max=150"`
	LastNames            string   `json:"last_names" validate:"required,min=2,max=150"`
	DocumentID           *string  `json:"document_id" validate:"omitempty,max=50"`
	BirthDate            string   `json:"birth_date" validate:"required"`
	ProfilePhotoURL      *string  `json:"profile_photo_url" validate:"omitempty,url"`
	Gender               *string  `json:"gender" validate:"omitempty,oneof=M F"`
	NationalityCountryID string   `json:"nationality_country_id" validate:"required,uuid"`
	ResidenceCountryID   string   `json:"residence_country_id" validate:"required,uuid"`
	ResidenceCityID      *string  `json:"residence_city_id" validate:"omitempty,uuid"`
	Emails               []string `json:"emails" validate:"required,min=1,dive,email"`
	Phones               []string `json:"phones" validate:"omitempty,dive,max=50"`
	CompanyID            *string  `json:"company_id" validate:"omitempty,uuid"`
	JobTitleCategoryID   *string  `json:"job_title_category_id" validate:"omitempty,uuid"`
	ProfessionID         *string  `json:"profession_id" validate:"omitempty,uuid"`
	StudentCode          *string  `json:"student_code" validate:"omitempty,len=9"`
	Status               string   `json:"status" validate:"required,oneof=active graduated withdrawn suspended"`
	Cohort               string   `json:"cohort" validate:"required,max=10"`
	EnrollmentDate       string   `json:"enrollment_date" validate:"required"`
}

// UpdateStudentRequest is the DTO for updating a student. All fields are optional.
type UpdateStudentRequest struct {
	FirstNames         *string  `json:"first_names" validate:"omitempty,min=2,max=150"`
	LastNames          *string  `json:"last_names" validate:"omitempty,min=2,max=150"`
	DocumentID         *string  `json:"document_id" validate:"omitempty,max=50"`
	ProfilePhotoURL    *string  `json:"profile_photo_url" validate:"omitempty,url"`
	Gender             *string  `json:"gender" validate:"omitempty,oneof=M F"`
	Emails             []string `json:"emails" validate:"omitempty,min=1,dive,email"`
	Phones             []string `json:"phones" validate:"omitempty"`
	CompanyID          *string  `json:"company_id" validate:"omitempty,uuid"`
	JobTitleCategoryID *string  `json:"job_title_category_id" validate:"omitempty,uuid"`
	ProfessionID       *string  `json:"profession_id" validate:"omitempty,uuid"`
	StudentCode        *string  `json:"student_code" validate:"omitempty,len=9"`
	Status             *string  `json:"status" validate:"omitempty,oneof=active graduated withdrawn suspended"`
}

// ImportRowError describes a validation or insertion error for a single row.
type ImportRowError struct {
	Row     int    `json:"row"`
	Field   string `json:"field"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// ImportResult holds the outcome of a bulk student import.
type ImportResult struct {
	TotalRows int              `json:"total_rows"`
	Created   int              `json:"created"`
	Errors    []ImportRowError `json:"errors"`
}
