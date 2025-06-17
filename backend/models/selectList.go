package models

type Item struct {
	ID         int     `json:"id" gorm:"primaryKey"`
	CategoryID int     `json:"category_id"`
	CategoryName string `json:"category_name"`
	Name       string  `json:"name"`
	Image      string  `json:"image"`
	Price      float64 `json:"price"`
	SortOrder  int     `json:"sort_order"`
	Enabled    bool    `json:"enabled"`
	IsPrint    bool    `json:"is_print"`
	IsQrcode   bool    `json:"is_qrcode"`
}

type Category struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Items []Item `json:"items"`
}
