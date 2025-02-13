package handler

import (
	"MerchandiseShop/internal/handler/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) ItemBuy(c *gin.Context) {
	itemValue := c.Param("item")
	if itemValue == "" {
		fmt.Println("param is empty")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

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

	item, err := h.storage.GetItem(c, itemValue)
	if err != nil {
		fmt.Println(fmt.Errorf("h.storage.GetItem: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userBalance, err := h.storage.GetUserBalance(c, userID)
	if err != nil {
		fmt.Println(fmt.Errorf("h.storage.GetUserBalance: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if userBalance < item.Price {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.storage.CreatePurchase(c, userID, item.ID, userBalance-item.Price)
	if err != nil {
		fmt.Println(fmt.Errorf("h.storage.CreatePurchase: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
