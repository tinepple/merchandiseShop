package handler

import (
	"MerchandiseShop/internal/storage"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) SendCoin(c *gin.Context) {
	userID, err := h.getUserIDFromHeaders(c)
	if err != nil {
		h.handleErr(c, err)
		return
	}

	var req SendCoinsRequest

	if err := c.BindJSON(&req); err != nil {
		h.handleErr(c, fmt.Errorf("%w: error parsing body", validationError))
		return
	}

	if req.Amount == 0 {
		h.handleErr(c, fmt.Errorf("%w: amount is empty", validationError))
		return
	}

	balanceUserFrom, err := h.storage.GetUserBalance(c, userID)
	if err != nil {
		h.handleErr(c, err)
		return
	}

	if balanceUserFrom < req.Amount {
		h.handleErr(c, fmt.Errorf("%w: not enough coins on balance", validationError))
		return
	}

	toUser, err := h.storage.GetUserByUsername(c, req.ToUser)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		h.handleErr(c, err)
		return
	}

	balanceUserTo, err := h.storage.GetUserBalance(c, toUser.ID)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		h.handleErr(c, err)
		return
	}

	transaction := storage.Transaction{
		UserIDFrom: userID,
		UserIDTo:   toUser.ID,
		Amount:     req.Amount,
	}

	err = h.storage.CreateTransaction(c, transaction, balanceUserFrom-req.Amount, balanceUserTo+req.Amount)
	if err != nil {
		h.handleErr(c, err)
		return
	}

	c.Status(http.StatusOK)
}
