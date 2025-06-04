package queue

import (
	"sync"

	"github.com/gofiber/websocket/v2"
)

// Clients 儲存所有連線中的 WebSocket 用戶端（例如 print-agent）
var Clients = make(map[*websocket.Conn]bool)

// ClientsMu 是操作 Clients 的同步鎖
var ClientsMu sync.Mutex
