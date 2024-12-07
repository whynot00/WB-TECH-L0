package cache

import (
	"kaf-interface/internal/orders/models"
	"kaf-interface/pkg/cacheMap"
)

type Orders interface {
	SetOrderInCache(order models.Order) error
	GetOrderFromCacheByID(orderID string) (*models.Order, error)
}

type Repository struct {
	Orders
}

func NewRepository(cache *cacheMap.CacheMap) *Repository {
	return &Repository{
		Orders: NewOrdersRepository(cache),
	}
}
