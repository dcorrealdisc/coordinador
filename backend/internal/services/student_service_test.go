package services_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/dcorreal/coordinador/internal/models"
	"github.com/dcorreal/coordinador/internal/repositories"
	"github.com/dcorreal/coordinador/internal/repositories/mocks"
	"github.com/dcorreal/coordinador/internal/services"
)

func validCreateRequest() *models.CreateStudentRequest {
	return &models.CreateStudentRequest{
		FirstNames:           "Juan Carlos",
		LastNames:            "Perez",
		BirthDate:            "1995-03-15",
		NationalityCountryID: uuid.New().String(),
		ResidenceCountryID:   uuid.New().String(),
		Emails:               []string{"juan@test.com"},
		Status:               "active",
		Cohort:               "2024-1",
		EnrollmentDate:       "2024-01-15",
	}
}

func sampleStudent() *models.Student {
	return &models.Student{
		ID:                   uuid.New(),
		FirstNames:           "Juan Carlos",
		LastNames:            "Perez",
		BirthDate:            time.Date(1995, 3, 15, 0, 0, 0, 0, time.UTC),
		NationalityCountryID: uuid.New(),
		ResidenceCountryID:   uuid.New(),
		Emails:               []string{"juan@test.com"},
		Status:               models.StudentStatusActive,
		Cohort:               "2024-1",
		EnrollmentDate:       time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
}

// =============================================================================
// CreateStudent
// =============================================================================

func TestCreateStudent_Success(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	req := validCreateRequest()
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	student, err := service.CreateStudent(context.Background(), req, nil)

	assert.NoError(t, err)
	assert.NotNil(t, student)
	assert.Equal(t, "Juan Carlos", student.FirstNames)
	assert.Equal(t, "Perez", student.LastNames)
	assert.Equal(t, models.StudentStatusActive, student.Status)
	assert.Equal(t, "2024-1", student.Cohort)
	assert.NotEqual(t, uuid.Nil, student.ID)
	mockRepo.AssertExpectations(t)
}

func TestCreateStudent_InvalidBirthDateFormat(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	req := validCreateRequest()
	req.BirthDate = "15-03-1995" // wrong format

	student, err := service.CreateStudent(context.Background(), req, nil)

	assert.Error(t, err)
	assert.Nil(t, student)
	assert.Contains(t, err.Error(), "birth_date")
}

func TestCreateStudent_UnderAge(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	req := validCreateRequest()
	req.BirthDate = time.Now().AddDate(-17, 0, 0).Format("2006-01-02") // 17 years old

	student, err := service.CreateStudent(context.Background(), req, nil)

	assert.Error(t, err)
	assert.Nil(t, student)
	assert.Contains(t, err.Error(), "18 years old")
}

func TestCreateStudent_InvalidEnrollmentDateFormat(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	req := validCreateRequest()
	req.EnrollmentDate = "not-a-date"

	student, err := service.CreateStudent(context.Background(), req, nil)

	assert.Error(t, err)
	assert.Nil(t, student)
	assert.Contains(t, err.Error(), "enrollment_date")
}

func TestCreateStudent_InvalidNationalityCountryID(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	req := validCreateRequest()
	req.NationalityCountryID = "not-a-uuid"

	student, err := service.CreateStudent(context.Background(), req, nil)

	assert.Error(t, err)
	assert.Nil(t, student)
	assert.Contains(t, err.Error(), "nationality_country_id")
}

func TestCreateStudent_InvalidResidenceCountryID(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	req := validCreateRequest()
	req.ResidenceCountryID = "not-a-uuid"

	student, err := service.CreateStudent(context.Background(), req, nil)

	assert.Error(t, err)
	assert.Nil(t, student)
	assert.Contains(t, err.Error(), "residence_country_id")
}

func TestCreateStudent_InvalidResidenceCityID(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	req := validCreateRequest()
	badID := "not-a-uuid"
	req.ResidenceCityID = &badID

	student, err := service.CreateStudent(context.Background(), req, nil)

	assert.Error(t, err)
	assert.Nil(t, student)
	assert.Contains(t, err.Error(), "residence_city_id")
}

func TestCreateStudent_InvalidCompanyID(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	req := validCreateRequest()
	badID := "not-a-uuid"
	req.CompanyID = &badID

	student, err := service.CreateStudent(context.Background(), req, nil)

	assert.Error(t, err)
	assert.Nil(t, student)
	assert.Contains(t, err.Error(), "company_id")
}

func TestCreateStudent_WithStudentCode(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	req := validCreateRequest()
	code := "202620190"
	req.StudentCode = &code
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	student, err := service.CreateStudent(context.Background(), req, nil)

	assert.NoError(t, err)
	assert.NotNil(t, student)
	assert.Equal(t, &code, student.StudentCode)
	mockRepo.AssertExpectations(t)
}

func TestCreateStudent_InvalidStudentCodeFormat(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	req := validCreateRequest()
	badCode := "ABC123456"
	req.StudentCode = &badCode

	student, err := service.CreateStudent(context.Background(), req, nil)

	assert.Error(t, err)
	assert.Nil(t, student)
	assert.Contains(t, err.Error(), "student_code")
}

func TestCreateStudent_InvalidStudentCodeSemester(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	req := validCreateRequest()
	badCode := "202630190" // semester 3 is invalid
	req.StudentCode = &badCode

	student, err := service.CreateStudent(context.Background(), req, nil)

	assert.Error(t, err)
	assert.Nil(t, student)
	assert.Contains(t, err.Error(), "student_code")
}

func TestCreateStudent_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	req := validCreateRequest()
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(fmt.Errorf("db connection failed"))

	student, err := service.CreateStudent(context.Background(), req, nil)

	assert.Error(t, err)
	assert.Nil(t, student)
	assert.Contains(t, err.Error(), "failed to create student")
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// GetStudent
// =============================================================================

func TestGetStudent_Success(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	expected := sampleStudent()
	mockRepo.On("GetByID", mock.Anything, expected.ID).Return(expected, nil)

	student, err := service.GetStudent(context.Background(), expected.ID)

	assert.NoError(t, err)
	assert.Equal(t, expected.ID, student.ID)
	assert.Equal(t, expected.FirstNames, student.FirstNames)
	mockRepo.AssertExpectations(t)
}

func TestGetStudent_NotFound(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	id := uuid.New()
	mockRepo.On("GetByID", mock.Anything, id).Return(nil, fmt.Errorf("student not found"))

	student, err := service.GetStudent(context.Background(), id)

	assert.Error(t, err)
	assert.Nil(t, student)
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// ListStudents
// =============================================================================

func TestListStudents_Success(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	filters := repositories.StudentFilters{Limit: 20, Offset: 0}
	expected := []*models.Student{sampleStudent(), sampleStudent()}

	mockRepo.On("List", mock.Anything, filters).Return(expected, nil)
	mockRepo.On("Count", mock.Anything, filters).Return(2, nil)

	students, total, err := service.ListStudents(context.Background(), filters)

	assert.NoError(t, err)
	assert.Len(t, students, 2)
	assert.Equal(t, 2, total)
	mockRepo.AssertExpectations(t)
}

func TestListStudents_Empty(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	filters := repositories.StudentFilters{Limit: 20, Offset: 0}

	mockRepo.On("List", mock.Anything, filters).Return([]*models.Student{}, nil)
	mockRepo.On("Count", mock.Anything, filters).Return(0, nil)

	students, total, err := service.ListStudents(context.Background(), filters)

	assert.NoError(t, err)
	assert.Empty(t, students)
	assert.Equal(t, 0, total)
	mockRepo.AssertExpectations(t)
}

func TestListStudents_ListError(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	filters := repositories.StudentFilters{}
	mockRepo.On("List", mock.Anything, filters).Return(nil, fmt.Errorf("db error"))

	students, total, err := service.ListStudents(context.Background(), filters)

	assert.Error(t, err)
	assert.Nil(t, students)
	assert.Equal(t, 0, total)
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// UpdateStudent
// =============================================================================

func TestUpdateStudent_Success(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	existing := sampleStudent()
	newFirstNames := "Juan Actualizado"
	req := &models.UpdateStudentRequest{
		FirstNames: &newFirstNames,
	}

	mockRepo.On("GetByID", mock.Anything, existing.ID).Return(existing, nil)
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	student, err := service.UpdateStudent(context.Background(), existing.ID, req, nil)

	assert.NoError(t, err)
	assert.Equal(t, "Juan Actualizado", student.FirstNames)
	mockRepo.AssertExpectations(t)
}

func TestUpdateStudent_NotFound(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	id := uuid.New()
	newName := "Inexistente"
	req := &models.UpdateStudentRequest{FirstNames: &newName}

	mockRepo.On("GetByID", mock.Anything, id).Return(nil, fmt.Errorf("student not found"))

	student, err := service.UpdateStudent(context.Background(), id, req, nil)

	assert.Error(t, err)
	assert.Nil(t, student)
	mockRepo.AssertExpectations(t)
}

func TestUpdateStudent_EmptyEmails(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	existing := sampleStudent()
	req := &models.UpdateStudentRequest{
		Emails: []string{}, // empty - should fail
	}

	mockRepo.On("GetByID", mock.Anything, existing.ID).Return(existing, nil)

	student, err := service.UpdateStudent(context.Background(), existing.ID, req, nil)

	assert.Error(t, err)
	assert.Nil(t, student)
	assert.Contains(t, err.Error(), "at least one email")
}

func TestUpdateStudent_WithStudentCode(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	existing := sampleStudent()
	code := "202510001"
	req := &models.UpdateStudentRequest{
		StudentCode: &code,
	}

	mockRepo.On("GetByID", mock.Anything, existing.ID).Return(existing, nil)
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	student, err := service.UpdateStudent(context.Background(), existing.ID, req, nil)

	assert.NoError(t, err)
	assert.Equal(t, &code, student.StudentCode)
	mockRepo.AssertExpectations(t)
}

func TestUpdateStudent_InvalidStudentCode(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	existing := sampleStudent()
	badCode := "12345"
	req := &models.UpdateStudentRequest{
		StudentCode: &badCode,
	}

	mockRepo.On("GetByID", mock.Anything, existing.ID).Return(existing, nil)

	student, err := service.UpdateStudent(context.Background(), existing.ID, req, nil)

	assert.Error(t, err)
	assert.Nil(t, student)
	assert.Contains(t, err.Error(), "student_code")
}

func TestUpdateStudent_PartialFields(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	existing := sampleStudent()
	newStatus := "graduated"
	req := &models.UpdateStudentRequest{
		Status: &newStatus,
		Emails: []string{"new@email.com", "other@email.com"},
	}

	mockRepo.On("GetByID", mock.Anything, existing.ID).Return(existing, nil)
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	student, err := service.UpdateStudent(context.Background(), existing.ID, req, nil)

	assert.NoError(t, err)
	assert.Equal(t, models.StudentStatus("graduated"), student.Status)
	assert.Equal(t, []string{"new@email.com", "other@email.com"}, student.Emails)
	// Names should remain unchanged
	assert.Equal(t, "Juan Carlos", student.FirstNames)
	assert.Equal(t, "Perez", student.LastNames)
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// DeleteStudent
// =============================================================================

func TestDeleteStudent_Success(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	id := uuid.New()
	mockRepo.On("Delete", mock.Anything, id, (*uuid.UUID)(nil)).Return(nil)

	err := service.DeleteStudent(context.Background(), id, nil)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteStudent_NotFound(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	service := services.NewStudentService(mockRepo)

	id := uuid.New()
	mockRepo.On("Delete", mock.Anything, id, (*uuid.UUID)(nil)).Return(fmt.Errorf("student not found"))

	err := service.DeleteStudent(context.Background(), id, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student not found")
	mockRepo.AssertExpectations(t)
}
