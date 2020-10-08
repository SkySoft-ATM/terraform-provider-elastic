package utils

import (
	"fmt"
	"testing"

	"gotest.tools/assert"
)

func TestParseTwoPartID(t *testing.T) {

	tests := []struct {
		input              string
		expectedLeftValue  string
		expectedRightValue string
	}{
		{"test:truc", "test", "truc"},
		{"machin:bidule", "machin", "bidule"},
	}

	for _, test := range tests {
		leftValue, rightValue, err := ParseTwoPartID(test.input, "leftValue", "rightValue")
		assert.Equal(t, test.expectedLeftValue, leftValue, "expecting left values to be equal")
		assert.Equal(t, test.expectedRightValue, rightValue, "expecting right values to be equal")
		assert.NilError(t, err, "expecting error to be nil")
	}
}

func TestParseTwoPartIDError(t *testing.T) {
	expectedError := fmt.Errorf("Unexpected ID format (%q). Expected %s:%s", "incorrectInput", "left", "right")
	leftValue, rightValue, err := ParseTwoPartID("incorrectInput", "left", "right")
	assert.Equal(t, leftValue, "", "expecting left value to be nil")
	assert.Equal(t, rightValue, "", "expecting right value to be nil")
	assert.Equal(t, expectedError.Error(), err.Error(), "expected errors to be equal")
}
