package service

import (
	"kaf-interface/internal/orders/config"
	"kaf-interface/internal/orders/models"
	"kaf-interface/internal/orders/repository/cache"
	"kaf-interface/internal/orders/repository/postgres"
)

type CacheService struct {
	cacheRepo *cache.Repository
	dbRepo    *postgres.Repository
	config    *config.Config
}

func NewOrdersService(cacheRepo *cache.Repository, dbRepo *postgres.Repository, cfg *config.Config) *CacheService {

	return &CacheService{
		cacheRepo: cacheRepo,
		dbRepo:    dbRepo,
		config:    cfg,
	}
}

func (c *CacheService) SetOrder(order models.Order) error {

	if order.OrderUID == "" {
		return models.EmptyOrderIDError
	}

	// запись в кэш
	if err := c.cacheRepo.Orders.SetOrderInCache(order); err != nil {
		return err
	}

	// запись в бд
	if err := c.dbRepo.Orders.SetOrder(order); err != nil {
		return err
	}

	return nil
}

func (c *CacheService) MigrateFromDB() error {
	orders, err := c.dbRepo.Orders.GetOrders()
	if err != nil {
		return err
	}

	// запись в кэш из бд при инициализации
	for i := 0; i < len(orders); i++ {
		orders[i].Delivery, err = c.dbRepo.Orders.GetDeliveryByOrderID(orders[i].DBOrderID)
		if err != nil {
			return err
		}

		orders[i].Items, err = c.dbRepo.Orders.GetItemsByOrderID(orders[i].DBOrderID)
		if err != nil {
			return err
		}

		orders[i].Payment, err = c.dbRepo.Orders.GetPaymentByOrderID(orders[i].DBOrderID)
		if err != nil {
			return err
		}
	}

	for _, order := range orders {
		if err := c.cacheRepo.Orders.SetOrderInCache(order); err != nil {
			return err
		}
	}

	return nil
}

func (c *CacheService) GetOrderByID(orderUID string) (*models.Order, error) {

	return c.cacheRepo.Orders.GetOrderFromCacheByID(orderUID)
}
