package services

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"net/mail"
	"strings"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"

	"github.com/dcorreal/coordinador/internal/models"
	"github.com/dcorreal/coordinador/internal/repositories"
)

// StudentImportService handles bulk student imports from files.
type StudentImportService interface {
	ImportFromFile(ctx context.Context, fileData []byte, format string, createdBy *uuid.UUID) (*models.ImportResult, error)
}

type studentImportService struct {
	studentService StudentService
	studentRepo    repositories.StudentRepository
}

// NewStudentImportService creates a new StudentImportService.
func NewStudentImportService(studentService StudentService, studentRepo repositories.StudentRepository) StudentImportService {
	return &studentImportService{
		studentService: studentService,
		studentRepo:    studentRepo,
	}
}

// expectedHeaders defines the column order for import files.
var expectedHeaders = []string{
	"first_names", "last_names", "document_id", "birth_date", "gender",
	"email", "phone",
	"nationality_country_id", "residence_country_id", "residence_city_id",
	"company_id", "job_title_category_id", "profession_id",
	"student_code", "status", "cohort", "enrollment_date",
}

func (s *studentImportService) ImportFromFile(ctx context.Context, fileData []byte, format string, createdBy *uuid.UUID) (*models.ImportResult, error) {
	var rows [][]string
	var err error

	switch format {
	case "csv":
		rows, err = parseCSV(fileData)
	case "xlsx":
		rows, err = parseXLSX(fileData)
	default:
		return nil, fmt.Errorf("unsupported format: %s, expected csv or xlsx", format)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("file must have a header row and at least one data row")
	}

	headerMap, err := mapHeaders(rows[0])
	if err != nil {
		return nil, err
	}

	dataRows := rows[1:]
	result := &models.ImportResult{
		TotalRows: len(dataRows),
		Errors:    []models.ImportRowError{},
	}

	// Collect all document_ids and emails for batch duplicate check
	var allDocIDs []string
	var allEmails []string
	for _, row := range dataRows {
		if docID := getField(row, headerMap, "document_id"); docID != "" {
			allDocIDs = append(allDocIDs, docID)
		}
		if email := getField(row, headerMap, "email"); email != "" {
			allEmails = append(allEmails, email)
		}
	}

	existingDocs, err := s.studentRepo.ExistingDocumentIDs(ctx, allDocIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing documents: %w", err)
	}
	existingEmails, err := s.studentRepo.ExistingEmails(ctx, allEmails)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing emails: %w", err)
	}

	// Track duplicates within the file itself
	seenDocs := make(map[string]int)
	seenEmails := make(map[string]int)

	for i, row := range dataRows {
		rowNum := i + 2 // 1-based, skip header

		rowErrors := s.validateAndImportRow(ctx, row, headerMap, rowNum, createdBy, existingDocs, existingEmails, seenDocs, seenEmails)
		if len(rowErrors) > 0 {
			result.Errors = append(result.Errors, rowErrors...)
		} else {
			result.Created++
		}
	}

	return result, nil
}

