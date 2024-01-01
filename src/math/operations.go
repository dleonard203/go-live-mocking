package math

// Add returns the sum of nums
func Add(nums ...int) int {
	total := 0
	for _, num := range nums {
		total += num
	}
	return total
}

// Multiply returns the product of nums
func Multiply(nums ...int) int {
	total := 0
	for _, num := range nums {
		total *= num
	}
	return total
}
