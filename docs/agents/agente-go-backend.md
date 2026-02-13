# Agente Go/Backend - Guía de Trabajo

## 🎯 Rol y Responsabilidades

El Agente Go/Backend es responsable de la implementación del servidor backend del sistema Coordinador usando Go y el framework Fiber. Su misión es construir APIs REST robustas, mantener código limpio y testeable, y seguir las mejores prácticas de Go.

### Responsabilidades Principales

1. **Implementación de APIs REST**
   - Endpoints siguiendo convenciones RESTful
   - Validación de datos de entrada
   - Manejo de errores consistente
   - Respuestas JSON estructuradas

2. **Acceso a Datos (Repositories)**
   - Interacción con PostgreSQL
   - Queries optimizadas
   - Manejo de transacciones
   - Prevención de SQL injection

3. **Lógica de Negocio (Services)**
   - Reglas de negocio centralizadas
   - Validaciones complejas
   - Orquestación de operaciones
   - Cálculos y transformaciones

4. **Testing**
   - Unit tests para services
   - Integration tests para repositories
   - API tests para handlers
   - Mocks cuando sea apropiado

## 📚 Contexto del Proyecto

### Stack Tecnológico

**Backend**:
- Go 1.21+
- Fiber v2 (framework web)
- pgx v5 (driver PostgreSQL)
- golang-migrate (migraciones)

**Base de Datos**:
- PostgreSQL 15
- Schema ya definido (ver ADR-002)
- 18 tablas + 7 vistas materializadas

### Arquitectura

Seguimos **Arquitectura por Capas** (Layered Architecture):

```
┌─────────────────────────────────┐
│   HTTP Handlers (Presentación)  │  ← Recibe requests, retorna responses
├─────────────────────────────────┤
│   Services (Lógica de Negocio)  │  ← Reglas, validaciones, orquestación
├─────────────────────────────────┤
│   Repositories (Datos)           │  ← Queries SQL, transacciones
├─────────────────────────────────┤
│   Models (Dominio)               │  ← Structs que mapean a tablas
└─────────────────────────────────┘
```

**Flujo de un Request**:
```
Request → Handler → Service → Repository → Database
                        ↓
Response ← Handler ← Service ← Repository ← Data
```

### Estructura de Directorios

```
backend/
├── cmd/
│   └── api/
│       └── main.go              # Entry point
├── internal/
│   ├── models/                  # Structs del dominio
│   │   ├── student.go
│   │   ├── course.go
│   │   └── ...
│   ├── repositories/            # Acceso a datos
│   │   ├── student_repository.go
│   │   ├── course_repository.go
│   │   └── ...
│   ├── services/                # Lógica de negocio
│   │   ├── student_service.go
│   │   ├── course_service.go
│   │   └── ...
│   ├── handlers/                # HTTP handlers
│   │   ├── student_handler.go
│   │   ├── course_handler.go
│   │   └── ...
│   ├── middleware/              # Middlewares personalizados
│   │   ├── auth.go
│   │   ├── logger.go
│   │   └── ...
│   ├── database/                # Conexión y configuración DB
│   │   └── postgres.go
│   └── shared/                  # Utilidades compartidas
│       ├── errors.go
│       ├── response.go
│       └── validator.go
├── pkg/                         # Código reutilizable público
│   └── ...
├── migrations/                  # Migraciones SQL (ya creadas)
├── go.mod
└── go.sum
```

## 🔧 Metodología de Trabajo

### Proceso de Implementación de un Módulo

Ejemplo: Implementar módulo de **Estudiantes**

#### 1. Model (Dominio)

Define la estructura de datos que mapea a la tabla.

