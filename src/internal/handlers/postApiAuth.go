package handlers

import (
	"MerchStore/src/internal/generated"
	"MerchStore/src/internal/repository"
	"MerchStore/src/internal/schemas"
	"encoding/json"
	"errors"
	"net/http"
)

func (s Server) PostApiAuth(w http.ResponseWriter, r *http.Request) {
	var req generated.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// Валидация
	if errors := schemas.ValidateAuthRequest(req.Username, req.Password); len(errors) > 0 {
		sendValidationErrors(w, errors)
		return
	}

	token, err := s.repo.PostAuthUser(r.Context(), req.Username, req.Password)
	switch {
	default:

	}
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrMsgWrongPass):
			sendError(w, http.StatusBadRequest, "Wrong password")
		default:
			sendError(w, http.StatusInternalServerError, "Auth failed")
		}
	}
	respondJSON(w, http.StatusOK, generated.AuthResponse{Token: token})
}
