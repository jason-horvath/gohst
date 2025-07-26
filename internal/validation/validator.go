package validation

import "fmt"

// Validator manages validation for a specific form
type Validator struct {
    errors map[string][]string
}

// New creates a new form validator
func NewValidator() *Validator {
    return &Validator{
        errors: make(map[string][]string),
    }
}

// Required validates a required field
func (v *Validator) Required(field, value, message string) *Validator {
    if IsEmpty(value) {
        if message == "" {
            message = "This field is required"
        }
        v.errors[field] = append(v.errors[field], message)
    }
    return v
}

// Email validates an email field
func (v *Validator) Email(field, value, message string) *Validator {
    // Skip if already has error or empty (should be caught by Required)
    if len(v.errors[field]) > 0 || value == "" {
        return v
    }

    if !IsEmail(value) {
        if message == "" {
            message = "Please enter a valid email address"
        }
        v.errors[field] = append(v.errors[field], message)
    }
    return v
}

// Password validates password strength
func (v *Validator) Password(field, value, message string) *Validator {
    if len(v.errors[field]) > 0 || value == "" {
        return v
    }

    if !IsStrongPassword(value) {
        if message == "" {
            message = "Password must include uppercase, lowercase, number, and special character"
        }
        v.errors[field] = append(v.errors[field], message)
    }
    return v
}

// Custom adds a custom validation
func (v *Validator) Custom(field string, isValid bool, message string) *Validator {
    if !isValid {
        v.errors[field] = append(v.errors[field], message)
    }
    return v
}

// Numeric validates that a field contains only numeric characters
func (v *Validator) Numeric(field, value, message string) *Validator {
	if len(v.errors[field]) > 0 || value == "" {
		return v
	}

	if !IsNumeric(value) {
		if message == "" {
			message = "This field must be numeric"
		}
		v.errors[field] = append(v.errors[field], message)
	}
	return v
}

// Checks if two values match
func (v *Validator) Matches(field, value, matchValue, message string) *Validator {
	if len(v.errors[field]) > 0 || value == "" {
		return v
	}

	if value != matchValue {
		if message == "" {
			message = "Values do not match"
		}
		v.errors[field] = append(v.errors[field], message)
	}
	return v
}

// MinLength checks if a string is at least minLength characters long
func (v *Validator) MinSelected(field string, values []string, minCount int, message string) *Validator {
    if len(values) < minCount {
        if message == "" {
            message = fmt.Sprintf("Please select at least %d options", minCount)
        }
        v.errors[field] = append(v.errors[field], message)
    }
    return v
}

// Helper for default minCount = 1
func (v *Validator) RequiredSelected(field string, values []string, message string) *Validator {
    return v.MinSelected(field, values, 1, message)
}

// Errors returns all validation errors
func (v *Validator) Errors() map[string][]string {
    return v.errors
}

// IsValid returns true if there are no validation errors
func (v *Validator) IsValid() bool {
    return len(v.errors) == 0
}