```go
// internal/models/student.go
package models

import (
    "time"
    "github.com/google/uuid"
)

type StudentStatus string

const (
    StudentStatusActive    StudentStatus = "active"
    StudentStatusGraduated StudentStatus = "graduated"
    StudentStatusWithdrawn StudentStatus = "withdrawn"
    StudentStatusSuspended StudentStatus = "suspended"
)

type Student struct {
    ID              uuid.UUID     `json:"id" db:"id"`
    FullName        string        `json:"full_name" db:"full_name"`
    DocumentID      *string       `json:"document_id,omitempty" db:"document_id"`
    BirthDate       time.Time     `json:"birth_date" db:"birth_date"`
    ProfilePhotoURL *string       `json:"profile_photo_url,omitempty" db:"profile_photo_url"`
    
    // Procedencia
    CityOriginID    *uuid.UUID    `json:"city_origin_id,omitempty" db:"city_origin_id"`
    CountryOriginID uuid.UUID     `json:"country_origin_id" db:"country_origin_id"`
    
    // Contacto (PostgreSQL arrays)
    Emails          []string      `json:"emails" db:"emails"`
    Phones          []string      `json:"phones" db:"phones"`
    
    // Laboral
    CompanyID       *uuid.UUID    `json:"company_id,omitempty" db:"company_id"`
    
    // Estado académico
    Status          StudentStatus `json:"status" db:"status"`
    Cohort          string        `json:"cohort" db:"cohort"`
    EnrollmentDate  time.Time     `json:"enrollment_date" db:"enrollment_date"`
    GraduationDate  *time.Time    `json:"graduation_date,omitempty" db:"graduation_date"`
    
    // Auditoría
    CreatedAt       time.Time     `json:"created_at" db:"created_at"`
    CreatedBy       *uuid.UUID    `json:"created_by,omitempty" db:"created_by"`
    UpdatedAt       time.Time     `json:"updated_at" db:"updated_at"`
    UpdatedBy       *uuid.UUID    `json:"updated_by,omitempty" db:"updated_by"`
    DeletedAt       *time.Time    `json:"deleted_at,omitempty" db:"deleted_at"`
    DeletedBy       *uuid.UUID    `json:"deleted_by,omitempty" db:"deleted_by"`
}

// CreateStudentRequest - DTO para crear estudiante
type CreateStudentRequest struct {
    FullName        string    `json:"full_name" validate:"required,min=3,max=255"`
    DocumentID      *string   `json:"document_id" validate:"omitempty,max=50"`
    BirthDate       string    `json:"birth_date" validate:"required"` // "1995-01-01"
    ProfilePhotoURL *string   `json:"profile_photo_url" validate:"omitempty,url"`
    CityOriginID    *string   `json:"city_origin_id" validate:"omitempty,uuid"`
    CountryOriginID string    `json:"country_origin_id" validate:"required,uuid"`
    Emails          []string  `json:"emails" validate:"required,min=1,dive,email"`
    Phones          []string  `json:"phones" validate:"omitempty,dive,max=50"`
    CompanyID       *string   `json:"company_id" validate:"omitempty,uuid"`
    Status          string    `json:"status" validate:"required,oneof=active graduated withdrawn suspended"`
    Cohort          string    `json:"cohort" validate:"required"`
    EnrollmentDate  string    `json:"enrollment_date" validate:"required"`
}

// UpdateStudentRequest - DTO para actualizar
type UpdateStudentRequest struct {
    FullName        *string   `json:"full_name" validate:"omitempty,min=3,max=255"`
    ProfilePhotoURL *string   `json:"profile_photo_url" validate:"omitempty,url"`
    Emails          []string  `json:"emails" validate:"omitempty,min=1,dive,email"`
    Phones          []string  `json:"phones" validate:"omitempty"`
    CompanyID       *string   `json:"company_id" validate:"omitempty,uuid"`
    Status          *string   `json:"status" validate:"omitempty,oneof=active graduated withdrawn suspended"`
}
```

#### 2. Repository (Acceso a Datos)

Maneja interacción con la base de datos.

