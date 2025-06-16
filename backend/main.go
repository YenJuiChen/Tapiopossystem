package main

import (
	"fmt"
	"hcj-fdg-pos/database"
	"hcj-fdg-pos/handlers"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	database.Init()
	app := fiber.New()

	// ✅ API 群組
	api := app.Group("/api")
	api.Get("/items", handlers.GetProducts)
	api.Post("/record", handlers.CreateRecord)
	api.Get("/records", handlers.ListRecords)
	api.Get("/certificate-records", handlers.ListCertificateRecords)
	api.Post("/scan-records", handlers.ScanRecords)
	api.Get("/category-items", handlers.GetCategoryItems)
	api.Get("/members", handlers.SearchMembers)
	api.Get("/item-detail", handlers.GetItemDetail)
	// ✅ WebSocket 路由
	app.Get("/ws/print", handlers.PrintWebSocketHandler)
	app.Use("/ws/print", websocket.New(handlers.HandlePrintWebSocket))

	// ✅ 提供前端靜態頁面
	if os.Getenv("ENVIRONMENT") == "dev" {
		app.Static("/", "../frontend/dist") // 使用容器外部 ../frontend/dist
	} else {
		app.Static("/", "./dist") // 使用容器內部 /root/dist
	}
	app.Static("/templates", "./templates") // 可提供 /templates/images/xx.png
	app.Static("/assets", "./templates/images")

	// ✅ fallback 回傳 index.html
	app.Use(func(c *fiber.Ctx) error {
		if c.Path() != "/" && len(c.Path()) >= 4 && c.Path()[:4] == "/api" {
			return c.SendStatus(fiber.StatusNotFound)
		}
		if os.Getenv("ENVIRONMENT") == "dev" {
			return c.SendFile("../frontend/dist/index.html")
		} else {
			return c.SendFile("./dist/index.html")
		}
	})

	port := fmt.Sprint(":", os.Getenv("PORT"))
	app.Listen(port)
}
