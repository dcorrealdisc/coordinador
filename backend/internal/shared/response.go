package shared

import "github.com/gofiber/fiber/v2"

// APIResponse is the standard JSON envelope for all API responses.
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *string     `json:"error,omitempty"`
}

// PaginatedData wraps a list with pagination metadata.
type PaginatedData struct {
	Items  interface{} `json:"items"`
	Total  int         `json:"total"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
}

// SuccessResponse sends a successful JSON response.
func SuccessResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error JSON response.
func ErrorResponse(c *fiber.Ctx, status int, message string, err error) error {
	errMsg := err.Error()
	return c.Status(status).JSON(APIResponse{
		Success: false,
		Message: message,
		Error:   &errMsg,
	})
}

// PaginatedResponse sends a paginated JSON response.
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
