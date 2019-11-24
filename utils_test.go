package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_SetPageBounds(t *testing.T) {
	tests := map[string]struct {
		input    map[string]int
		expected map[string]int
	}{
		"limit equals zero": {input: map[string]int{"limit": 0, "page": 0, "len": 7},
			expected: map[string]int{"lower": 0, "upper": 7}},
		"last page, len limits upper bound": {input: map[string]int{"limit": 4, "page": 3, "len": 13},
			expected: map[string]int{"lower": 12, "upper": 13}},
		"default limit, len limits upper bound": {input: map[string]int{"limit": 100, "page": 0, "len": 56},
			expected: map[string]int{"lower": 0, "upper": 56}},
		"lower beyond len, empty slice": {input: map[string]int{"limit": 10, "page": 4, "len": 34},
			expected: map[string]int{"lower": 0, "upper": 0}},
		"lower equals len, one element": {input: map[string]int{"limit": 10, "page": 4, "len": 40},
			expected: map[string]int{"lower": 40, "upper": 40}},
		"normal page": {input: map[string]int{"limit": 5, "page": 2, "len": 100},
			expected: map[string]int{"lower": 10, "upper": 15}},
	}

	var actualLower, acutalUpper int
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actualLower, acutalUpper = setPageBounds(tc.input["limit"], tc.input["page"], tc.input["len"])

			assert.Equal(t, tc.expected["lower"], actualLower)
			assert.Equal(t, tc.expected["upper"], acutalUpper)
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
