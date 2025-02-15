package handlers

import (
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

	username := r.Context().Value("username").(string)

	// Выполняем покупку
	err := s.repo.BuyItem(r.Context(), username, item)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrMsgNotEnoughCoins):
			sendError(w, http.StatusBadRequest, "Not enough coins")
		case errors.Is(err, repository.ErrMsgOutOfStock):
			sendError(w, http.StatusBadRequest, "Item out of stock")
		case errors.Is(err, repository.ErrMsgProductNotExist):
			sendError(w, http.StatusNotFound, "Product not found")
		default:
			sendError(w, http.StatusInternalServerError, "Purchase failed")
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Item purchased successfully",
	})
	w.WriteHeader(http.StatusOK)
}
