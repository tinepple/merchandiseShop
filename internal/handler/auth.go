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
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if req.Password == "" || req.Username == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pass, _ := utils.HashPassword(req.Password)
	fmt.Println("hashedPass:", pass)

	user, err := h.storage.GetUserByUsername(c, req.Username)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if errors.Is(err, storage.ErrNotFound) {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		user, err = h.storage.CreateUser(c, req.Username, hashedPassword)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	} else {
		if !utils.CheckPasswordHash(req.Password, user.Password) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}

	token, err := h.authService.GenerateJWT(user.ID)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
	})
}
