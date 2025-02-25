package handlers

import (
	"MerchStore/src/internal/generated"
	"MerchStore/src/internal/logger"
	"MerchStore/src/internal/repository"
	"MerchStore/src/internal/schemas"
	"encoding/json"
	"errors"
	"net/http"
)

func (s Server) PostApiAuth(w http.ResponseWriter, r *http.Request) {
	var req generated.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Warn("Invalid request format", "error", err)
		sendError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// Валидация
	if fieldErrors := schemas.ValidateAuthRequest(req.Username, req.Password); len(fieldErrors) > 0 {
		sendValidationErrors(w, fieldErrors)
		return
	}

	token, err := s.repo.PostAuthUser(r.Context(), req.Username, req.Password)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrMsgWrongPass):
			sendError(w, http.StatusBadRequest, "Wrong password")
		default:
			logger.Logger.Error("Authentication failed", "username", req.Username, "error", err)
			sendError(w, http.StatusInternalServerError, "Auth failed")
		}
		return
	}
	respondJSON(w, http.StatusOK, generated.AuthResponse{Token: token})
}
