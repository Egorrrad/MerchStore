package handlers

import (
	"MerchStore/src/internal/datastorage/model"
	"MerchStore/src/internal/generated"
	"net/http"
)

func buildInventory(purchases []model.Purchase) *[]struct {
	Quantity *int    `json:"quantity,omitempty"`
	Type     *string `json:"type,omitempty"`
} {
	inventory := make(map[string]int)
	for _, p := range purchases {
		inventory[p.ProductName] += p.Quantity
	}

	result := make([]struct {
		Quantity *int    `json:"quantity,omitempty"`
		Type     *string `json:"type,omitempty"`
	}, 0, len(inventory))

	for name, qty := range inventory {
		nameCopy := name
		qtyCopy := qty
		result = append(result, struct {
			Quantity *int    `json:"quantity,omitempty"`
			Type     *string `json:"type,omitempty"`
		}{
			Type:     &nameCopy,
			Quantity: &qtyCopy,
		})
	}
	return &result
}

func (s *Server) GetApiInfo(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value("user").(string)
	if !ok {
		sendError(w, http.StatusUnauthorized, "Invalid user context")
		return
	}

	// Получаем основную информацию о пользователе
	user, err := s.repo.GetUserByUsername(r.Context(), username)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "User not found")
		return
	}

	purchases, _, err := s.repo.GetUserInfo(r.Context(), user.Username)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to get purchases and operations")
		return
	}

	// Формируем ответ
	response := generated.InfoResponse{
		Coins:     &user.Coins,
		Inventory: buildInventory(purchases),
	}

	respondJSON(w, http.StatusOK, response)
}