```go
// internal/repositories/student_repository.go
package repositories

import (
    "context"
    "fmt"
    
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    
    "coordinador/internal/models"
)

type StudentRepository interface {
    Create(ctx context.Context, student *models.Student) error
    GetByID(ctx context.Context, id uuid.UUID) (*models.Student, error)
    List(ctx context.Context, filters StudentFilters) ([]*models.Student, error)
    Update(ctx context.Context, student *models.Student) error
    Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
    Count(ctx context.Context, filters StudentFilters) (int, error)
}

type studentRepository struct {
    db *pgxpool.Pool
}

func NewStudentRepository(db *pgxpool.Pool) StudentRepository {
    return &studentRepository{db: db}
}

type StudentFilters struct {
    Status          *string
    Cohort          *string
    CountryOriginID *uuid.UUID
    Search          *string // Búsqueda por nombre
    Limit           int
    Offset          int
}

func (r *studentRepository) Create(ctx context.Context, student *models.Student) error {
    query := `
        INSERT INTO students (
            id, full_name, document_id, birth_date, profile_photo_url,
            city_origin_id, country_origin_id, emails, phones,
            company_id, status, cohort, enrollment_date,
            created_by
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
        student.Status,
        student.Cohort,
        student.EnrollmentDate,
        student.CreatedBy,
    ).Scan(&student.CreatedAt, &student.UpdatedAt)
    
    if err != nil {
        return fmt.Errorf("error creating student: %w", err)
    }
    
    return nil
}

func (r *studentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Student, error) {
    query := `
        SELECT 
            id, full_name, document_id, birth_date, profile_photo_url,
            city_origin_id, country_origin_id, emails, phones,
            company_id, status, cohort, enrollment_date, graduation_date,
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
        return nil, fmt.Errorf("error getting student: %w", err)
    }
    
    return student, nil
}

func (r *studentRepository) List(ctx context.Context, filters StudentFilters) ([]*models.Student, error) {
    query := `
        SELECT 
            id, full_name, document_id, birth_date, profile_photo_url,
            city_origin_id, country_origin_id, emails, phones,
            company_id, status, cohort, enrollment_date, graduation_date,
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
    }
    
    rows, err := r.db.Query(ctx, query, args...)
    if err != nil {
        return nil, fmt.Errorf("error listing students: %w", err)
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
            return nil, fmt.Errorf("error scanning student: %w", err)
        }
        students = append(students, student)
    }
    
    return students, nil
}

func (r *studentRepository) Update(ctx context.Context, student *models.Student) error {
    query := `
        UPDATE students
        SET 
            full_name = $2,
            profile_photo_url = $3,
            emails = $4,
            phones = $5,
            company_id = $6,
            status = $7,
            updated_by = $8
        WHERE id = $1 AND deleted_at IS NULL
        RETURNING updated_at
    `
    
    err := r.db.QueryRow(ctx, query,
        student.ID,
        student.FullName,
        student.ProfilePhotoURL,
        student.Emails,
        student.Phones,
        student.CompanyID,
        student.Status,
        student.UpdatedBy,
    ).Scan(&student.UpdatedAt)
    
    if err != nil {
        if err == pgx.ErrNoRows {
            return fmt.Errorf("student not found")
        }
        return fmt.Errorf("error updating student: %w", err)
    }
    
    return nil
}

func (r *studentRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
    query := `
        UPDATE students
        SET deleted_at = NOW(), deleted_by = $2
        WHERE id = $1 AND deleted_at IS NULL
    `
    
    result, err := r.db.Exec(ctx, query, id, deletedBy)
    if err != nil {
        return fmt.Errorf("error deleting student: %w", err)
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
    }
    
    var count int
    err := r.db.QueryRow(ctx, query, args...).Scan(&count)
    if err != nil {
        return 0, fmt.Errorf("error counting students: %w", err)
    }
    
    return count, nil
}
```

#### 3. Service (Lógica de Negocio)

Orquesta operaciones, valida reglas de negocio.

```go
// internal/services/student_service.go
package services

