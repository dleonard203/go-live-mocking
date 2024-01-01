package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	testCases := []struct {
		name     string
		nums     []int
		expected int
	}{
		{
			name:     "empty set adds up to 0",
			nums:     []int{},
			expected: 0,
		},
		{
			name:     "all positive numbers",
			nums:     []int{1, 2, 3},
			expected: 6,
		},
		{
			name:     "all negative numbers",
			nums:     []int{-1, -2, -3},
			expected: -6,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Add(tc.nums...)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMultiply(t *testing.T) {
	testCases := []struct {
		name     string
		nums     []int
		expected int
	}{
		{
			name:     "empty set multiplies to 0",
			nums:     []int{},
			expected: 0,
		},
		{
			name:     "all positive numbers",
			nums:     []int{1, 2, 3},
			expected: 6,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Multiply(tc.nums...)
			assert.Equal(t, tc.expected, result)
		})
	}
}
