package models

import "time"

type (
	Order struct {
		DBOrderID         int64     `json:"-" db:"id"`
		OrderUID          string    `json:"order_uid" db:"order_uid"`
		TrackNumber       string    `json:"track_number" db:"track_number"`
		Entry             string    `json:"entry" db:"entry"`
		Delivery          Delivery  `json:"delivery"`
		Payment           Payment   `json:"payment"`
		Items             []Items   `json:"items"`
		Locale            string    `json:"locale" db:"locale"`
		InternalSignature string    `json:"internal_signature" db:"internal_signature"`
		CustomerID        string    `json:"customer_id" db:"customer_id"`
		DeliveryService   string    `json:"delivery_service" db:"delivery_service"`
		ShardKey          string    `json:"shardkey" db:"shardkey"`
		SMID              int64     `json:"sm_id" db:"sm_id"`
		DateCreated       time.Time `json:"date_created" db:"date_created"`
		OofShard          string    `json:"oof_shard" db:"oof_shard"`
	}

	Delivery struct {
		ID      int64  `json:"-" db:"id"`
		OrderID int64  `json:"-" db:"order_id"`
		Name    string `json:"name" db:"delivery_name"`
		Phone   string `json:"phone" db:"phone"`
		ZIP     string `json:"zip" db:"zip"`
		City    string `json:"city" db:"city"`
		Address string `json:"address" db:"address"`
		Region  string `json:"region" db:"region"`
		Email   string `json:"email" db:"email"`
	}

	Payment struct {
		ID           int64  `json:"-" db:"id"`
		OrderID      int64  `json:"-" db:"order_id"`
		Transaction  string `json:"transaction" db:"transaction"`
		RequestID    string `json:"request_id" db:"request_id"`
		Currency     string `json:"currency" db:"currency"`
		Provider     string `json:"provider" db:"provider"`
		Amount       int    `json:"amount" db:"amount"`
		PaymentDT    int64  `json:"payment_dt" db:"payment_dt"`
		Bank         string `json:"bank" db:"bank"`
		DeliveryCost int    `json:"delivery_cost" db:"delivery_cost"`
		GoodsTotal   int    `json:"goods_total" db:"goods_total"`
		CustomFee    int    `json:"custom_fee" db:"custom_fee"`
	}

	Items struct {
		ID          int64  `json:"-" db:"id"`
		OrderID     int64  `json:"-" db:"order_id"`
		ChrtID      int64  `json:"chrt_id" db:"chrt_id"`
		TrackNumber string `json:"track_number" db:"track_number"`
		Price       int64  `json:"price" db:"price"`
		RID         string `json:"rid" db:"rid"`
		Name        string `json:"name" db:"item_name"`
		Sale        int    `json:"sale" db:"sale"`
		Size        string `json:"size" db:"item_size"`
		TotalPrice  int64  `json:"total_price" db:"total_price"`
		NMID        int64  `json:"nm_id" db:"nm_id"`
		Brand       string `json:"brand" db:"brand"`
		Status      int    `json:"status" db:"status"`
	}
)