import (
    "context"
    "fmt"
    "time"
    
    "github.com/google/uuid"
    
    "coordinador/internal/models"
    "coordinador/internal/repositories"
)

type StudentService interface {
    CreateStudent(ctx context.Context, req *models.CreateStudentRequest, createdBy uuid.UUID) (*models.Student, error)
    GetStudent(ctx context.Context, id uuid.UUID) (*models.Student, error)
    ListStudents(ctx context.Context, filters repositories.StudentFilters) ([]*models.Student, int, error)
    UpdateStudent(ctx context.Context, id uuid.UUID, req *models.UpdateStudentRequest, updatedBy uuid.UUID) (*models.Student, error)
    DeleteStudent(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
}

type studentService struct {
    studentRepo repositories.StudentRepository
}

func NewStudentService(studentRepo repositories.StudentRepository) StudentService {
    return &studentService{
        studentRepo: studentRepo,
    }
}

func (s *studentService) CreateStudent(ctx context.Context, req *models.CreateStudentRequest, createdBy uuid.UUID) (*models.Student, error) {
    // Validar edad mínima (18 años)
    birthDate, err := time.Parse("2006-01-02", req.BirthDate)
    if err != nil {
        return nil, fmt.Errorf("invalid birth date format")
    }
    
    age := time.Now().Year() - birthDate.Year()
    if age < 18 {
        return nil, fmt.Errorf("student must be at least 18 years old")
    }
    
    // Parsear UUIDs
    countryOriginID, err := uuid.Parse(req.CountryOriginID)
    if err != nil {
        return nil, fmt.Errorf("invalid country origin ID")
    }
    
    var cityOriginID *uuid.UUID
    if req.CityOriginID != nil {
        parsed, err := uuid.Parse(*req.CityOriginID)
        if err != nil {
            return nil, fmt.Errorf("invalid city origin ID")
        }
        cityOriginID = &parsed
    }
    
    var companyID *uuid.UUID
    if req.CompanyID != nil {
        parsed, err := uuid.Parse(*req.CompanyID)
        if err != nil {
            return nil, fmt.Errorf("invalid company ID")
        }
        companyID = &parsed
    }
    
    enrollmentDate, err := time.Parse("2006-01-02", req.EnrollmentDate)
    if err != nil {
        return nil, fmt.Errorf("invalid enrollment date format")
    }
    
    // Crear modelo
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
        Status:          models.StudentStatus(req.Status),
        Cohort:          req.Cohort,
        EnrollmentDate:  enrollmentDate,
        CreatedBy:       &createdBy,
    }
    
    // Guardar
    err = s.studentRepo.Create(ctx, student)
    if err != nil {
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

func (s *studentService) UpdateStudent(ctx context.Context, id uuid.UUID, req *models.UpdateStudentRequest, updatedBy uuid.UUID) (*models.Student, error) {
    // Obtener estudiante actual
    student, err := s.studentRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Actualizar campos
    if req.FullName != nil {
        student.FullName = *req.FullName
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
            return nil, fmt.Errorf("invalid company ID")
        }
        student.CompanyID = &parsed
    }
    if req.Status != nil {
        student.Status = models.StudentStatus(*req.Status)
    }
    
    student.UpdatedBy = &updatedBy
    
    // Guardar cambios
    err = s.studentRepo.Update(ctx, student)
    if err != nil {
        return nil, fmt.Errorf("failed to update student: %w", err)
    }
    
    return student, nil
}

func (s *studentService) DeleteStudent(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
    return s.studentRepo.Delete(ctx, id, deletedBy)
}
```

#### 4. Handler (HTTP Controller)

Maneja requests HTTP, validación, responses.

```go
// internal/handlers/student_handler.go
package handlers

import (
    "strconv"
    
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    
    "coordinador/internal/models"
    "coordinador/internal/repositories"
    "coordinador/internal/services"
    "coordinador/internal/shared"
)

type StudentHandler struct {
    studentService services.StudentService
}

func NewStudentHandler(studentService services.StudentService) *StudentHandler {
    return &StudentHandler{
        studentService: studentService,
    }
}

// POST /api/v1/students
func (h *StudentHandler) CreateStudent(c *fiber.Ctx) error {
    var req models.CreateStudentRequest
    
    if err := c.BodyParser(&req); err != nil {
        return shared.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
    }
    
    // TODO: Obtener usuario actual de contexto (después de implementar auth)
    createdBy := uuid.New() // Temporal
    
    student, err := h.studentService.CreateStudent(c.Context(), &req, createdBy)
    if err != nil {
        return shared.ErrorResponse(c, fiber.StatusBadRequest, "Failed to create student", err)
    }
    
    return shared.SuccessResponse(c, fiber.StatusCreated, "Student created successfully", student)
}

// GET /api/v1/students/:id
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

// GET /api/v1/students
func (h *StudentHandler) ListStudents(c *fiber.Ctx) error {
    filters := repositories.StudentFilters{}
    
    // Query params
    if status := c.Query("status"); status != "" {
        filters.Status = &status
    }
    if cohort := c.Query("cohort"); cohort != "" {
        filters.Cohort = &cohort
    }
    if search := c.Query("search"); search != "" {
        filters.Search = &search
    }
    if countryID := c.Query("country_id"); countryID != "" {
        parsed, err := uuid.Parse(countryID)
        if err == nil {
            filters.CountryOriginID = &parsed
        }
    }
    
    // Paginación
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

// PUT /api/v1/students/:id
func (h *StudentHandler) UpdateStudent(c *fiber.Ctx) error {
    id, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return shared.ErrorResponse(c, fiber.StatusBadRequest, "Invalid student ID", err)
    }
    
    var req models.UpdateStudentRequest
    if err := c.BodyParser(&req); err != nil {
        return shared.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
    }
    
    // TODO: Obtener usuario actual de contexto
    updatedBy := uuid.New() // Temporal
    
    student, err := h.studentService.UpdateStudent(c.Context(), id, &req, updatedBy)
    if err != nil {
        return shared.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update student", err)
    }
    
    return shared.SuccessResponse(c, fiber.StatusOK, "Student updated successfully", student)
}

// DELETE /api/v1/students/:id
func (h *StudentHandler) DeleteStudent(c *fiber.Ctx) error {
    id, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return shared.ErrorResponse(c, fiber.StatusBadRequest, "Invalid student ID", err)
    }
    
    // TODO: Obtener usuario actual de contexto
    deletedBy := uuid.New() // Temporal
    
    err = h.studentService.DeleteStudent(c.Context(), id, deletedBy)
    if err != nil {
        return shared.ErrorResponse(c, fiber.StatusNotFound, "Failed to delete student", err)
    }
    
    return shared.SuccessResponse(c, fiber.StatusOK, "Student deleted successfully", nil)
}

