package handlers

import (
	"MerchStore/src/internal/generated"
	"MerchStore/src/internal/schemas"
	"net/http"
)

func sendValidationErrors(w http.ResponseWriter, errors []schemas.FieldError) {
	errorsStr := ""
	for i, err := range errors {
		if i > 0 {
			errorsStr += ","
		}
		errorsStr += err.Error()
	}
	respondJSON(w, http.StatusBadRequest, generated.ErrorResponse{Errors: &errorsStr})
}

func sendError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, generated.ErrorResponse{Errors: &message})
}
