package utils

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

// ValidateValueFunc check arguments based on list of strings
func ValidateValueFunc(values []string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (we []string, errors []error) {
		value := v.(string)
		valid := false
		for _, role := range values {
			if value == role {
				valid = true
				break
			}
		}

		if !valid {
			errors = append(errors, fmt.Errorf("%s is an invalid value for argument %s", value, k))
		}
		return
	}
}

// ParseTwoPartID returns the pieces of id `left:right` as left, right
func ParseTwoPartID(id, left, right string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Unexpected ID format (%q). Expected %s:%s", id, left, right)
	}

	return parts[0], parts[1], nil
}
