package validator

import (
	"fmt"
	"regexp"
)

// ValidateHash validate hash format
// Parameters:
//   - hashString: hash string
//
// Returns:
//   - error
func ValidateHash(hashString string) error {
	regex := regexp.MustCompile(RegexHash)
	if !regex.MatchString(hashString) {
		return fmt.Errorf("invalid hash string: %s", hashString)
	}
	return nil
}
