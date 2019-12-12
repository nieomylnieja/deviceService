package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerifyMongoDBName_DifferentLength(t *testing.T) {
	var longName string
	for longName = ""; len(longName) < 64; longName = longName + "a" {
	}
	tests := map[string]struct {
		input        string
		returnsError bool
	}{
		"correct name len":         {input: "mydb", returnsError: false},
		"lower boundary condition": {input: "", returnsError: true},
		"upper boundary condition": {input: longName, returnsError: true},
	}

	var err error
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err = verifyMongoDBName(tc.input)

			if tc.returnsError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVerifyMongoDBName_GivenNameWithNotAllowedChars_FuncReturnsError(t *testing.T) {
	tests := map[string]struct {
		input        string
		returnsError bool
	}{
		"case1": {input: "mydb\\12", returnsError: true},
		"case2": {input: "mydb\"12", returnsError: true},
		"case3": {input: "my.db12", returnsError: true},
		"case4": {input: "my$db12", returnsError: true},
		"case5": {input: "my db", returnsError: true},
	}

	var err error
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err = verifyMongoDBName(tc.input)

			assert.Error(t, err)
		})
	}
}
