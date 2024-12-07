package postgres

import (
	"kaf-interface/internal/orders/models"

	"github.com/jmoiron/sqlx"
)

type Orders interface {
	SetOrder(order models.Order) error
	GetOrders() ([]models.Order, error)
	GetPaymentByOrderID(orderID int64) (models.Payment, error)
	GetItemsByOrderID(orderID int64) ([]models.Items, error)
	GetDeliveryByOrderID(orderID int64) (models.Delivery, error)
}

type Repository struct {
	Orders
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Orders: NewOrdersRepository(db),
	}
}
