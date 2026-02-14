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

// Student maps to the students table defined in migration 003.
type Student struct {
	ID              uuid.UUID     `json:"id" db:"id"`
	FullName        string        `json:"full_name" db:"full_name"`
	DocumentID      *string       `json:"document_id,omitempty" db:"document_id"`
	BirthDate       time.Time     `json:"birth_date" db:"birth_date"`
	ProfilePhotoURL *string       `json:"profile_photo_url,omitempty" db:"profile_photo_url"`

	// Procedencia
	CityOriginID    *uuid.UUID `json:"city_origin_id,omitempty" db:"city_origin_id"`
	CountryOriginID uuid.UUID  `json:"country_origin_id" db:"country_origin_id"`

	// Contacto (PostgreSQL arrays)
	Emails []string `json:"emails" db:"emails"`
	Phones []string `json:"phones" db:"phones"`

	// Laboral
	CompanyID *uuid.UUID `json:"company_id,omitempty" db:"company_id"`

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
	FullName        string   `json:"full_name" validate:"required,min=3,max=255"`
	DocumentID      *string  `json:"document_id" validate:"omitempty,max=50"`
	BirthDate       string   `json:"birth_date" validate:"required"`
	ProfilePhotoURL *string  `json:"profile_photo_url" validate:"omitempty,url"`
	CityOriginID    *string  `json:"city_origin_id" validate:"omitempty,uuid"`
	CountryOriginID string   `json:"country_origin_id" validate:"required,uuid"`
	Emails          []string `json:"emails" validate:"required,min=1,dive,email"`
	Phones          []string `json:"phones" validate:"omitempty,dive,max=50"`
	CompanyID       *string  `json:"company_id" validate:"omitempty,uuid"`
	StudentCode     *string  `json:"student_code" validate:"omitempty,len=9"`
	Status          string   `json:"status" validate:"required,oneof=active graduated withdrawn suspended"`
	Cohort          string   `json:"cohort" validate:"required,max=10"`
	EnrollmentDate  string   `json:"enrollment_date" validate:"required"`
}

// UpdateStudentRequest is the DTO for updating a student. All fields are optional.
type UpdateStudentRequest struct {
	FullName        *string  `json:"full_name" validate:"omitempty,min=3,max=255"`
	DocumentID      *string  `json:"document_id" validate:"omitempty,max=50"`
	ProfilePhotoURL *string  `json:"profile_photo_url" validate:"omitempty,url"`
	Emails          []string `json:"emails" validate:"omitempty,min=1,dive,email"`
	Phones          []string `json:"phones" validate:"omitempty"`
	CompanyID       *string  `json:"company_id" validate:"omitempty,uuid"`
	StudentCode     *string  `json:"student_code" validate:"omitempty,len=9"`
	Status          *string  `json:"status" validate:"omitempty,oneof=active graduated withdrawn suspended"`
}
