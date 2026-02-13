package validation

import (
	"strings"
)

// password should contain
// 1. at least 8 characters
// 2. at least one uppercase letter
// 3. at least one lowercase letter
// 4. at least one number
// 5. at least one special character
func ValidatePassword(password string) ValidationErrors {
	var errs ValidationErrors
	if len(password) < 8 {
		errs = append(errs, ValidationError{
			Field: "password",
			Param: "min",
			Code:  "min",
			Value: "8",
		})
	}
	if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		errs = append(errs, ValidationError{
			Field: "password",
			Param: "uppercase",
			Code:  "uppercase",
			Value: password,
		})
	}
	if !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
		errs = append(errs, ValidationError{
			Field: "password",
			Param: 1,
			Code:  "lowercase",
			Value: password,
		})
	}
	if !strings.ContainsAny(password, "0123456789") {
		errs = append(errs, ValidationError{
			Field: "password",
			Param: 1,
			Code:  "number",
			Value: password,
		})
	}
	if !strings.ContainsAny(password, "!@#$%^&*()_+"+"-=[]{}|;:,.<>?") {
		errs = append(errs, ValidationError{
			Field: "password",
			Param: 1,
			Code:  "special",
			Value: password,
		})
	}
	return errs
}
