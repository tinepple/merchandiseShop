package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) ItemBuy(c *gin.Context) {
	itemValue := c.Param("item")
	if itemValue == "" {
		h.handleErr(c, fmt.Errorf("%w: item param is empty", validationError))
		return
	}

	userID, err := h.getUserIDFromHeaders(c)
	if err != nil {
		h.handleErr(c, err)
		return
	}

	item, err := h.storage.GetItem(c, itemValue)
	if err != nil {
		h.handleErr(c, err)
		return
	}

	userBalance, err := h.storage.GetUserBalance(c, userID)
	if err != nil {
		h.handleErr(c, err)
		return
	}

	if userBalance < item.Price {
		h.handleErr(c, fmt.Errorf("%w: not enough coins on balance", validationError))
		return
	}

	err = h.storage.CreatePurchase(c, userID, item.ID, userBalance-item.Price)
	if err != nil {
		h.handleErr(c, err)
		return
	}

	c.Status(http.StatusOK)
}
