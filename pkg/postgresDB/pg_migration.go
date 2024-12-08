package postgresDB

import (
	"github.com/jmoiron/sqlx"
)

func migration(db *sqlx.DB) error {

	var creationTables [4]*sqlx.Stmt

	creationTables[0], _ = db.Preparex(
		`CREATE TABLE IF NOT EXISTS orders 
		(
			id SERIAL PRIMARY KEY,
			order_uid VARCHAR(255),
			track_number VARCHAR(255),
			entry VARCHAR(255),
			locale VARCHAR(255),
			internal_signature VARCHAR(255),
			customer_id VARCHAR(255),
			delivery_service VARCHAR(255),
			shardkey VARCHAR(255),
			sm_id BIGINT,
			date_created TIMESTAMP,
			oof_shard VARCHAR(255)
		);`,
	)

	creationTables[1], _ = db.Preparex(
		`CREATE TABLE IF NOT EXISTS delivery
		(
			id SERIAL PRIMARY KEY,
			order_id BIGINT REFERENCES orders(id),
			delivery_name VARCHAR(255),
			phone VARCHAR(255),
			zip VARCHAR(255),
			city VARCHAR(255),
			address VARCHAR(255),
			region VARCHAR(255),
			email VARCHAR(255)
		);`,
	)

	creationTables[2], _ = db.Preparex(
		`CREATE TABLE IF NOT EXISTS payment
		(
			id SERIAL PRIMARY KEY,
			order_id BIGINT REFERENCES orders(id),
			transaction VARCHAR(255),
			request_id VARCHAR(255),
			currency VARCHAR(255),
			provider VARCHAR(255),
			amount INTEGER,
			payment_dt BIGINT,
			bank VARCHAR(255),
			delivery_cost INTEGER,
			goods_total INTEGER,
			custom_fee INTEGER
		);`,
	)

	creationTables[3], _ = db.Preparex(
		`CREATE TABLE IF NOT EXISTS items
		(
			id SERIAL primary key,
			order_id BIGINT references orders(id),
			chrt_id BIGINT,
			track_number VARCHAR(255),
			price BIGINT,
			rid VARCHAR(255),
			item_name VARCHAR(255),
			sale INTEGER,
			item_size VARCHAR(255),
			total_price BIGINT,
			nm_id BIGINT,
			brand VARCHAR(255),
			status INTEGER
		);`,
	)

	defer func() {
		for _, stmt := range creationTables {
			stmt.Close()
		}
	}()

	for _, stmt := range creationTables {

		if _, err := stmt.Exec(); err != nil {
			return err
		}
	}

	return nil
}
