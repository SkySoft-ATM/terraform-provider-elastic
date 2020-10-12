package utils

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gotest.tools/assert"
)

type testCase struct {
	val         interface{}
	f           schema.SchemaValidateFunc
	expectedErr *regexp.Regexp
}

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

func TestValidationStringInSlice(t *testing.T) {
	runTestCases(t, []testCase{
		{
			val: "ValidValue",
			f:   StringInSlice([]string{"ValidValue", "AnotherValidValue"}, false),
		},
		// ignore case
		{
			val: "VALIDVALUE",
			f:   StringInSlice([]string{"ValidValue", "AnotherValidValue"}, true),
		},
		{
			val:         "VALIDVALUE",
			f:           StringInSlice([]string{"ValidValue", "AnotherValidValue"}, false),
			expectedErr: regexp.MustCompile("expected [\\w]+ to be one of \\[ValidValue AnotherValidValue\\], got VALIDVALUE"),
		},
		{
			val:         "InvalidValue",
			f:           StringInSlice([]string{"ValidValue", "AnotherValidValue"}, false),
			expectedErr: regexp.MustCompile("expected [\\w]+ to be one of \\[ValidValue AnotherValidValue\\], got InvalidValue"),
		},
		{
			val:         1,
			f:           StringInSlice([]string{"ValidValue", "AnotherValidValue"}, false),
			expectedErr: regexp.MustCompile("expected type of [\\w]+ to be string"),
		},
	})
}

func runTestCases(t *testing.T, cases []testCase) {
	matchErr := func(errs []error, r *regexp.Regexp) bool {
		// err must match one provided
		for _, err := range errs {
			if r.MatchString(err.Error()) {
				return true
			}
		}

		return false
	}

	for i, tc := range cases {
		_, errs := tc.f(tc.val, "test_property")

		if len(errs) == 0 && tc.expectedErr == nil {
			continue
		}

		if len(errs) != 0 && tc.expectedErr == nil {
			t.Fatalf("expected test case %d to produce no errors, got %v", i, errs)
		}

		if !matchErr(errs, tc.expectedErr) {
			t.Fatalf("expected test case %d to produce error matching \"%s\", got %v", i, tc.expectedErr, errs)
		}
	}
}
