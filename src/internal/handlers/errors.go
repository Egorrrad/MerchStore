package handlers

import (
	"MerchStore/src/internal/generated"
	"MerchStore/src/internal/schemas"
	"net/http"
)

func sendValidationErrors(w http.ResponseWriter, errors []schemas.FieldError) {
	errorMap := make(map[string]string)
	for _, err := range errors {
		errorMap[err.Field] = err.Message
	}
	respondJSON(w, http.StatusBadRequest, map[string]interface{}{"errors": errorMap})
}

func sendError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, generated.ErrorResponse{Errors: &message})
}
