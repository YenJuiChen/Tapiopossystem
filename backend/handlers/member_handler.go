package handlers

import (
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
