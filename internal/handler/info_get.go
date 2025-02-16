package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetInfo(c *gin.Context) {
	userID, err := h.getUserIDFromHeaders(c)
	if err != nil {
		h.handleErr(c, err)
		return
	}

	userBalance, err := h.storage.GetUserBalance(c, userID)
	if err != nil {
		h.handleErr(c, err)
		return
	}

	var inventories []Inventory

	purchases, err := h.storage.GetPurchasesByUserID(c, userID)
	if err != nil {
		h.handleErr(c, err)
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
		h.handleErr(c, err)
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
