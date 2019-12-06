package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_SetPageBoundsToInt64(t *testing.T) {
	tests := map[string]struct {
		input    map[string]int
		expected map[string]int64
	}{
		"case 1": {input: map[string]int{"limit": 0, "page": 0},
			expected: map[string]int64{"lower": 0, "upper": 0}},
		"case 2": {input: map[string]int{"limit": 4, "page": 3},
			expected: map[string]int64{"lower": 12, "upper": 16}},
		"case 3": {input: map[string]int{"limit": 100, "page": 0},
			expected: map[string]int64{"lower": 0, "upper": 100}},
		"case 4": {input: map[string]int{"limit": 10, "page": 4},
			expected: map[string]int64{"lower": 40, "upper": 50}},
		"case 5": {input: map[string]int{"limit": 5, "page": 2},
			expected: map[string]int64{"lower": 10, "upper": 15}},
	}

	var actualLower, actualUpper int64
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actualLower, actualUpper = setPageBoundsToInt64(tc.input["limit"], tc.input["page"])

			assert.Equal(t, tc.expected["lower"], actualLower)
			assert.Equal(t, tc.expected["upper"], actualUpper)
		})
	}
}

func Test_ConvertToPositiveInteger_GivenWrongInput_FuncReturnsError(t *testing.T) {
	tests := map[string]string{
		"char":            "a",
		"negative number": "-2",
		"float":           "0.0",
		"interface":       "{}",
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := convertToPositiveInteger(tc)
			assert.Error(t, err)
		})
	}
}

func Test_ConvertToPositiveInteger_GivenCorrectInput_FuncReturnsPositiveInt(t *testing.T) {
	tests := map[string]string{
		"zero":              "0",
		"non zero positive": "14",
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := convertToPositiveInteger(tc)

			assert.NoError(t, err)
			assert.IsType(t, 1, actual)
			assert.GreaterOrEqual(t, actual, 0)
		})
	}
}

func Test_ReadIntFromQueryParameter_GivenNoValueInParam_FuncReturnsDefault(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	result, err := readIntFromQueryParameter(req.URL, "limit", 100)

	assert.NoError(t, err)
	assert.Equal(t, 100, result)
}
