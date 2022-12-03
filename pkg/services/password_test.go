package services

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidatePassword(t *testing.T) {
	testTable := []struct {
		password string
		result   error
	}{
		{"12345", nil},
		{"123456", nil},
		{"123457easd", nil},
		{"12", fmt.Errorf("password min length 5 char")},
		{"1234", fmt.Errorf("password min length 5 char")},
		{"", fmt.Errorf("password min length 5 char")},
		{"      ", fmt.Errorf("password min length 5 char")},
	}
	t.Parallel()
	for _, test := range testTable {
		result := ValidatePassword(test.password)
		assert.Equal(t, test.result, result, fmt.Sprintf("Incorect return result expect : %v, result: %v, testData : %v", test.result, result, test))
	}
}

func TestCheckPasswordHash(t *testing.T) {
	testTable := []struct {
		password string
		hash     string
		result   bool
	}{
		{"1234", "$2a$14$vBdp5MC8yl3ULhi6qHeIS.CVWfr.m2JmX0XH9IDi9H9JdjNxH9it6", true},
		{"123", "$2a$14$vBdp5MC8yl3ULhi6qHeIS.CVWfr.m2JmX0XH9IDi9H9JdjNxH9it6", false},
	}
	t.Parallel()
	for _, test := range testTable {
		result := CheckPasswordHash(test.password, test.hash)
		assert.Equal(t, test.result, result, fmt.Sprintf("Incorect return result expect : %v, result: %v, testData : %v", test.result, result, test))
	}
}
