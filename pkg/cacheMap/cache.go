package cacheMap

import (
	"kaf-interface/internal/orders/models"
	"sync"
)

type CacheMap struct {
	m  map[string]models.Order
	mu sync.Mutex
}

// реализация простого in memory кэша, имеет только сеттер и геттер
func NewCacheMap() *CacheMap {

	return &CacheMap{
		m: make(map[string]models.Order),
	}
}

func (c *CacheMap) Set(order models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.m[order.OrderUID] = order
}

func (c *CacheMap) Get(key string) (*models.Order, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.m[key]
	if !ok {
		return nil, models.KeyNotExistsError
	}

	return &val, nil
}
