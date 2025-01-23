package main

import (
	pb "ServerExample/generated" // 引入 Protobuf 生成的 Go 包，替換為你的 Protobuf 包路徑
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"     // 用於處理 WebSocket 連接
	"google.golang.org/protobuf/proto" // 用於 Protobuf 的序列化與反序列化
)

// 升級器，用於將 HTTP 請求升級為 WebSocket 連接
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 忽略跨域檢查，允許所有來源的請求升級為 WebSocket
		return true
	},
}

// WebSocket 處理邏輯
func wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebSocket 處理邏輯")

	// 將 HTTP 請求升級為 WebSocket 連接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket Upgrade Error: %v", err) // 輸出升級錯誤
		return
	}
	defer conn.Close() // 當函數返回時關閉連接

	log.Printf("WebSocket 連線成功")
	for {
		// 讀取來自客戶端的消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket Read Error: %v", err) // 輸出讀取錯誤
			break                                       // 如果出錯，退出循環
		}

		// 使用 Protobuf 解碼客戶端發送的消息
		var chatMsg pb.ChatMessage
		if err := proto.Unmarshal(msg, &chatMsg); err != nil {
			log.Printf("Failed to parse WebSocket message: %v", err) // 解碼失敗
			continue                                                 // 跳過本次處理
		}

		// 輸出收到的消息
		log.Printf("Received WebSocket message from %s: %s", chatMsg.Sender, chatMsg.Content)

		// 構建回應消息
		response := pb.ChatMessage{
			Sender:  "Server",            // 回應消息的發送者
			Content: "Message received!", // 回應的內容
		}

		respData, err := proto.Marshal(&response) // 將回應消息序列化為二進制格式
		if err != nil {
			log.Printf("Failed to marshal response: %v", err)
			return
		}

		// 發送回應消息給客戶端
		err = conn.WriteMessage(websocket.BinaryMessage, respData)
		if err != nil {
			log.Printf("Failed to send WebSocket message: %v", err)
			return
		}
		log.Printf("Sent WebSocket message: %s", response.Content)
	}
}

// REST API 處理邏輯
func restHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("REST API 處理邏輯")

	// 確保請求方法為 POST，否則返回 405 錯誤
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 讀取請求體
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// 使用 Protobuf 解碼請求體
	var req pb.HelloRequest
	if err := proto.Unmarshal(body, &req); err != nil {
		http.Error(w, "Failed to parse Protobuf", http.StatusBadRequest) // 如果解碼失敗，返回 400 錯誤
		return
	}

	log.Printf("REST API連線成功")

	// 構建回應消息
	response := pb.HelloResponse{
		Message: "Hello, " + req.Name + "!", // 根據請求中的 Name 構建回應消息
	}
	respData, _ := proto.Marshal(&response) // 將回應消息序列化為二進制格式

	// 設置回應頭，指定內容類型為 Protobuf 格式
	w.Header().Set("Content-Type", "application/x-protobuf")
	w.Write(respData) // 發送回應
	log.Printf(" REST API Received name: %s", req.Name)
}

func main() {
	// 註冊 WebSocket 路由
	http.HandleFunc("/ws", wsHandler)

	// 註冊 REST API 路由
	http.HandleFunc("/api/hello", restHandler)

	// 啟動 HTTP 服務器，監聽 8888 端口
	log.Println("Server is running on :8888")
	log.Fatal(http.ListenAndServe(":8888", nil)) // 如果服務器啟動失敗，輸出錯誤並退出
}
