package validation

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

// Common validation errors
var (
    ErrRequired     = errors.New("this field is required")
    ErrInvalidEmail = errors.New("invalid email format")
    ErrWeakPassword = errors.New("password must include uppercase, lowercase, number, and special character")
)

// IsEmpty checks if a string is empty (after trimming whitespace)
func IsEmpty(value string) bool {
    return strings.TrimSpace(value) == ""
}

// IsEmail validates an email address format
func IsEmail(email string) bool {
    pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
    match, _ := regexp.MatchString(pattern, email)
    return match
}

// IsStrongPassword checks if a password meets security requirements
func IsStrongPassword(password string) bool {
    if len(password) < 8 {
        return false
    }

    var hasUpper, hasLower, hasNumber, hasSpecial bool
    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsNumber(char):
            hasNumber = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }

    return hasUpper && hasLower && hasNumber && hasSpecial
}

// IsNumeric checks if a string contains only numeric characters
func IsNumeric(s string) bool {
    for _, c := range s {
        if c < '0' || c > '9' {
            return false
        }
    }
    return len(s) > 0
}
