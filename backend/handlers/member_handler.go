package handlers

import (
	"strings"
	"time"

	"hcj-fdg-pos/database"

	"github.com/gofiber/fiber/v2"
)

// SearchMembers 搜尋會員資料
func SearchMembers(c *fiber.Ctx) error {
	q := c.Query("query")

	db, err := database.POSRECORDS.DB()
	if err != nil {
		return c.Status(500).SendString("資料庫連線失敗")
	}

	rows, err := db.Query(`
        SELECT id, name, gender, phone, address, created_at, amount, product_name, category, payment_method
        FROM Records
        WHERE name LIKE ? OR phone LIKE ?
        ORDER BY created_at DESC
        LIMIT 50`,
		"%"+q+"%", "%"+q+"%",
	)
	if err != nil {
		return c.Status(500).SendString("查詢失敗")
	}
	defer rows.Close()

	type member struct {
		ID            int       `json:"id"`
		Name          string    `json:"name"`
		Gender        string    `json:"gender"`
		Phone         string    `json:"phone"`
		Address       string    `json:"address"`
		CreatedAt     time.Time `json:"created_at"`
		Amount        int       `json:"amount"`
		ProductName   string    `json:"product_name"`
		Category      string    `json:"category"`
		PaymentMethod string    `json:"payment_method"`
	}

	var members []member
	for rows.Next() {
		var m member
		if err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Gender,
			&m.Phone,
			&m.Address,
			&m.CreatedAt,
			&m.Amount,
			&m.ProductName,
			&m.Category,
			&m.PaymentMethod,
		); err != nil {
			return c.Status(500).SendString("資料格式錯誤")
		}
		members = append(members, m)
	}

	return c.JSON(members)
}

// ListMemberOrders 查詢會員訂單紀錄
func ListMemberOrders(c *fiber.Ctx) error {
	phone := c.Query("phone")
	name := c.Query("name")

	db, err := database.POSRECORDS.DB()
	if err != nil {
		return c.Status(500).SendString("資料庫連線失敗")
	}

	var conditions []string
	var args []interface{}
	if phone != "" {
		conditions = append(conditions, "phone = ?")
		args = append(args, phone)
	}
	if name != "" {
		conditions = append(conditions, "name = ?")
		args = append(args, name)
	}

	if len(conditions) == 0 {
		return c.Status(400).SendString("phone or name required")
	}

	query := "SELECT id, created_at, amount, product_name, category, payment_method FROM Records WHERE " + strings.Join(conditions, " OR ") + " ORDER BY created_at DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return c.Status(500).SendString("查詢失敗")
	}
	defer rows.Close()

	type order struct {
		ID            int       `json:"id"`
		CreatedAt     time.Time `json:"created_at"`
		Amount        int       `json:"amount"`
		ProductName   string    `json:"product_name"`
		Category      string    `json:"category"`
		PaymentMethod string    `json:"payment_method"`
	}

	var orders []order
	for rows.Next() {
		var o order
		if err := rows.Scan(
			&o.ID,
			&o.CreatedAt,
			&o.Amount,
			&o.ProductName,
			&o.Category,
			&o.PaymentMethod,
		); err != nil {
			return c.Status(500).SendString("資料格式錯誤")
		}
		orders = append(orders, o)
	}

	return c.JSON(orders)
}
