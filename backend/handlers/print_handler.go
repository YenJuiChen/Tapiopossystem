package handlers

import (
	"fmt"
	"hcj-fdg-pos/models"
	"hcj-fdg-pos/queue"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// GenerateAndServePrintHTML 產生並提供 HTML 檔案
func GenerateAndServePrintHTML(c *fiber.Ctx) error {
	var data models.PrintRequest
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "無效資料"})
	}

	// ⏱ 產生唯一檔名
	timestamp := time.Now().Unix()
	htmlName := fmt.Sprintf("print_%d.html", timestamp)
	htmlPath := filepath.Join("static/generated", htmlName)

	// 🛠 確保資料夾存在
	if err := os.MkdirAll("static/generated", os.ModePerm); err != nil {
		return c.Status(500).SendString("無法建立目錄")
	}

	// 📄 載入 HTML 範本
	tmpl, err := template.ParseFiles("templates/print_template.html")
	if err != nil {
		return c.Status(500).SendString("無法載入範本: " + err.Error())
	}

	// ✏️ 產生 HTML 並儲存
	file, err := os.Create(htmlPath)
	if err != nil {
		return c.Status(500).SendString("無法建立 HTML: " + err.Error())
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return c.Status(500).SendString("無法渲染 HTML: " + err.Error())
	}

	// ✅ 回傳產出的檔案 URL
	return c.JSON(fiber.Map{
		"status":   "success",
		"html_url": "/generated/" + htmlName,
	})
}

// ✅ WebSocket handler
func PrintWebSocketHandler(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func HandlePrintWebSocket(conn *websocket.Conn) {
	queue.ClientsMu.Lock()
	queue.Clients[conn] = true
	queue.ClientsMu.Unlock()

	for {
		// 保持連線，不處理訊息內容
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}

	queue.ClientsMu.Lock()
	delete(queue.Clients, conn)
	queue.ClientsMu.Unlock()
	conn.Close()
}