// RegisterRoutes - Registrar rutas
func (h *StudentHandler) RegisterRoutes(router fiber.Router) {
    students := router.Group("/students")
    
    students.Post("/", h.CreateStudent)
    students.Get("/", h.ListStudents)
    students.Get("/:id", h.GetStudent)
    students.Put("/:id", h.UpdateStudent)
    students.Delete("/:id", h.DeleteStudent)
}
```

#### 5. Shared Utilities

```go
// internal/shared/response.go
package shared

import "github.com/gofiber/fiber/v2"

type APIResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   *string     `json:"error,omitempty"`
}

type PaginatedData struct {
    Items  interface{} `json:"items"`
    Total  int         `json:"total"`
    Limit  int         `json:"limit"`
    Offset int         `json:"offset"`
}

func SuccessResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
    return c.Status(status).JSON(APIResponse{
        Success: true,
        Message: message,
        Data:    data,
    })
}

func ErrorResponse(c *fiber.Ctx, status int, message string, err error) error {
    errMsg := err.Error()
    return c.Status(status).JSON(APIResponse{
        Success: false,
        Message: message,
        Error:   &errMsg,
    })
}

func PaginatedResponse(c *fiber.Ctx, status int, message string, items interface{}, total, limit, offset int) error {
    return c.Status(status).JSON(APIResponse{
        Success: true,
        Message: message,
        Data: PaginatedData{
            Items:  items,
            Total:  total,
            Limit:  limit,
            Offset: offset,
        },
    })
}
```

## 🧪 Testing

### Unit Test (Service)

```go
// internal/services/student_service_test.go
package services_test

