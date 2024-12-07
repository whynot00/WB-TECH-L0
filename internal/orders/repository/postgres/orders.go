package postgres

import (
	"fmt"
	"kaf-interface/internal/orders/models"

	"github.com/jmoiron/sqlx"
)

type OrdersRepository struct {
	db *sqlx.DB
}

func NewOrdersRepository(db *sqlx.DB) *OrdersRepository {
	return &OrdersRepository{
		db: db,
	}
}

func (r *OrdersRepository) SetOrder(order models.Order) error {
	var (
		lastInsertedID int64
	)

	// запись основной информации о заказе в БД
	// получаем ID последней добавленной записи - он же order_id для таблиц items, delivery, payment
	stmtNamed, _ := r.db.PrepareNamed(`
		INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, 
		customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES (:order_uid, :track_number, :entry, :locale, :internal_signature, 
		:customer_id, :delivery_service, :shardkey, :sm_id, :date_created, :oof_shard)
		RETURNING id;
		`)

	defer stmtNamed.Close()

	if err := stmtNamed.Get(&lastInsertedID, order); err != nil {
		return fmt.Errorf("error main order: %w", err)
	}

	// обновляем модель, добавляя в нее order_id
	order.DBOrderID = lastInsertedID
	order.Delivery.OrderID = lastInsertedID
	order.Payment.OrderID = lastInsertedID
	for i := 0; i < len(order.Items); i++ {
		order.Items[i].OrderID = lastInsertedID
	}

	stmtNamed, _ = r.db.PrepareNamed(`
		INSERT INTO delivery (order_id, delivery_name, phone, zip, city, address, region, email)
		VALUES (:order_id, :delivery_name, :phone, :zip, :city, :address, :region, :email)
	`)

	if _, err := stmtNamed.Exec(order.Delivery); err != nil {
		return fmt.Errorf("error delivery: %w", err)
	}

	stmtNamed, _ = r.db.PrepareNamed(`
		INSERT INTO payment (order_id, transaction, request_id, currency, provider, amount,
		payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES (:order_id, :transaction, :request_id, :currency, :provider, :amount,
		:payment_dt, :bank, :delivery_cost, :goods_total, :custom_fee)
	`)

	if _, err := stmtNamed.Exec(order.Payment); err != nil {
		return fmt.Errorf("error payment: %w", err)
	}

	stmtNamed, _ = r.db.PrepareNamed(`
		INSERT INTO items (order_id, chrt_id, track_number, price, rid, item_name,
		sale, item_size, total_price, nm_id, brand, status)
		VALUES (:order_id, :chrt_id, :track_number, :price, :rid, :item_name,
		:sale, :item_size, :total_price, :nm_id, :brand, :status)
	`)

	for _, item := range order.Items {
		if _, err := stmtNamed.Exec(item); err != nil {
			return fmt.Errorf("error items: %w", err)
		}
	}

	return nil
}

func (r *OrdersRepository) GetOrders() ([]models.Order, error) {
	var orders []models.Order

	stmt, _ := r.db.Preparex(
		`
		SELECT * FROM orders
	`,
	)

	rows, err := stmt.Queryx()
	if err != nil {
		return orders, err
	}

	order := models.Order{}
	for rows.Next() {
		if err := rows.StructScan(&order); err != nil {
			return orders, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrdersRepository) GetPaymentByOrderID(orderID int64) (models.Payment, error) {

	stmt, _ := r.db.Preparex(`
		SELECT * FROM payment WHERE order_id=$1
	`)

	var payment models.Payment

	if err := stmt.QueryRowx(orderID).StructScan(&payment); err != nil {
		return payment, err
	}

	return payment, nil
}

func (r *OrdersRepository) GetItemsByOrderID(orderID int64) ([]models.Items, error) {

	stmt, _ := r.db.Preparex(`
		SELECT * FROM items WHERE order_id=$1
	`)

	var items []models.Items

	rows, err := stmt.Queryx(orderID)
	if err != nil {
		return items, nil
	}

	for rows.Next() {
		var item models.Items

		if err := rows.StructScan(&item); err != nil {
			return items, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *OrdersRepository) GetDeliveryByOrderID(orderID int64) (models.Delivery, error) {

	stmt, _ := r.db.Preparex(`
		SELECT * FROM delivery WHERE order_id=$1
	`)

	var delivery models.Delivery

	if err := stmt.QueryRowx(orderID).StructScan(&delivery); err != nil {
		return delivery, err
	}

	return delivery, nil
}
