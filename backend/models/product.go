package models

import "time"

// Product 產品
type Product struct {
	ID       int     `json:"id"`
	Category string  `json:"category"`
	Name     string  `json:"name"`
	Image    string  `json:"image"`
	Price    float64 `json:"price"`
}

// Record 紀錄
type Record struct {
	ID              int       `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Gender          string    `json:"gender" db:"gender"`
	Address         string    `json:"address" db:"address"`
	Phone           string    `json:"phone" db:"phone"`
	Category        string    `json:"category" db:"category"`
	ProductID       int       `json:"product_id" db:"product_id"`
	ProductName     string    `json:"product_name" db:"product_name"`
	Price           int       `json:"price" db:"price"`
	Quantity        int       `json:"quantity" db:"quantity"`
	Amount          int       `json:"amount" db:"amount"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	ServiceItem     string    `json:"service_item" db:"service_item"`
	PaymentMethod   string    `json:"payment_method" db:"payment_method"`
	Code            string    `json:"code" db:"code"`
	NeedCertificate bool      `json:"need_certificate" db:"need_certificate"`
	Info            string    `json:"info" db:"info"`
}
