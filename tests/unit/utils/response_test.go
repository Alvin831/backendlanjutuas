package utils_test

import (
	"testing"
	"uas_backend/app/utils"

	"github.com/stretchr/testify/assert"
)

func TestSuccessResponse(t *testing.T) {
	tests := []struct {
		name    string
		message string
		code    int
		data    interface{}
	}{
		{
			name:    "Success with data",
			message: "Operation successful",
			code:    200,
			data:    map[string]string{"key": "value"},
		},
		{
			name:    "Success with nil data",
			message: "Success",
			code:    201,
			data:    nil,
		},
		{
			name:    "Success with array data",
			message: "List retrieved",
			code:    200,
			data:    []string{"item1", "item2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := utils.SuccessResponse(tt.message, tt.code, tt.data)
			
			// Check meta fields
			assert.Equal(t, tt.message, response.Meta.Message)
			assert.Equal(t, tt.code, response.Meta.Code)
			assert.Equal(t, "success", response.Meta.Status)
			
			// Check data
			assert.Equal(t, tt.data, response.Data)
		})
	}
}

func TestErrorResponse(t *testing.T) {
	tests := []struct {
		name    string
		message string
		code    int
		data    interface{}
	}{
		{
			name:    "Error with data",
			message: "Validation failed",
			code:    400,
			data:    map[string]string{"field": "required"},
		},
		{
			name:    "Error with nil data",
			message: "Not found",
			code:    404,
			data:    nil,
		},
		{
			name:    "Server error",
			message: "Internal server error",
			code:    500,
			data:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := utils.ErrorResponse(tt.message, tt.code, tt.data)
			
			// Check meta fields
			assert.Equal(t, tt.message, response.Meta.Message)
			assert.Equal(t, tt.code, response.Meta.Code)
			assert.Equal(t, "error", response.Meta.Status)
			
			// Check data
			assert.Equal(t, tt.data, response.Data)
		})
	}
}

func TestResponseStructure(t *testing.T) {
	// Test success response structure
	successResp := utils.SuccessResponse("test", 200, "data")
	assert.NotNil(t, successResp.Meta)
	assert.Equal(t, "data", successResp.Data)

	// Test error response structure
	errorResp := utils.ErrorResponse("error", 400, nil)
	assert.NotNil(t, errorResp.Meta)
	assert.Nil(t, errorResp.Data)
}