import (
    "context"
    "testing"
    
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    
    "coordinador/internal/models"
    "coordinador/internal/services"
    "coordinador/internal/repositories/mocks"
)

func TestCreateStudent_Success(t *testing.T) {
    mockRepo := new(mocks.StudentRepository)
    service := services.NewStudentService(mockRepo)
    
    req := &models.CreateStudentRequest{
        FullName:        "Juan Test",
        BirthDate:       "1995-01-01",
        CountryOriginID: uuid.New().String(),
        Emails:          []string{"juan@test.com"},
        Status:          "active",
        Cohort:          "2024-1",
        EnrollmentDate:  "2024-01-15",
    }
    
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
    
    student, err := service.CreateStudent(context.Background(), req, uuid.New())
    
    assert.NoError(t, err)
    assert.NotNil(t, student)
    assert.Equal(t, "Juan Test", student.FullName)
    mockRepo.AssertExpectations(t)
}

func TestCreateStudent_InvalidAge(t *testing.T) {
    mockRepo := new(mocks.StudentRepository)
    service := services.NewStudentService(mockRepo)
    
    req := &models.CreateStudentRequest{
        FullName:        "Minor Test",
        BirthDate:       "2010-01-01", // Menor de 18
        CountryOriginID: uuid.New().String(),
        Emails:          []string{"minor@test.com"},
        Status:          "active",
        Cohort:          "2024-1",
        EnrollmentDate:  "2024-01-15",
    }
    
    student, err := service.CreateStudent(context.Background(), req, uuid.New())
    
    assert.Error(t, err)
    assert.Nil(t, student)
    assert.Contains(t, err.Error(), "18 years old")
}
```

### Integration Test (Repository)

```go
// internal/repositories/student_repository_test.go
package repositories_test

import (
    "context"
    "testing"
    "time"
    
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    
    "coordinador/internal/database"
    "coordinador/internal/models"
    "coordinador/internal/repositories"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
    // Setup test database connection
    // Similar al setup de producción pero con DB de test
}

func TestStudentRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    repo := repositories.NewStudentRepository(db)
    
    student := &models.Student{
        ID:              uuid.New(),
        FullName:        "Test Student",
        BirthDate:       time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC),
        CountryOriginID: uuid.New(),
        Emails:          []string{"test@example.com"},
        Status:          models.StudentStatusActive,
        Cohort:          "2024-1",
        EnrollmentDate:  time.Now(),
    }
    
    err := repo.Create(context.Background(), student)
    require.NoError(t, err)
    assert.NotZero(t, student.CreatedAt)
}
```

## 📋 Checklist de Implementación

Para cada módulo nuevo:

- [ ] Model creado con tags JSON y DB
- [ ] DTOs de request/response definidos
- [ ] Repository interface definida
- [ ] Repository implementación con queries SQL
- [ ] Service interface definida
- [ ] Service implementación con lógica de negocio
- [ ] Handler creado con endpoints REST
- [ ] Rutas registradas en main.go
- [ ] Unit tests para service (>80% coverage)
- [ ] Integration tests para repository
- [ ] Documentación de API actualizada
- [ ] Validaciones implementadas
- [ ] Error handling consistente

## 🎯 Convenciones y Best Practices

### Naming

```go
// Interfaces: Singular + "Repository" o "Service"
type StudentRepository interface {}
type StudentService interface {}

