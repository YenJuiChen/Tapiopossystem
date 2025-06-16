package handlers

import (
	"encoding/base64"
	"fmt"
	"hcj-fdg-pos/database"
	"hcj-fdg-pos/models"
	"hcj-fdg-pos/queue"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var staticPath string

// allowList defines permitted columns and sort orders for listing APIs
var allowedSortFields = map[string]bool{
	"id":             true,
	"name":           true,
	"gender":         true,
	"address":        true,
	"phone":          true,
	"category":       true,
	"product_id":     true,
	"product_name":   true,
	"price":          true,
	"quantity":       true,
	"amount":         true,
	"payment_method": true,
	"info":           true,
	"created_at":     true,
	"code":           true,
}

var allowedSortOrders = map[string]bool{
	"ASC":  true,
	"DESC": true,
}

func New(path string) {
	staticPath = path
}

// ✅ 建立紀錄並產生 Excel 並透過 WebSocket 推送 HTML
func CreateRecord(c *fiber.Ctx) error {
	var record models.Record
	if err := c.BodyParser(&record); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Invalid data"})
	}

	db, err := database.POSRECORDS.DB()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "fail", "message": "無法連接資料庫"})
	}

	record.CreatedAt = time.Now()

	// 查詢是否重複 code
	skipPrint := false
	if record.Code != "" {
		var count int
		db.QueryRow("SELECT COUNT(*) FROM Records WHERE code = ?", record.Code).Scan(&count)
		if count > 0 {
			skipPrint = true
		}
	}

	// 寫入 DB
	insertResult, err := db.Exec(`INSERT INTO Records 
        (name, gender, address, phone, category, product_id, product_name, price, quantity, amount, payment_method, code, need_certificate, info, created_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		record.Name, record.Gender, record.Address, record.Phone, record.Category,
		record.ProductID, record.ProductName, record.Price, record.Quantity,
		record.Amount, record.PaymentMethod, record.Code, record.NeedCertificate,
		record.Info, record.CreatedAt,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "fail", "message": "儲存資料失敗"})
	} else {
		id, err := insertResult.LastInsertId()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "fail", "message": "無法取得新增的 ID"})
		}
		record.ID = int(id)
	}

	// ✅ 若有 code 且已存在，或是選擇不列印
	if skipPrint || !record.NeedCertificate {
		return c.JSON(fiber.Map{"status": "success", "message": "已存在，不列印"})
	}

	tmpl, err := template.ParseFiles("templates/template_test2.html")
	// tmpl, err := template.ParseFiles("templates/print_template.html")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "fail", "message": "無法載入模板"})
	}
	title := ""
	switch record.Gender {
	case "男":
		title = "先生"
	case "女":
		title = "小姐"
	default:
		title = "先生/小姐"
	}
	// 組出要給模板用的資料
	data := models.PrintRequest{
		SerialNo:    fmt.Sprintf("%s%05d", time.Now().Format("20060102"), record.ID),
		Name:        record.Name,
		Title:       title,
		Amount:      fmt.Sprintf("%d", record.Amount),
		Phoen:       record.Phone,
		Date:        record.CreatedAt.Format("2006-01-02"),
		Address:     record.Address,
		ServiceName: record.ProductName, // ✅ 對應 template 裡用的 .ServiceName
	}

	var htmlBuilder strings.Builder
	if err := tmpl.Execute(&htmlBuilder, data); err != nil {
		fmt.Printf("Error executing template: %s\n", err.Error())
		return c.Status(500).JSON(fiber.Map{"status": "fail", "message": "HTML 組裝失敗"})
	}

	html := htmlBuilder.String()
	html = inlineImageBase64All(html, "templates/assets") // 或你實際圖片所在資料夾

	// ✅ 將產出的 HTML 寫入本地檔案（測試用）
	// testPath := fmt.Sprintf("print_test_%d.html", time.Now().Unix())
	// err = os.WriteFile(testPath, []byte(html), 0644)
	// if err != nil {
	// 	fmt.Printf("⚠️ 無法寫入測試 HTML 檔案: %v\n", err)
	// } else {
	// 	fmt.Printf("✅ 測試 HTML 已寫入: %s\n", testPath)
	// }

	// outputPath := fmt.Sprintf("output_%d.pdf", time.Now().Unix())
	// if err := exportHTMLToPDF(html, outputPath); err != nil {
	// 	fmt.Printf("⚠️ 產出 PDF 失敗: %v\n", err)
	// }

	// ✅ 發送給 WebSocket 所有已連線的前端
	queue.ClientsMu.Lock()
	for conn := range queue.Clients {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(html)); err != nil {
			conn.Close()
			delete(queue.Clients, conn)
		}
	}
	queue.ClientsMu.Unlock()

	return c.JSON(fiber.Map{"status": "success", "message": "訂單已送出"})
}

// 批次寫入紀錄
func ScanRecords(c *fiber.Ctx) error {
	type ScanReq struct {
		Items []struct {
			ItemId int    `json:"item_id"`
			Code   string `json:"code"`
		} `json:"items"`
		PaymentType string `json:"payment_type"`
	}
	var req ScanReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Invalid request"})
	}
	fmt.Printf("ScanRecords: %v\n", req)

	if len(req.Items) == 0 {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Items cannot be empty"})
	}

	db, err := database.POSRECORDS.DB()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "fail", "message": "資料庫連線失敗"})
	}

	// 將 item_id 去重查詢商品資訊
	idSet := make(map[int]bool)
	var distinctIDs []int
	for _, item := range req.Items {
		if !idSet[item.ItemId] {
			idSet[item.ItemId] = true
			distinctIDs = append(distinctIDs, item.ItemId)
		}
	}

	query := fmt.Sprintf(
		"SELECT id, name, category_id, price FROM Items WHERE id IN (%s)",
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(distinctIDs)), ","), "[]"),
	)

	rows, err := db.Query(query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "fail", "message": "查詢商品失敗"})
	}
	defer rows.Close()

	// 建立 map[id]商品資訊
	type ItemData struct {
		Name       string
		CategoryID string
		Price      int
	}
	itemMap := make(map[int]ItemData)
	for rows.Next() {
		var id, price int
		var name, category string
		if err := rows.Scan(&id, &name, &category, &price); err != nil {
			continue
		}
		itemMap[id] = ItemData{name, category, price}
	}

	var values []string
	var args []interface{}
	now := time.Now()

	for _, item := range req.Items {
		data, ok := itemMap[item.ItemId]
		if !ok {
			continue
		}

		values = append(values, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		args = append(args,
			"", "", "", "", // name, gender, address, phone
			data.CategoryID, // category
			item.ItemId,     // product_id
			data.Name,       // product_name
			data.Price,      // price
			1, data.Price,   // quantity, amount
			req.PaymentType, // payment_method
			item.Code,       // code
			0,               // need_certificate
			"",              // info
			now,             // created_at
		)
	}

	if len(values) == 0 {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "找不到對應商品"})
	}

	insertSQL := fmt.Sprintf(`INSERT INTO Records
        (name, gender, address, phone, category, product_id, product_name, price, quantity, amount, payment_method, code, need_certificate, info, created_at)
        VALUES %s`, strings.Join(values, ","))

	_, err = db.Exec(insertSQL, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "fail", "message": "寫入資料失敗", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "批次寫入完成"})
}

// ✅ 取得紀錄列表
func ListRecords(c *fiber.Ctx) error {
	db, err := database.POSRECORDS.DB()
	if err != nil {
		return c.Status(500).SendString("資料庫連線失敗")
	}

	// 查詢參數
	category := c.Query("category")
	item := c.Query("item")
	startDate := c.Query("startDate") // 格式：YYYY-MM-DD HH:mm
	endDate := c.Query("endDate")     // 格式：YYYY-MM-DD HH:mm
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	sortBy := c.Query("sortBy", "created_at")
	order := strings.ToUpper(c.Query("order", "DESC"))

	if !allowedSortFields[sortBy] {
		sortBy = "created_at"
	}
	if !allowedSortOrders[order] {
		order = "DESC"
	}

	var conditions []string
	var args []interface{}

	if category != "" {
		conditions = append(conditions, "category = ?")
		args = append(args, category)
	}
	if item != "" {
		conditions = append(conditions, "product_name = ?")
		args = append(args, item)
	}

	// ✅ 精準時間區間過濾
	if startDate != "" && endDate != "" {
		conditions = append(conditions, "created_at BETWEEN ? AND ?")
		args = append(args, startDate+":00", endDate+":59") // 加上秒
	} else if startDate != "" {
		conditions = append(conditions, "created_at >= ?")
		args = append(args, startDate+":00")
	} else if endDate != "" {
		conditions = append(conditions, "created_at <= ?")
		args = append(args, endDate+":59")
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// 查總筆數
	countSQL := fmt.Sprintf(`
		SELECT COUNT(*) FROM (
			SELECT * FROM (
				SELECT * FROM Records WHERE code IS NULL OR code = ''
				UNION ALL
				SELECT r.* FROM Records r
				JOIN (
					SELECT code, MIN(created_at) AS created_at
					FROM Records
					WHERE code IS NOT NULL AND code != ''
					GROUP BY code
				) sub ON r.code = sub.code AND r.created_at = sub.created_at
			) all_records
			%s
		) final`, whereClause)

	var total int
	if err := db.QueryRow(countSQL, args...).Scan(&total); err != nil {
		return c.Status(500).SendString("計算筆數失敗")
	}

	// 查資料
	sql := fmt.Sprintf(`
        SELECT id, name, gender, address, phone, category, product_id, product_name, price, quantity, amount, payment_method, need_certificate, info, created_at
        FROM (
                SELECT id, name, gender, address, phone, category, product_id, product_name, price, quantity, amount, payment_method, need_certificate, info, created_at
                FROM Records
                WHERE code IS NULL OR code = ''

                UNION ALL

                SELECT r.id, r.name, r.gender, r.address, r.phone, r.category, r.product_id, r.product_name, r.price, r.quantity, r.amount, r.payment_method, r.need_certificate, r.info, r.created_at
                FROM Records r
                JOIN (
                        SELECT code, MIN(created_at) AS created_at
                        FROM Records
                        WHERE code IS NOT NULL AND code != ''
			GROUP BY code
		) sub ON r.code = sub.code AND r.created_at = sub.created_at
	) all_records
	%s
	ORDER BY %s %s
	LIMIT ? OFFSET ?`, whereClause, sortBy, order)

	argsWithLimit := append(args, limit, offset)

	rows, err := db.Query(sql, argsWithLimit...)
	if err != nil {
		return c.Status(500).SendString("查詢資料失敗")
	}
	defer rows.Close()

	result := []models.Record{}
	for rows.Next() {
		var r models.Record
		err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.Gender,
			&r.Address,
			&r.Phone,
			&r.Category,
			&r.ProductID,
			&r.ProductName,
			&r.Price,
			&r.Quantity,
			&r.Amount,
			&r.PaymentMethod,
			&r.NeedCertificate,
			&r.Info,
			&r.CreatedAt,
		)
		if err != nil {
			return c.Status(500).SendString("資料格式錯誤")
		}

		switch r.PaymentMethod {
		case "cash":
			r.PaymentMethod = "現金"
		case "linepay":
			r.PaymentMethod = "LINE Pay"
		default:
			r.PaymentMethod = "其他"
		}

		result = append(result, r)
	}

	return c.JSON(fiber.Map{
		"total": total,
		"data":  result,
	})
}

// ListCertificateRecords returns records that have a certificate code
func ListCertificateRecords(c *fiber.Ctx) error {
	db, err := database.POSRECORDS.DB()
	if err != nil {
		return c.Status(500).SendString("資料庫連線失敗")
	}

	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	sortBy := c.Query("sortBy", "created_at")
	order := strings.ToUpper(c.Query("order", "DESC"))

	if !allowedSortFields[sortBy] {
		sortBy = "created_at"
	}
	if !allowedSortOrders[order] {
		order = "DESC"
	}

	var total int
	if err := db.QueryRow("SELECT COUNT(*) FROM Records WHERE need_certificate = 1").Scan(&total); err != nil {
		return c.Status(500).SendString("計算筆數失敗")
	}

	sql := fmt.Sprintf(`
               SELECT id, name, gender, address, phone, category, product_id, product_name, price, quantity, amount, payment_method, info, created_at
               FROM Records
               WHERE need_certificate = 1
               ORDER BY %s %s
               LIMIT ? OFFSET ?`, sortBy, order)

	rows, err := db.Query(sql, limit, offset)
	if err != nil {
		return c.Status(500).SendString("查詢資料失敗")
	}
	defer rows.Close()

	result := []models.Record{}
	for rows.Next() {
		var r models.Record
		err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.Gender,
			&r.Address,
			&r.Phone,
			&r.Category,
			&r.ProductID,
			&r.ProductName,
			&r.Price,
			&r.Quantity,
			&r.Amount,
			&r.PaymentMethod,
			&r.Info,
			&r.CreatedAt,
		)
		if err != nil {
			return c.Status(500).SendString("資料格式錯誤")
		}

		switch r.PaymentMethod {
		case "cash":
			r.PaymentMethod = "現金"
		case "linepay":
			r.PaymentMethod = "LINE Pay"
		default:
			r.PaymentMethod = "其他"
		}

		result = append(result, r)
	}

	return c.JSON(fiber.Map{
		"total": total,
		"data":  result,
	})
}

func inlineImageBase64All(html string, assetsDir string) string {
	// src="/assets/xxx.png"
	imgRe := regexp.MustCompile(`src="/assets/([^"]+)"`)
	html = imgRe.ReplaceAllStringFunc(html, func(match string) string {
		matches := imgRe.FindStringSubmatch(match)
		return embedImage(matches, assetsDir, "src")
	})

	// background-image: url("/assets/xxx.png")
	bgRe := regexp.MustCompile(`url\("/assets/([^"]+)"\)`)
	html = bgRe.ReplaceAllStringFunc(html, func(match string) string {
		matches := bgRe.FindStringSubmatch(match)
		return embedImage(matches, assetsDir, "background")
	})

	return html
}

func embedImage(matches []string, assetsDir string, mode string) string {
	if len(matches) < 2 {
		return matches[0]
	}

	imagePath := filepath.Join(assetsDir, matches[1])
	imgBytes, err := os.ReadFile(imagePath)
	if err != nil {
		fmt.Printf("讀圖失敗：%s\n", imagePath)
		return matches[0]
	}

	ext := filepath.Ext(imagePath)
	mimeType := "image/png"
	if ext == ".jpg" || ext == ".jpeg" {
		mimeType = "image/jpeg"
	} else if ext == ".gif" {
		mimeType = "image/gif"
	}

	base64Str := base64.StdEncoding.EncodeToString(imgBytes)

	if mode == "src" {
		return fmt.Sprintf(`src="data:%s;base64,%s"`, mimeType, base64Str)
	}
	return fmt.Sprintf(`url("data:%s;base64,%s")`, mimeType, base64Str)

}
