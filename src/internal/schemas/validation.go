package schemas

import (
	"fmt"
	"regexp"
)

type FieldError struct {
	Field   string
	Message string
}

func (e FieldError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func ValidateAuthRequest(username, password string) []FieldError {
	var fieldErrors []FieldError

	// Валидация username
	if username == "" {
		fieldErrors = append(fieldErrors, FieldError{"username", ErrRequired.Error()})
	} else {
		if len(username) < 3 {
			fieldErrors = append(fieldErrors, FieldError{"username", ErrMinLength.Error()})
		}
		if len(username) > 50 {
			fieldErrors = append(fieldErrors, FieldError{"username", ErrMaxLength.Error()})
		}
		if !isAlphanumeric(username) {
			fieldErrors = append(fieldErrors, FieldError{"username", ErrInvalidChars.Error()})
		}
	}

	// Валидация password
	if password == "" {
		fieldErrors = append(fieldErrors, FieldError{"password", ErrRequired.Error()})
	} else if len(password) < 6 {
		fieldErrors = append(fieldErrors, FieldError{"password", ErrMinLength.Error()})
	}

	return fieldErrors
}

// ValidateSendCoinRequest проверяет структуру SendCoinRequest
func ValidateSendCoinRequest(toUser string, amount int) []FieldError {
	var errors []FieldError

	// Валидация toUser
	if toUser == "" {
		errors = append(errors, FieldError{"toUser", ErrRequired.Error()})
	} else {
		if len(toUser) < 3 {
			errors = append(errors, FieldError{"toUser", ErrMinLength.Error()})
		}
		if len(toUser) > 50 {
			errors = append(errors, FieldError{"toUser", ErrMaxLength.Error()})
		}
		if !isAlphanumeric(toUser) {
			errors = append(errors, FieldError{"toUser", ErrInvalidChars.Error()})
		}
	}

	// Валидация amount
	if amount < 1 {
		errors = append(errors, FieldError{"amount", ErrInvalidAmount.Error()})
	}

	return errors
}

// ValidateItemName проверяет название товара
func ValidateItemName(item string) []FieldError {
	var errors []FieldError

	if item == "" {
		errors = append(errors, FieldError{"item", ErrRequired.Error()})
	} else {
		if len(item) < 2 {
			errors = append(errors, FieldError{"item", ErrMinLength.Error()})
		}
		if len(item) > 50 {
			errors = append(errors, FieldError{"item", ErrMaxLength.Error()})
		}
	}

	return errors
}

// Вспомогательные функции
func isAlphanumeric(s string) bool {
	matched, _ := regexp.MatchString("^[a-zA-Z0-9]+$", s)
	return matched
}
