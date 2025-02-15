package handler

import (
	"MerchandiseShop/internal/handler/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetInfo(c *gin.Context) {
	userID, err := h.authService.GetUserID(utils.GetTokenFromRequest(c))
	if err != nil {
		fmt.Println(fmt.Errorf("authService.GetUserID: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if userID == 0 {
		fmt.Println("userID is empty")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	//получить баланс
	userBalance, err := h.storage.GetUserBalance(c, userID)
	if err != nil {
		fmt.Println(fmt.Errorf("h.storage.GetUserBalance: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	//список покупок
	var inventories []Inventory

	purchases, err := h.storage.GetPurchasesByUserID(c, userID)
	if err != nil {
		fmt.Println(fmt.Errorf("h.storage.GetPurchasesByUserID: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	for _, purchase := range purchases {
		inventories = append(inventories, Inventory{
			Type:     purchase.Name,
			Quantity: purchase.Quantity,
		})
	}

	//история транзакций
	transactions, err := h.storage.GetTransactionsByUserID(c, userID)
	if err != nil {
		fmt.Println(fmt.Errorf("h.storage.GetTransactionsByUserID: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var (
		receivedCoins []ReceivedCoins
		sentCoins     []SentCoins
	)
	for _, transaction := range transactions {
		if transaction.UserIDTo == userID {
			receivedCoins = append(receivedCoins, ReceivedCoins{
				FromUser: transaction.UserNameFrom,
				Amount:   transaction.Amount,
			})
		}
		if transaction.UserIDFrom == userID {
			sentCoins = append(sentCoins, SentCoins{
				ToUser: transaction.UserNameTo,
				Amount: transaction.Amount,
			})
		}
	}

	c.JSON(http.StatusOK, InfoResponse{
		Coins:     userBalance,
		Inventory: inventories,
		CoinHistory: CoinHistory{
			Received: receivedCoins,
			Sent:     sentCoins,
		},
	})
}
