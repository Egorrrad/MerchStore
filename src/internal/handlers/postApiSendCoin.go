package handlers

import (
	"MerchStore/src/internal/generated"
	"MerchStore/src/internal/logger"
	"MerchStore/src/internal/middleware"
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

	sender, ok := r.Context().Value(middleware.UsernameKey).(string)

	if !ok {
		logger.Logger.Error("Failed to extract username from context", "error", "invalid user context")
		sendError(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	// Выполняем перевод
	if err := s.repo.SendCoins(r.Context(), sender, req.ToUser, req.Amount); err != nil {
		switch {
		case errors.Is(err, repository.ErrMsgNotEnoughCoins):
			sendError(w, http.StatusBadRequest, "Not enough coins")
		case errors.Is(err, repository.ErrMsgSentToSelf):
			sendError(w, http.StatusBadRequest, "You can't transfer coins to self")
		case errors.Is(err, repository.ErrMsgUserNotExist):
			sendError(w, http.StatusBadRequest, "Receiver not found")
		default:
			logger.Logger.Error("Transfer failed", "sender", sender, "toUser", req.ToUser, "amount", req.Amount, "error", err)
			sendError(w, http.StatusInternalServerError, "Transfer failed")
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Coins transferred successfully"})
}
