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
	studentService  StudentService
	studentRepo     repositories.StudentRepository
	catalogRepo     repositories.CatalogRepository
	catalogResolver *CatalogResolver
}

// NewStudentImportService creates a new StudentImportService.
func NewStudentImportService(
	studentService StudentService,
	studentRepo repositories.StudentRepository,
	catalogRepo repositories.CatalogRepository,
) StudentImportService {
	return &studentImportService{
		studentService:  studentService,
		studentRepo:     studentRepo,
		catalogRepo:     catalogRepo,
		catalogResolver: NewCatalogResolver(catalogRepo),
	}
}

// statusMap translates Spanish status values to English.
var statusMap = map[string]string{
	"activo":     "active",
	"graduado":   "graduated",
	"retirado":   "withdrawn",
	"suspendido": "suspended",
}

func normalizeStatus(s string) string {
	lower := strings.ToLower(strings.TrimSpace(s))
	if mapped, ok := statusMap[lower]; ok {
		return mapped
	}
	return s
}

// isUUID checks if a string looks like a valid UUID.
func isUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
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
	nationalityRaw := strings.TrimSpace(getField(row, headerMap, "nationality_country_id"))
	residenceRaw := strings.TrimSpace(getField(row, headerMap, "residence_country_id"))
	residenceCityRaw := strings.TrimSpace(getField(row, headerMap, "residence_city_id"))
	companyRaw := strings.TrimSpace(getField(row, headerMap, "company_id"))
	jobTitleRaw := strings.TrimSpace(getField(row, headerMap, "job_title_category_id"))
	professionRaw := strings.TrimSpace(getField(row, headerMap, "profession_id"))
	studentCode := strings.TrimSpace(getField(row, headerMap, "student_code"))
	status := normalizeStatus(getField(row, headerMap, "status"))
	cohort := strings.TrimSpace(getField(row, headerMap, "cohort"))
	enrollmentDate := strings.TrimSpace(getField(row, headerMap, "enrollment_date"))

	// University columns (optional)
	universityName := strings.TrimSpace(getField(row, headerMap, "universidad"))
	universityCityName := strings.TrimSpace(getField(row, headerMap, "universidad-ciudad"))
	universityCountryName := strings.TrimSpace(getField(row, headerMap, "universidad-pais"))

	// Validate required fields
	if firstNames == "" {
		addError("first_names", "", "required field is empty")
	}
	if lastNames == "" {
		addError("last_names", "", "required field is empty")
	}
	if email != "" {
		if _, err := mail.ParseAddress(email); err != nil {
			addError("email", email, "invalid email format")
		}
	}
	if nationalityRaw == "" {
		addError("nationality_country_id", "", "required field is empty")
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

	// If there are validation errors, don't attempt to resolve catalogs or insert
	if len(errors) > 0 {
		return errors
	}

	// --- Resolve names to UUIDs via CatalogResolver ---

	// Nationality country (required)
	var nationalityCountryUUID string
	if isUUID(nationalityRaw) {
		nationalityCountryUUID = nationalityRaw
	} else {
		id, err := s.catalogResolver.ResolveCountry(ctx, nationalityRaw)
		if err != nil {
			addError("nationality_country_id", nationalityRaw, err.Error())
			return errors
		}
		nationalityCountryUUID = id.String()
	}

	// Residence country: use nationality if not provided
	var residenceCountryUUID string
	if residenceRaw == "" {
		residenceCountryUUID = nationalityCountryUUID
	} else if isUUID(residenceRaw) {
		residenceCountryUUID = residenceRaw
	} else {
		id, err := s.catalogResolver.ResolveCountry(ctx, residenceRaw)
		if err != nil {
			addError("residence_country_id", residenceRaw, err.Error())
			return errors
		}
		residenceCountryUUID = id.String()
	}

	// Residence city (optional)
	var residenceCityUUID *string
	if residenceCityRaw != "" {
		if isUUID(residenceCityRaw) {
			residenceCityUUID = &residenceCityRaw
		} else {
			resCountryID, _ := uuid.Parse(residenceCountryUUID)
			id, err := s.catalogResolver.ResolveCity(ctx, residenceCityRaw, resCountryID)
			if err != nil {
				addError("residence_city_id", residenceCityRaw, err.Error())
				return errors
			}
			if id != uuid.Nil {
				s := id.String()
				residenceCityUUID = &s
			}
		}
	}

	// Profession (optional)
	var professionUUID *string
	if professionRaw != "" {
		if isUUID(professionRaw) {
			professionUUID = &professionRaw
		} else {
			id, err := s.catalogResolver.ResolveProfession(ctx, professionRaw)
			if err != nil {
				addError("profession_id", professionRaw, err.Error())
				return errors
			}
			if id != uuid.Nil {
				s := id.String()
				professionUUID = &s
			}
		}
	}

	// Job title category (optional)
	var jobTitleUUID *string
	if jobTitleRaw != "" {
		if isUUID(jobTitleRaw) {
			jobTitleUUID = &jobTitleRaw
		} else {
			id, err := s.catalogResolver.ResolveJobTitleCategory(ctx, jobTitleRaw)
			if err != nil {
				addError("job_title_category_id", jobTitleRaw, err.Error())
				return errors
			}
			if id != uuid.Nil {
				s := id.String()
				jobTitleUUID = &s
			}
		}
	}

	// Build CreateStudentRequest
	req := &models.CreateStudentRequest{
		FirstNames:           firstNames,
		LastNames:            lastNames,
		BirthDate:            birthDate,
		NationalityCountryID: nationalityCountryUUID,
		ResidenceCountryID:   residenceCountryUUID,
		ResidenceCityID:      residenceCityUUID,
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
	if email != "" {
		req.Emails = []string{email}
	}
	if phone != "" {
		req.Phones = []string{phone}
	}
	if companyRaw != "" && isUUID(companyRaw) {
		req.CompanyID = &companyRaw
	}
	if jobTitleUUID != nil {
		req.JobTitleCategoryID = jobTitleUUID
	}
	if professionUUID != nil {
		req.ProfessionID = professionUUID
	}
	if studentCode != "" {
		req.StudentCode = &studentCode
	}

	// Create the student
	student, err := s.studentService.CreateStudent(ctx, req, createdBy)
	if err != nil {
		addError("_row", "", err.Error())
		return errors
	}

	// --- Create student_university relationship ---
	if universityName != "" {
		// Resolve university country (default to nationality if not provided)
		uniCountryID, _ := uuid.Parse(nationalityCountryUUID)
		if universityCountryName != "" {
			resolved, err := s.catalogResolver.ResolveCountry(ctx, universityCountryName)
			if err == nil && resolved != uuid.Nil {
				uniCountryID = resolved
			}
		}

		// Resolve university city (optional)
		var uniCityID *uuid.UUID
		if universityCityName != "" {
			resolved, err := s.catalogResolver.ResolveCity(ctx, universityCityName, uniCountryID)
			if err == nil && resolved != uuid.Nil {
				uniCityID = &resolved
			}
		}

		// Resolve university
		uniID, err := s.catalogResolver.ResolveUniversity(ctx, universityName, uniCityID, uniCountryID)
		if err == nil && uniID != uuid.Nil {
			_ = s.catalogRepo.CreateStudentUniversity(ctx, student.ID, uniID)
		}
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

	// Required headers — residence_country_id is no longer required since we fall back to nationality
	required := []string{"first_names", "last_names",
		"nationality_country_id", "status", "cohort", "enrollment_date"}
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
