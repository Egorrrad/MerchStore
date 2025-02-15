package handlers

import (
	"MerchStore/src/internal/generated"
	"MerchStore/src/internal/repository"
	"MerchStore/src/internal/schemas"
	"encoding/json"
	"errors"
	"net/http"
)

func (s Server) PostApiSendCoin(w http.ResponseWriter, r *http.Request) {
	var req generated.PostApiSendCoinJSONRequestBody

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// Валидация
	if errorsFields := schemas.ValidateSendCoinRequest(req.ToUser, req.Amount); len(errorsFields) > 0 {
		sendValidationErrors(w, errorsFields)
		return
	}

	sender := r.Context().Value("username").(string)

	// Выполняем перевод
	if err := s.repo.SendCoins(r.Context(), sender, req.ToUser, req.Amount); err != nil {
		switch {
		case errors.Is(err, repository.ErrMsgNotEnoughCoins):
			sendError(w, http.StatusBadRequest, "Not enough coins")
		case errors.Is(err, repository.ErrMsgUserNotExist):
			sendError(w, http.StatusNotFound, "Receiver not found")
		default:
			sendError(w, http.StatusInternalServerError, "Transfer failed")
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Coins transferred successfully"})
}
