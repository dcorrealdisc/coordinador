package handlers

import (
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/dcorreal/coordinador/internal/models"
	"github.com/dcorreal/coordinador/internal/repositories"
	"github.com/dcorreal/coordinador/internal/services"
	"github.com/dcorreal/coordinador/internal/shared"
)

// StudentHandler handles HTTP requests for student endpoints.
type StudentHandler struct {
	studentService       services.StudentService
	studentImportService services.StudentImportService
}

// NewStudentHandler creates a new StudentHandler.
func NewStudentHandler(studentService services.StudentService, studentImportService services.StudentImportService) *StudentHandler {
	return &StudentHandler{
		studentService:       studentService,
		studentImportService: studentImportService,
	}
}

// RegisterRoutes registers all student routes on the given router group.
func (h *StudentHandler) RegisterRoutes(router fiber.Router) {
	students := router.Group("/students")

	students.Post("/", h.CreateStudent)
	students.Post("/import", h.ImportStudents)
	students.Get("/", h.ListStudents)
	students.Get("/:id", h.GetStudent)
	students.Put("/:id", h.UpdateStudent)
	students.Delete("/:id", h.DeleteStudent)
}

// CreateStudent handles POST /api/v1/students
func (h *StudentHandler) CreateStudent(c *fiber.Ctx) error {
	var req models.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return shared.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// TODO: Get authenticated user from context once auth is implemented
	var createdBy *uuid.UUID // nil until auth is implemented

	student, err := h.studentService.CreateStudent(c.Context(), &req, createdBy)
	if err != nil {
		return shared.ErrorResponse(c, fiber.StatusBadRequest, "Failed to create student", err)
	}

	return shared.SuccessResponse(c, fiber.StatusCreated, "Student created successfully", student)
}

// GetStudent handles GET /api/v1/students/:id
func (h *StudentHandler) GetStudent(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return shared.ErrorResponse(c, fiber.StatusBadRequest, "Invalid student ID", err)
	}

	student, err := h.studentService.GetStudent(c.Context(), id)
	if err != nil {
		return shared.ErrorResponse(c, fiber.StatusNotFound, "Student not found", err)
	}

	return shared.SuccessResponse(c, fiber.StatusOK, "Student retrieved successfully", student)
}

// ListStudents handles GET /api/v1/students
func (h *StudentHandler) ListStudents(c *fiber.Ctx) error {
	filters := repositories.StudentFilters{}

	if status := c.Query("status"); status != "" {
		filters.Status = &status
	}
	if cohort := c.Query("cohort"); cohort != "" {
		filters.Cohort = &cohort
	}
	if search := c.Query("search"); search != "" {
		filters.Search = &search
	}
	if countryID := c.Query("residence_country_id"); countryID != "" {
		parsed, err := uuid.Parse(countryID)
		if err == nil {
			filters.ResidenceCountryID = &parsed
		}
	}

	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	filters.Limit = limit
	filters.Offset = offset

	students, total, err := h.studentService.ListStudents(c.Context(), filters)
	if err != nil {
		return shared.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to list students", err)
	}

	return shared.PaginatedResponse(c, fiber.StatusOK, "Students retrieved successfully", students, total, limit, offset)
}

// UpdateStudent handles PUT /api/v1/students/:id
func (h *StudentHandler) UpdateStudent(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return shared.ErrorResponse(c, fiber.StatusBadRequest, "Invalid student ID", err)
	}

	var req models.UpdateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return shared.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// TODO: Get authenticated user from context once auth is implemented
	var updatedBy *uuid.UUID // nil until auth is implemented

	student, err := h.studentService.UpdateStudent(c.Context(), id, &req, updatedBy)
	if err != nil {
		return shared.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update student", err)
	}

	return shared.SuccessResponse(c, fiber.StatusOK, "Student updated successfully", student)
}

// DeleteStudent handles DELETE /api/v1/students/:id
func (h *StudentHandler) DeleteStudent(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return shared.ErrorResponse(c, fiber.StatusBadRequest, "Invalid student ID", err)
	}

	// TODO: Get authenticated user from context once auth is implemented
	var deletedBy *uuid.UUID // nil until auth is implemented

	if err := h.studentService.DeleteStudent(c.Context(), id, deletedBy); err != nil {
		return shared.ErrorResponse(c, fiber.StatusNotFound, "Failed to delete student", err)
	}

	return shared.SuccessResponse(c, fiber.StatusOK, "Student deleted successfully", nil)
}

// ImportStudents handles POST /api/v1/students/import
func (h *StudentHandler) ImportStudents(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return shared.ErrorResponse(c, fiber.StatusBadRequest, "File is required", fmt.Errorf("missing 'file' field in multipart form"))
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	var format string
	switch ext {
	case ".csv":
		format = "csv"
	case ".xlsx":
		format = "xlsx"
	default:
		return shared.ErrorResponse(c, fiber.StatusBadRequest, "Unsupported file format", fmt.Errorf("expected .csv or .xlsx, got %s", ext))
	}

	file, err := fileHeader.Open()
	if err != nil {
		return shared.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to open file", err)
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return shared.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to read file", err)
	}

	// TODO: Get authenticated user from context once auth is implemented
	var createdBy *uuid.UUID

	result, err := h.studentImportService.ImportFromFile(c.Context(), fileData, format, createdBy)
	if err != nil {
		return shared.ErrorResponse(c, fiber.StatusBadRequest, "Import failed", err)
	}

	message := fmt.Sprintf("Import completed: %d created, %d errors", result.Created, len(result.Errors))
	return shared.SuccessResponse(c, fiber.StatusOK, message, result)
}
