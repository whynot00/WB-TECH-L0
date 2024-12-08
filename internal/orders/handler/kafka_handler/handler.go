package kafkaHandler

import (
	"encoding/json"
	"kaf-interface/internal/orders/models"
	"kaf-interface/internal/orders/service"
	"log/slog"
)

type Handler struct {
	logger  *slog.Logger
	service *service.Service
}

func NewHandler(service *service.Service, logger *slog.Logger) *Handler {

	return &Handler{
		logger:  logger,
		service: service,
	}
}

func (h *Handler) MessageHandler(message []byte) error {

	var order models.Order

	if err := json.Unmarshal(message, &order); err != nil {
		return models.JSONUnmarshalError
	}

	// запись в кэш и бд
	if err := h.service.Orders.SetOrder(order); err != nil {
		return err
	}

	return nil
}

func (h *Handler) MigrateOrdersFromDB() {

	// подгрузка из бд в кэш при создании консьюмера
	if err := h.service.Orders.MigrateFromDB(); err != nil {
		h.logger.Error("error with migrate db to cache", slog.Any("error", err))
	}
}
