package handlers

import (
	"database/sql"
	"hcj-fdg-pos/database"
	"hcj-fdg-pos/models"

	"github.com/gofiber/fiber/v2"
)

func GetCategoryItems(c *fiber.Ctx) error {
	db, err := database.POSRECORDS.DB()
	if err != nil {
		return c.Status(500).SendString("資料庫連線失敗")
	}

	rows, err := db.Query(`
               SELECT c.id, c.name, i.id, i.name, i.is_print, i.is_qrcode
               FROM Categories c
               LEFT JOIN Items i ON i.category_id = c.id AND i.is_print = 1
               ORDER BY c.id, i.id
       `)
	if err != nil {
		return c.Status(500).SendString("查詢失敗")
	}
	defer rows.Close()

	categoryMap := make(map[int]*models.Category)
	for rows.Next() {
		var cID int
		var cName string
		var iID sql.NullInt64
		var iName sql.NullString
		var isPrint bool
		var isQrcode bool
		if err := rows.Scan(&cID, &cName, &iID, &iName, &isPrint, &isQrcode); err != nil {
			return c.Status(500).SendString("資料格式錯誤")
		}

		if _, ok := categoryMap[cID]; !ok {
			categoryMap[cID] = &models.Category{ID: cID, Name: cName, Items: []models.Item{}}
		}
		if iID.Valid && iName.Valid {
			categoryMap[cID].Items = append(categoryMap[cID].Items, models.Item{
				ID:       int(iID.Int64),
				Name:     iName.String,
				IsPrint:  isPrint,
				IsQrcode: isQrcode,
			})
		}
	}

	var result []models.Category
	for _, cat := range categoryMap {
		if len(cat.Items) == 0 {
			// omit categories without printable items
			continue
		}
		result = append(result, *cat)
	}
	return c.JSON(result)
}
