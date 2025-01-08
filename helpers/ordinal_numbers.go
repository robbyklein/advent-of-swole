package helpers

import "fmt"

func OrdinalNumbers(num int) string {
	if num <= 0 {
		return fmt.Sprintf("%dth", num) // Handle non-positive integers gracefully
	}

	// Handle special cases for 11, 12, and 13
	if num%100 >= 11 && num%100 <= 13 {
		return fmt.Sprintf("%dth", num)
	}

	switch num % 10 {
	case 1:
		return fmt.Sprintf("%dst", num)
	case 2:
		return fmt.Sprintf("%dnd", num)
	case 3:
		return fmt.Sprintf("%drd", num)
	default:
		return fmt.Sprintf("%dth", num)
	}
}