func (s *studentImportService) validateAndImportRow(
	ctx context.Context,
	row []string,
	headerMap map[string]int,
	rowNum int,
	createdBy *uuid.UUID,
	existingDocs, existingEmails map[string]bool,
	seenDocs map[string]int,
	seenEmails map[string]int,
) []models.ImportRowError {
	var errors []models.ImportRowError

	addError := func(field, value, message string) {
		errors = append(errors, models.ImportRowError{
			Row:     rowNum,
			Field:   field,
			Value:   value,
			Message: message,
		})
	}

	// Extract fields
	firstNames := strings.TrimSpace(getField(row, headerMap, "first_names"))
	lastNames := strings.TrimSpace(getField(row, headerMap, "last_names"))
	documentID := strings.TrimSpace(getField(row, headerMap, "document_id"))
	birthDate := strings.TrimSpace(getField(row, headerMap, "birth_date"))
	gender := strings.TrimSpace(getField(row, headerMap, "gender"))
	email := strings.TrimSpace(getField(row, headerMap, "email"))
	phone := strings.TrimSpace(getField(row, headerMap, "phone"))
	nationalityCountryID := strings.TrimSpace(getField(row, headerMap, "nationality_country_id"))
	residenceCountryID := strings.TrimSpace(getField(row, headerMap, "residence_country_id"))
	residenceCityID := strings.TrimSpace(getField(row, headerMap, "residence_city_id"))
	companyID := strings.TrimSpace(getField(row, headerMap, "company_id"))
	jobTitleCategoryID := strings.TrimSpace(getField(row, headerMap, "job_title_category_id"))
	professionID := strings.TrimSpace(getField(row, headerMap, "profession_id"))
	studentCode := strings.TrimSpace(getField(row, headerMap, "student_code"))
	status := strings.TrimSpace(getField(row, headerMap, "status"))
	cohort := strings.TrimSpace(getField(row, headerMap, "cohort"))
	enrollmentDate := strings.TrimSpace(getField(row, headerMap, "enrollment_date"))

	// Validate required fields
	if firstNames == "" {
		addError("first_names", "", "required field is empty")
	}
	if lastNames == "" {
		addError("last_names", "", "required field is empty")
	}
	if birthDate == "" {
		addError("birth_date", "", "required field is empty")
	}
	if email == "" {
		addError("email", "", "required field is empty")
	} else if _, err := mail.ParseAddress(email); err != nil {
		addError("email", email, "invalid email format")
	}
	if nationalityCountryID == "" {
		addError("nationality_country_id", "", "required field is empty")
	} else if _, err := uuid.Parse(nationalityCountryID); err != nil {
		addError("nationality_country_id", nationalityCountryID, "invalid UUID format")
	}
	if residenceCountryID == "" {
		addError("residence_country_id", "", "required field is empty")
	} else if _, err := uuid.Parse(residenceCountryID); err != nil {
		addError("residence_country_id", residenceCountryID, "invalid UUID format")
	}
	if status == "" {
		addError("status", "", "required field is empty")
	}
	if cohort == "" {
		addError("cohort", "", "required field is empty")
	}
	if enrollmentDate == "" {
		addError("enrollment_date", "", "required field is empty")
	}

	// Validate gender if provided
	if gender != "" {
		g := strings.ToUpper(gender)
		if g != "M" && g != "F" {
			addError("gender", gender, "must be M or F")
		} else {
			gender = g
		}
	}

	// Check duplicates against DB
	if documentID != "" {
		if existingDocs[documentID] {
			addError("document_id", documentID, "duplicate: student with this document already exists")
		} else if prevRow, ok := seenDocs[documentID]; ok {
			addError("document_id", documentID, fmt.Sprintf("duplicate: same document_id as row %d in this file", prevRow))
		}
	}
	if email != "" {
		if existingEmails[email] {
			addError("email", email, "duplicate: student with this email already exists")
		} else if prevRow, ok := seenEmails[email]; ok {
			addError("email", email, fmt.Sprintf("duplicate: same email as row %d in this file", prevRow))
		}
	}

	// If there are validation errors, don't attempt to insert
	if len(errors) > 0 {
		return errors
	}

	// Build CreateStudentRequest
	req := &models.CreateStudentRequest{
		FirstNames:           firstNames,
		LastNames:            lastNames,
		BirthDate:            birthDate,
		NationalityCountryID: nationalityCountryID,
		ResidenceCountryID:   residenceCountryID,
		Emails:               []string{email},
		Status:               status,
		Cohort:               cohort,
		EnrollmentDate:       enrollmentDate,
	}

	if documentID != "" {
		req.DocumentID = &documentID
	}
	if gender != "" {
		req.Gender = &gender
	}
	if phone != "" {
		req.Phones = []string{phone}
	}
	if residenceCityID != "" {
		req.ResidenceCityID = &residenceCityID
	}
	if companyID != "" {
		req.CompanyID = &companyID
	}
	if jobTitleCategoryID != "" {
		req.JobTitleCategoryID = &jobTitleCategoryID
	}
	if professionID != "" {
		req.ProfessionID = &professionID
	}
	if studentCode != "" {
		req.StudentCode = &studentCode
	}

	// Use existing CreateStudent which does full validation + insert
	_, err := s.studentService.CreateStudent(ctx, req, createdBy)
	if err != nil {
		addError("_row", "", err.Error())
		return errors
	}

	// Track as seen for intra-file duplicate detection
	if documentID != "" {
		seenDocs[documentID] = rowNum
	}
	seenEmails[email] = rowNum

	return nil
}

func parseCSV(data []byte) ([][]string, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true
	return reader.ReadAll()
}

func parseXLSX(data []byte) ([][]string, error) {
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("no sheets found in xlsx file")
	}

	return f.GetRows(sheetName)
}

func mapHeaders(headerRow []string) (map[string]int, error) {
	m := make(map[string]int)
	for i, h := range headerRow {
		normalized := strings.ToLower(strings.TrimSpace(h))
		m[normalized] = i
	}

	// Check that at least the required headers are present
	required := []string{"first_names", "last_names", "birth_date", "email",
		"nationality_country_id", "residence_country_id", "status", "cohort", "enrollment_date"}
	var missing []string
	for _, r := range required {
		if _, ok := m[r]; !ok {
			missing = append(missing, r)
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required columns: %s", strings.Join(missing, ", "))
	}

	return m, nil
}

func getField(row []string, headerMap map[string]int, field string) string {
	idx, ok := headerMap[field]
	if !ok || idx >= len(row) {
		return ""
	}
	return row[idx]
}
