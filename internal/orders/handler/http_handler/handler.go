package httpHandler

import (
	"kaf-interface/internal/orders/config"
	"kaf-interface/internal/orders/service"
	"log/slog"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
	logger  *slog.Logger
}

func NewHandler(service *service.Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) InitRouters(config *config.Config) *gin.Engine {
	router := gin.New()
	gin.SetMode(gin.ReleaseMode)

	router.LoadHTMLGlob("html/index.html")

	router.Use(requestid.New()) // создание requestid для трейсинга запросов
	router.Use(h.Logging)       // логгирование запроса

	router.GET("/orders_id/:id", h.GetOrderByID) // http://localhost/orders_id/orderidtest1111

	return router
}
