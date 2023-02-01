package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

type (
	Order struct {
		ID                string    `json:"id"`
		TrackNumber       string    `json:"track_number"`
		Entry             string    `json:"entry"`
		Delivery          Delivery  `json:"delivery"`
		Payment           Payment   `json:"payment"`
		Items             []Item    `json:"items"`
		Locale            string    `json:"locale"`
		InternalSignature string    `json:"internal_signature"`
		CustomerId        string    `json:"customer_id"`
		DeliveryService   string    `json:"delivery_service"`
		Shardkey          string    `json:"shardkey"`
		SmId              int       `json:"sm_id"`
		DateCreated       time.Time `json:"date_created"`
		OofShard          string    `json:"oof_shard"`
	}
	Delivery struct {
		Name    string `json:"name"`
		Phone   string `json:"phone"`
		Zip     string `json:"zip"`
		City    string `json:"city"`
		Address string `json:"address"`
		Region  string `json:"region"`
		Email   string `json:"email"`
	}
	Payment struct {
		Transaction  string `json:"transaction"`
		RequestId    string `json:"request_id"`
		Currency     string `json:"currency"`
		Provider     string `json:"provider"`
		Amount       int    `json:"amount"`
		PaymentDt    int    `json:"payment_dt"`
		Bank         string `json:"bank"`
		DeliveryCost int    `json:"delivery_cost"`
		GoodsTotal   int    `json:"goods_total"`
		CustomFee    int    `json:"custom_fee"`
	}
	Item struct {
		ChrtId      int    `json:"chrt_id"`
		TrackNumber string `json:"track_number"`
		Price       int    `json:"price"`
		Rid         string `json:"rid"`
		Name        string `json:"name"`
		Sale        int    `json:"sale"`
		Size        string `json:"size"`
		TotalPrice  int    `json:"total_price"`
		NmId        int    `json:"nm_id"`
		Brand       string `json:"brand"`
		Status      int    `json:"status"`
	}
)

// Value - implementation for driver.Valuer
func (a Order) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan - implementation for sql.Scanner
func (a *Order) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.Errorf("type assertion to []byte failed for %T = %v", value, value)
	}

	return json.Unmarshal(b, &a)
}
