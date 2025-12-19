package utils

// CalculatePoints calculates achievement points based on competition level
func CalculatePoints(level string) int {
	switch level {
	case "Lokal (1-25 poin)":
		return 25
	case "Regional (26-50 poin)":
		return 50
	case "Nasional (51-75 poin)":
		return 75
	case "Internasional (76+ poin)":
		return 100
	default:
		return 0
	}
}

// CalculateAverage calculates average from slice of numbers
func CalculateAverage(numbers []int) float64 {
	if len(numbers) == 0 {
		return 0
	}
	
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	
	return float64(sum) / float64(len(numbers))
}

// GetMaxValue returns maximum value from slice
func GetMaxValue(numbers []int) int {
	if len(numbers) == 0 {
		return 0
	}
	
	max := numbers[0]
	for _, num := range numbers {
		if num > max {
			max = num
		}
	}
	
	return max
}

// GetMinValue returns minimum value from slice
func GetMinValue(numbers []int) int {
	if len(numbers) == 0 {
		return 0
	}
	
	min := numbers[0]
	for _, num := range numbers {
		if num < min {
			min = num
		}
	}
	
	return min
}