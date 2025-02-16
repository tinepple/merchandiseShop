package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	authorizationError = errors.New("произошла ошибка авторизации")
	validationError    = errors.New("произошла ошибка валидации")
	internalError      = errors.New("произошла внутренняя ошибка")
)

func (h *Handler) handleErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, authorizationError):
		c.JSON(http.StatusUnauthorized, ErrorResponse{Errors: err.Error()})
	case errors.Is(err, validationError):
		c.JSON(http.StatusBadRequest, ErrorResponse{Errors: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{Errors: internalError.Error()})
	}
}