// Implementaciones: lowercase + "Repository" o "Service"
type studentRepository struct {}
type studentService struct {}

// Handlers: Singular + "Handler"
type StudentHandler struct {}

// Models: Singular
type Student struct {}

// DTOs: Action + Model + "Request/Response"
type CreateStudentRequest struct {}
type UpdateStudentRequest struct {}
```

### Error Handling

```go
// Siempre wrap errors con contexto
return fmt.Errorf("failed to create student: %w", err)

// Usar errores específicos cuando sea apropiado
var ErrStudentNotFound = errors.New("student not found")
var ErrInvalidAge = errors.New("student must be at least 18 years old")

// En handlers, convertir a status HTTP apropiado
if errors.Is(err, ErrStudentNotFound) {
    return shared.ErrorResponse(c, fiber.StatusNotFound, "Student not found", err)
}
```

### SQL Queries

```go
// ✅ BUENO: Usar placeholders ($1, $2, etc.)
query := "SELECT * FROM students WHERE id = $1"
db.QueryRow(ctx, query, id)

// ❌ MALO: Concatenación de strings (SQL injection!)
query := fmt.Sprintf("SELECT * FROM students WHERE id = '%s'", id)

// ✅ BUENO: Queries legibles (multi-línea)
query := `
    SELECT 
        id, full_name, email
    FROM students
    WHERE deleted_at IS NULL
        AND status = $1
    ORDER BY created_at DESC
`

// ✅ BUENO: Siempre filtrar soft-deleted
WHERE deleted_at IS NULL
```

### Validación

```go
// En DTOs usar tags de validación
type CreateStudentRequest struct {
    FullName string `json:"full_name" validate:"required,min=3,max=255"`
    Email    string `json:"email" validate:"required,email"`
    Age      int    `json:"age" validate:"required,min=18"`
}

// En services validar reglas de negocio complejas
if age < 18 {
    return nil, ErrInvalidAge
}
```

### Contexto

```go
// Siempre pasar context como primer parámetro
func (r *repository) Create(ctx context.Context, student *Student) error

// Usar context para timeouts
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Usar context para cancelación
if ctx.Err() == context.Canceled {
    return ctx.Err()
}
```

## 🔄 Interacción con Otros Agentes

### Agente DBA
- **Le consulto**: Optimizaciones de queries, índices necesarios
- **Le proporciono**: Patrones de acceso reales, queries lentos
- **Coordino**: Cuando necesito cambios en el schema

### Agente Arquitecto
- **Le consulto**: Decisiones de estructura, nuevos patrones
- **Le proporciono**: Feedback sobre arquitectura actual
- **Coordino**: Cuando encuentro limitaciones arquitectónicas

### Agente Svelte
- **Le proporciono**: Documentación de API, contratos JSON
- **Recibo de él**: Requerimientos de endpoints faltantes
- **Coordino**: Formato de responses, validaciones

## 📚 Recursos

- [Go Documentation](https://go.dev/doc/)
- [Fiber Framework](https://docs.gofiber.io/)
- [pgx PostgreSQL Driver](https://github.com/jackc/pgx)
- [Go Testing](https://go.dev/doc/tutorial/add-a-test)
- [Effective Go](https://go.dev/doc/effective_go)

## 🎓 Principios de Go

- **Simplicidad**: Keep it simple, stupid (KISS)
- **Composición sobre herencia**: Use interfaces
- **Errores son valores**: Handle them explicitly
- **Concurrencia**: Goroutines y channels cuando apropiado
- **Testing**: Test coverage >80%

---

**Última actualización**: 2026-02-13
