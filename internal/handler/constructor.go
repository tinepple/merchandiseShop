package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	router      *gin.Engine
	storage     Storage
	authService authService
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) initRoutes() {
	h.router.POST("/api/auth", h.Auth)

	h.router.GET("/api/buy/:item", h.ItemBuy)
	h.router.POST("/api/sendCoin", h.SendCoin)
}

func New(storage Storage, authService authService) *Handler {
	h := &Handler{
		router:      gin.New(),
		storage:     storage,
		authService: authService,
	}

	h.initRoutes()

	return h
}
