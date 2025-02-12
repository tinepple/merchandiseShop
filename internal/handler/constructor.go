package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	router  *gin.Engine
	storage Storage
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) initRoutes() {
	h.router.GET("/test", h.Test)
}

func New(storage Storage) *Handler {
	h := &Handler{
		router:  gin.New(),
		storage: storage,
	}

	h.initRoutes()

	return h
}

func (h *Handler) Test(c *gin.Context) {
	c.JSON(http.StatusOK, "lol")
}
