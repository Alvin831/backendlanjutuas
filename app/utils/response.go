package utils

// Response adalah struktur standar output JSON untuk Swagger
// @Description Standard API response format
type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

// APIResponse adalah struktur standar output JSON
type APIResponse struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

// Meta contains response metadata
// @Description Response metadata
type Meta struct {
	Message string `json:"message" example:"Success"`
	Code    int    `json:"code" example:"200"`
	Status  string `json:"status" example:"success"`
}

// SuccessResponse untuk format sukses (Code 200-299)
func SuccessResponse(message string, code int, data interface{}) APIResponse {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  "success",
	}

	return APIResponse{
		Meta: meta,
		Data: data,
	}
}

// ErrorResponse untuk format gagal (Code 400-500)
func ErrorResponse(message string, code int, data interface{}) APIResponse {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  "error",
	}

	return APIResponse{
		Meta: meta,
		Data: data,
	}
}	