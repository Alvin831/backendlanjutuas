package simple_test

import (
	"testing"
	"uas_backend/app/utils"

	"github.com/stretchr/testify/assert"
)

// Test SuccessResponse function
func TestSuccessResponse(t *testing.T) {
	message := "Data retrieved successfully"
	code := 200
	data := map[string]string{"name": "John"}

	response := utils.SuccessResponse(message, code, data)

	assert.Equal(t, message, response.Meta.Message)
	assert.Equal(t, code, response.Meta.Code)
	assert.Equal(t, "success", response.Meta.Status)
	assert.Equal(t, data, response.Data)
}

// Test ErrorResponse function
func TestErrorResponse(t *testing.T) {
	message := "Invalid input"
	code := 400
	data := map[string]string{"error": "field required"}

	response := utils.ErrorResponse(message, code, data)

	assert.Equal(t, message, response.Meta.Message)
	assert.Equal(t, code, response.Meta.Code)
	assert.Equal(t, "error", response.Meta.Status)
	assert.Equal(t, data, response.Data)
}

// Test SuccessResponse with nil data
func TestSuccessResponseNilData(t *testing.T) {
	response := utils.SuccessResponse("Success", 200, nil)

	assert.Equal(t, "Success", response.Meta.Message)
	assert.Equal(t, 200, response.Meta.Code)
	assert.Equal(t, "success", response.Meta.Status)
	assert.Nil(t, response.Data)
}

// Test ErrorResponse with nil data
func TestErrorResponseNilData(t *testing.T) {
	response := utils.ErrorResponse("Error occurred", 500, nil)

	assert.Equal(t, "Error occurred", response.Meta.Message)
	assert.Equal(t, 500, response.Meta.Code)
	assert.Equal(t, "error", response.Meta.Status)
	assert.Nil(t, response.Data)
}
