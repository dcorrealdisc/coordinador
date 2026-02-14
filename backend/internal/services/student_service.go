package services

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"

	"github.com/dcorreal/coordinador/internal/models"
	"github.com/dcorreal/coordinador/internal/repositories"
)

// StudentService defines the business logic interface for students.
type StudentService interface {
	CreateStudent(ctx context.Context, req *models.CreateStudentRequest, createdBy *uuid.UUID) (*models.Student, error)
	GetStudent(ctx context.Context, id uuid.UUID) (*models.Student, error)
	ListStudents(ctx context.Context, filters repositories.StudentFilters) ([]*models.Student, int, error)
	UpdateStudent(ctx context.Context, id uuid.UUID, req *models.UpdateStudentRequest, updatedBy *uuid.UUID) (*models.Student, error)
	DeleteStudent(ctx context.Context, id uuid.UUID, deletedBy *uuid.UUID) error
}

var studentCodeRegex = regexp.MustCompile(`^[0-9]{4}[12][0-9]{4}$`)

type studentService struct {
	studentRepo repositories.StudentRepository
}

// NewStudentService creates a new StudentService.
func NewStudentService(studentRepo repositories.StudentRepository) StudentService {
	return &studentService{studentRepo: studentRepo}
}

func (s *studentService) CreateStudent(ctx context.Context, req *models.CreateStudentRequest, createdBy *uuid.UUID) (*models.Student, error) {
	// Parse and validate birth date
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return nil, fmt.Errorf("invalid birth_date format, expected YYYY-MM-DD: %w", err)
	}

	// Validate minimum age (18 years)
	age := time.Now().Year() - birthDate.Year()
	if time.Now().YearDay() < birthDate.YearDay() {
		age--
	}
	if age < 18 {
		return nil, fmt.Errorf("student must be at least 18 years old")
	}

	// Parse enrollment date
	enrollmentDate, err := time.Parse("2006-01-02", req.EnrollmentDate)
	if err != nil {
		return nil, fmt.Errorf("invalid enrollment_date format, expected YYYY-MM-DD: %w", err)
	}

	// Parse country origin ID (required)
	countryOriginID, err := uuid.Parse(req.CountryOriginID)
	if err != nil {
		return nil, fmt.Errorf("invalid country_origin_id: %w", err)
	}

	// Parse optional UUIDs
	var cityOriginID *uuid.UUID
	if req.CityOriginID != nil {
		parsed, err := uuid.Parse(*req.CityOriginID)
		if err != nil {
			return nil, fmt.Errorf("invalid city_origin_id: %w", err)
		}
		cityOriginID = &parsed
	}

	var companyID *uuid.UUID
	if req.CompanyID != nil {
		parsed, err := uuid.Parse(*req.CompanyID)
		if err != nil {
			return nil, fmt.Errorf("invalid company_id: %w", err)
		}
		companyID = &parsed
	}

	// Validate student_code format if provided
	if req.StudentCode != nil {
		if !studentCodeRegex.MatchString(*req.StudentCode) {
			return nil, fmt.Errorf("invalid student_code format, expected YYYYS#### (e.g. 202620190)")
		}
	}

	student := &models.Student{
		ID:              uuid.New(),
		FullName:        req.FullName,
		DocumentID:      req.DocumentID,
		BirthDate:       birthDate,
		ProfilePhotoURL: req.ProfilePhotoURL,
		CityOriginID:    cityOriginID,
		CountryOriginID: countryOriginID,
		Emails:          req.Emails,
		Phones:          req.Phones,
		CompanyID:       companyID,
		StudentCode:     req.StudentCode,
		Status:          models.StudentStatus(req.Status),
		Cohort:          req.Cohort,
		EnrollmentDate:  enrollmentDate,
		CreatedBy:       createdBy,
	}

	if err := s.studentRepo.Create(ctx, student); err != nil {
		return nil, fmt.Errorf("failed to create student: %w", err)
	}

	return student, nil
}

func (s *studentService) GetStudent(ctx context.Context, id uuid.UUID) (*models.Student, error) {
	return s.studentRepo.GetByID(ctx, id)
}

func (s *studentService) ListStudents(ctx context.Context, filters repositories.StudentFilters) ([]*models.Student, int, error) {
	students, err := s.studentRepo.List(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.studentRepo.Count(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	return students, count, nil
}

func (s *studentService) UpdateStudent(ctx context.Context, id uuid.UUID, req *models.UpdateStudentRequest, updatedBy *uuid.UUID) (*models.Student, error) {
	student, err := s.studentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply partial updates
	if req.FullName != nil {
		student.FullName = *req.FullName
	}
	if req.DocumentID != nil {
		student.DocumentID = req.DocumentID
	}
	if req.ProfilePhotoURL != nil {
		student.ProfilePhotoURL = req.ProfilePhotoURL
	}
	if req.Emails != nil {
		if len(req.Emails) == 0 {
			return nil, fmt.Errorf("at least one email is required")
		}
		student.Emails = req.Emails
	}
	if req.Phones != nil {
		student.Phones = req.Phones
	}
	if req.CompanyID != nil {
		parsed, err := uuid.Parse(*req.CompanyID)
		if err != nil {
			return nil, fmt.Errorf("invalid company_id: %w", err)
		}
		student.CompanyID = &parsed
	}
	if req.StudentCode != nil {
		if !studentCodeRegex.MatchString(*req.StudentCode) {
			return nil, fmt.Errorf("invalid student_code format, expected YYYYS#### (e.g. 202620190)")
		}
		student.StudentCode = req.StudentCode
	}
	if req.Status != nil {
		student.Status = models.StudentStatus(*req.Status)
	}

	student.UpdatedBy = updatedBy

	if err := s.studentRepo.Update(ctx, student); err != nil {
		return nil, fmt.Errorf("failed to update student: %w", err)
	}

	return student, nil
}

func (s *studentService) DeleteStudent(ctx context.Context, id uuid.UUID, deletedBy *uuid.UUID) error {
	return s.studentRepo.Delete(ctx, id, deletedBy)
}
