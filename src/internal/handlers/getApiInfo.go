package handlers

import (
	"MerchStore/src/internal/repository"
	"MerchStore/src/internal/schemas"
	"MerchStore/src/internal/storage/model"
	"net/http"
)

func ptrInt(i int) *int       { return &i }
func ptrStr(s string) *string { return &s }

func buildInventory(purchases []model.Purchase) []schemas.InventoryItem {
	inventory := make(map[string]int)
	for _, p := range purchases {
		inventory[p.ProductName] += p.Quantity
	}

	result := make([]schemas.InventoryItem, 0, len(inventory))
	for name, qty := range inventory {
		result = append(result, schemas.InventoryItem{
			Type:     ptrStr(name),
			Quantity: ptrInt(qty),
		})
	}
	return result
}

func buildCoinHistory(operations []model.Operation, currentUserID int) *schemas.CoinHistory {
	history := &schemas.CoinHistory{
		Received: &[]schemas.Transaction{},
		Sent:     &[]schemas.Transaction{},
	}

	for _, op := range operations {
		// Определяем тип операции относительно текущего пользователя
		if op.ReceiverUserID == currentUserID {
			*history.Received = append(*history.Received, schemas.Transaction{
				FromUser: ptrStr(op.SenderUsername),
				Amount:   ptrInt(op.Amount),
			})
		} else if op.SenderUserID == currentUserID {
			*history.Sent = append(*history.Sent, schemas.Transaction{
				ToUser: ptrStr(op.ReceiverUsername),
				Amount: ptrInt(op.Amount),
			})
		}
	}

	return history
}

func (s Server) GetApiInfo(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value("username").(string)
	if !ok {
		sendError(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	user, purchases, operations, err := s.repo.GetUserInfo(r.Context(), username)
	if err == repository.ErrMsgUserNotExist {
		sendError(w, http.StatusInternalServerError, "user not found")
		return
	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, "failed to get user data")
		return
	}

	inventory := buildInventory(purchases)

	response := schemas.InfoResponse{
		Coins:       &user.Coins,
		Inventory:   &inventory,
		CoinHistory: buildCoinHistory(operations, user.UserID),
	}

	// Гарантируем не-nil значения для всех полей
	if response.Inventory == nil {
		emptyInventory := make([]schemas.InventoryItem, 0)
		response.Inventory = &emptyInventory
	}

	respondJSON(w, http.StatusOK, response)
}
