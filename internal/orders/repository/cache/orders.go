package cache

import (
	"kaf-interface/internal/orders/models"
	"kaf-interface/pkg/cacheMap"
)

type OrdersRepository struct {
	cache *cacheMap.CacheMap
}

func NewOrdersRepository(cache *cacheMap.CacheMap) *OrdersRepository {
	return &OrdersRepository{
		cache: cache,
	}
}

func (r *OrdersRepository) SetOrderInCache(order models.Order) error {
	if order.OrderUID == "" {
		return models.EmptyOrderIDError
	}

	r.cache.Set(order)

	return nil
}

func (r *OrdersRepository) GetOrderFromCacheByID(orderUID string) (*models.Order, error) {

	return r.cache.Get(orderUID)
}
