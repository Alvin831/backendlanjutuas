package simple_test

import (
	"testing"
	"uas_backend/app/utils"

	"github.com/stretchr/testify/assert"
)

// Test CalculatePoints function
func TestCalculatePoints(t *testing.T) {
	assert.Equal(t, 25, utils.CalculatePoints("Lokal (1-25 poin)"))
	assert.Equal(t, 50, utils.CalculatePoints("Regional (26-50 poin)"))
	assert.Equal(t, 75, utils.CalculatePoints("Nasional (51-75 poin)"))
	assert.Equal(t, 100, utils.CalculatePoints("Internasional (76+ poin)"))
	assert.Equal(t, 0, utils.CalculatePoints("Unknown Level"))
	assert.Equal(t, 0, utils.CalculatePoints(""))
}

// Test CalculateAverage function
func TestCalculateAverage(t *testing.T) {
	// Normal cases
	assert.Equal(t, 3.0, utils.CalculateAverage([]int{1, 2, 3, 4, 5}))
	assert.Equal(t, 10.0, utils.CalculateAverage([]int{10, 10, 10}))
	assert.Equal(t, 2.5, utils.CalculateAverage([]int{1, 4}))
	
	// Edge cases
	assert.Equal(t, 0.0, utils.CalculateAverage([]int{})) // empty slice
	assert.Equal(t, 5.0, utils.CalculateAverage([]int{5})) // single element
}

// Test GetMaxValue function
func TestGetMaxValue(t *testing.T) {
	// Normal cases
	assert.Equal(t, 5, utils.GetMaxValue([]int{1, 2, 3, 4, 5}))
	assert.Equal(t, 10, utils.GetMaxValue([]int{10, 5, 8, 3}))
	assert.Equal(t, 100, utils.GetMaxValue([]int{100, 50, 75}))
	
	// Edge cases
	assert.Equal(t, 0, utils.GetMaxValue([]int{})) // empty slice
	assert.Equal(t, 5, utils.GetMaxValue([]int{5})) // single element
	assert.Equal(t, -1, utils.GetMaxValue([]int{-5, -3, -1})) // negative numbers
}

// Test GetMinValue function
func TestGetMinValue(t *testing.T) {
	// Normal cases
	assert.Equal(t, 1, utils.GetMinValue([]int{1, 2, 3, 4, 5}))
	assert.Equal(t, 3, utils.GetMinValue([]int{10, 5, 8, 3}))
	assert.Equal(t, 50, utils.GetMinValue([]int{100, 50, 75}))
	
	// Edge cases
	assert.Equal(t, 0, utils.GetMinValue([]int{})) // empty slice
	assert.Equal(t, 5, utils.GetMinValue([]int{5})) // single element
	assert.Equal(t, -5, utils.GetMinValue([]int{-5, -3, -1})) // negative numbers
}