package validation

// Validator manages validation for a specific form
type Validator struct {
    errors map[string]string
}

// New creates a new form validator
func New() *Validator {
    return &Validator{
        errors: make(map[string]string),
    }
}

// Required validates a required field
func (v *Validator) Required(field, value, message string) *Validator {
    if IsEmpty(value) {
        if message == "" {
            message = "This field is required"
        }
        v.errors[field] = message
    }
    return v
}

// Email validates an email field
func (v *Validator) Email(field, value, message string) *Validator {
    // Skip if already has error or empty (should be caught by Required)
    if _, hasError := v.errors[field]; hasError || value == "" {
        return v
    }

    if !IsEmail(value) {
        if message == "" {
            message = "Please enter a valid email address"
        }
        v.errors[field] = message
    }
    return v
}

// Password validates password strength
func (v *Validator) Password(field, value, message string) *Validator {
    if _, hasError := v.errors[field]; hasError || value == "" {
        return v
    }

    if !IsStrongPassword(value) {
        if message == "" {
            message = "Password must include uppercase, lowercase, number, and special character"
        }
        v.errors[field] = message
    }
    return v
}

// Custom adds a custom validation
func (v *Validator) Custom(field string, isValid bool, message string) *Validator {
    if !isValid {
        v.errors[field] = message
    }
    return v
}

// Errors returns all validation errors
func (v *Validator) Errors() map[string]string {
    return v.errors
}

// IsValid returns true if there are no validation errors
func (v *Validator) IsValid() bool {
    return len(v.errors) == 0
}
