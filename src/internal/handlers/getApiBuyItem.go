package handlers

import (
	"MerchStore/src/internal/logger"
	"MerchStore/src/internal/middleware"
	"MerchStore/src/internal/repository"
	"MerchStore/src/internal/schemas"
	"errors"
	"net/http"
)

func (s Server) GetApiBuyItem(w http.ResponseWriter, r *http.Request, item string) {
	// Валидация названия товара
	if fieldErrors := schemas.ValidateItemName(item); len(fieldErrors) > 0 {
		sendValidationErrors(w, fieldErrors)
		return
	}

	username, ok := r.Context().Value(middleware.UsernameKey).(string)
	if !ok {
		logger.Logger.Error("Failed to extract username from context", "error", "invalid user context")
		sendError(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	// Выполняем покупку
	err := s.repo.BuyItem(r.Context(), username, item)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrMsgNotEnoughCoins):
			sendError(w, http.StatusBadRequest, "Not enough coins")
		case errors.Is(err, repository.ErrMsgOutOfStock):
			sendError(w, http.StatusBadRequest, "Item out of stock")
		case errors.Is(err, repository.ErrMsgProductNotExist):
			sendError(w, http.StatusBadRequest, "Product not found")
		default:
			logger.Logger.Error("Unexpected error during purchase", "username", username, "item", item, "error", err)
			sendError(w, http.StatusInternalServerError, "Purchase failed")
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Item purchased successfully",
	})
}
