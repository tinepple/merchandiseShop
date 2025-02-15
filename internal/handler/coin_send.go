package handler

import (
	"MerchandiseShop/internal/handler/utils"
	"MerchandiseShop/internal/storage"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) SendCoin(c *gin.Context) {
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

	var req SendCoinsRequest

	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if req.Amount == 0 {
		fmt.Println("amount is empty")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userBalance, err := h.storage.GetUserBalance(c, userID)
	if err != nil {
		fmt.Println(fmt.Errorf("h.storage.GetUserBalance: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if userBalance < req.Amount {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	toUser, err := h.storage.GetUserByUsername(c, req.ToUser)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	balanceUserFrom, err := h.storage.GetUserBalance(c, userID)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	balanceUserTo, err := h.storage.GetUserBalance(c, toUser.ID)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	transaction := storage.Transaction{
		UserIDFrom: userID,
		UserIDTo:   toUser.ID,
		Amount:     req.Amount,
	}

	err = h.storage.CreateTransaction(c, transaction, balanceUserFrom-req.Amount, balanceUserTo+req.Amount)
	if err != nil {
		fmt.Println(fmt.Errorf("h.storage.CreateTransaction: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
