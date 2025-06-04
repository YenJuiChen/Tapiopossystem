package handlers

import (
	"database/sql"
	"fmt"
	"hcj-fdg-pos/database"
	"hcj-fdg-pos/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// 取得產品列表
func GetProducts(c *fiber.Ctx) error {
	category := c.Query("category_id")
	if category == "" {
		return c.Status(400).SendString("category is required")
	}
	categoryInt, err := strconv.Atoi(category)
	if err != nil {
		return c.Status(400).SendString("category must be an integer")
	}

	products, err := getProducts(categoryInt)
	if err != nil {
		return c.Status(500).SendString("Internal Server Error")
	}

	return c.JSON(products)
}

// 取得指定種類的產品
func getProducts(categoryID int) ([]models.Item, error) {
	db, err := database.POSRECORDS.DB()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, category_id, name, image, price 
		FROM Items 
		WHERE category_id = ?
		ORDER BY sort_order ASC
	`

	rows, err := db.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(&item.ID, &item.CategoryID, &item.Name, &item.Image, &item.Price)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// 取得產品詳細資訊
func GetItemDetail(c *fiber.Ctx) error {
	itemIDStr := c.Query("item_id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil || itemID <= 0 {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "無效的 item_id"})
	}
	fmt.Println("取得商品請求 itemID:", itemID)

	db, err := database.POSRECORDS.DB()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "fail", "message": "資料庫連線失敗"})
	}

	code := c.Query("code")
	if code != "" {
		var count int
		db.QueryRow("SELECT COUNT(*) FROM Records WHERE code = ?", code).Scan(&count)
		if count > 0 {
			return c.JSON(fiber.Map{
				"status": "fail",
				"data":   fiber.Map{},
			})
		}
	}

	query := `
		SELECT 
			i.id, i.name, i.price, i.image, i.is_print, i.is_qrcode,
			c.id AS category_id, c.name AS category_name
		FROM Items i
		JOIN Categories c ON i.category_id = c.id
		WHERE i.id = ?
		ORDER BY i.id ASC
	`

	var (
		id           int
		name         string
		price        int
		image        string
		categoryID   int
		categoryName string
		isPrint      bool
		isQrcode     bool
	)

	err = db.QueryRow(query, itemID).Scan(&id, &name, &price, &image, &isPrint, &isQrcode, &categoryID, &categoryName)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"status": "fail", "message": "找不到此項目"})
		}
		return c.Status(500).JSON(fiber.Map{"status": "fail", "message": "查詢失敗"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"id":            id,
			"name":          name,
			"price":         price,
			"image":         image,
			"category_id":   categoryID,
			"category_name": categoryName,
			"is_print":      isPrint,
			"is_qrcode":     isQrcode,
		},
	})
}

// 模擬資料庫資料
func getMockProducts(category string) ([]models.Product, error) {
	products, ok := mockData[category]
	if !ok {
		return []models.Product{}, nil // 或視情況 return error
	}
	return products, nil
}

// 模擬資料庫資料
var mockData = map[string][]models.Product{
	"service": {
		{ID: 1, Category: "service", Name: "開光加持", Image: "/images/service1.png", Price: 500},
		{ID: 2, Category: "service", Name: "超渡法會", Image: "/images/service2.png", Price: 1000},
	},
	"joss-paper": {
		{ID: 3, Category: "joss-paper", Name: "壽金100", Image: "/images/joss1.png", Price: 100},
		{ID: 4, Category: "joss-paper", Name: "刈金200", Image: "/images/joss2.png", Price: 200},
	},
	"donation": {
		{ID: 5, Category: "donation", Name: "香油錢 100", Image: "/images/donation1.png", Price: 100},
		{ID: 6, Category: "donation", Name: "隨喜捐款", Image: "/images/donation2.png", Price: 300},
	},
}
