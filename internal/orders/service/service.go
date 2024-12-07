package service

import (
	"kaf-interface/internal/orders/config"
	"kaf-interface/internal/orders/models"
	"kaf-interface/internal/orders/repository/cache"
	"kaf-interface/internal/orders/repository/postgres"
)

type Orders interface {
	GetOrderByID(orderID string) (*models.Order, error)
	SetOrder(order models.Order) error
	MigrateFromDB() error
}

type Service struct {
	Orders
}

type Deps struct {
	DBRepo    *postgres.Repository
	CacheRepo *cache.Repository
	Config    *config.Config
}

func NewService(deps Deps) *Service {
	return &Service{
		Orders: NewOrdersService(deps.CacheRepo, deps.DBRepo, deps.Config),
	}
}
