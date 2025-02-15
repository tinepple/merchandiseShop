package handler

import (
	"MerchandiseShop/internal/handler/utils"
	"MerchandiseShop/internal/storage"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) Auth(c *gin.Context) {
	var req AuthRequest

	if err := c.BindJSON(&req); err != nil {
		h.handleErr(c, fmt.Errorf("%w: error parsing body", validationError))
		return
	}

	if req.Password == "" || req.Username == "" {
		h.handleErr(c, fmt.Errorf("%w: username or password is empty", validationError))
		return
	}

	user, err := h.storage.GetUserByUsername(c, req.Username)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		h.handleErr(c, err)
		return
	}

	if errors.Is(err, storage.ErrNotFound) {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			h.handleErr(c, err)
			return
		}
		user, err = h.storage.CreateUser(c, req.Username, hashedPassword)
		if err != nil {
			h.handleErr(c, err)
			return
		}
	} else {
		if !utils.CheckPasswordHash(req.Password, user.Password) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}

	token, err := h.authService.GenerateJWT(user.ID)
	if err != nil {
		h.handleErr(c, err)
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
	})
}

func (h *Handler) getUserIDFromHeaders(ctx *gin.Context) (int, error) {
	userID, err := h.authService.GetUserID(utils.GetTokenFromRequest(ctx))
	if err != nil || userID <= 0 {
		return 0, authorizationError
	}

	return userID, nil
}
