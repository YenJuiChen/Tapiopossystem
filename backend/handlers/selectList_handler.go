package handlers

import (
	"database/sql"
	"hcj-fdg-pos/database"
	"hcj-fdg-pos/models"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetCategoryItems(c *fiber.Ctx) error {
	db, err := database.POSRECORDS.DB()
	if err != nil {
		return c.Status(500).SendString("資料庫連線失敗")
	}

	printable := c.Query("printableOnly")
	qrcode := c.Query("qrcodeOnly")

	printableOnly := printable == "1" || strings.ToLower(printable) == "true"
	qrcodeOnly := qrcode == "1" || strings.ToLower(qrcode) == "true"

	var conditions []string
	if printableOnly {
		conditions = append(conditions, "i.is_print = 1")
	}
	if qrcodeOnly {
		conditions = append(conditions, "i.is_qrcode = 1")
	}

	query := `
               SELECT c.id, c.name, i.id, i.name, i.is_print, i.is_qrcode
               FROM Categories c
               LEFT JOIN Items i ON i.category_id = c.id`
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY c.id, i.id"

	rows, err := db.Query(query)
	if err != nil {
		return c.Status(500).SendString("查詢失敗")
	}
	defer rows.Close()

	var (
		categories []models.Category
		currentCat *models.Category
	)
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

		if currentCat == nil || currentCat.ID != cID {
			if currentCat != nil && len(currentCat.Items) > 0 {
				categories = append(categories, *currentCat)
			}
			currentCat = &models.Category{ID: cID, Name: cName, Items: []models.Item{}}
		}
		if iID.Valid && iName.Valid {
			currentCat.Items = append(currentCat.Items, models.Item{
				ID:       int(iID.Int64),
				Name:     iName.String,
				IsPrint:  isPrint,
				IsQrcode: isQrcode,
			})
		}
	}

	if currentCat != nil && len(currentCat.Items) > 0 {
		categories = append(categories, *currentCat)
	}
	return c.JSON(categories)
}
