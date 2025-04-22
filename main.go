package main

import (
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"net/http"
	"wsprotGame/command"
	regis "wsprotGame/command/register"
	"wsprotGame/command/response"
	proto2 "wsprotGame/proto/gen"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func handleWebSocket(registry *regis.CommandRegistry, w http.ResponseWriter, r *http.Request) {
	// 添加跨域头
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()
	// 获取远程地址并转换为 *net.TCPAddr 类型
	tcpAddr, ok := conn.RemoteAddr().(*net.TCPAddr)
	if !ok {
		log.Println("Failed to convert RemoteAddr to *net.TCPAddr")
	} else {
		// 直接获取完整的 IP 地址
		log.Println("Client IP:", tcpAddr.IP.String())
	}

	log.Println("Client connected") // 添加日志

	// 创建 ResponseSender 实例
	sender := &response.ResponseSender{}

	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		// 这里打印出来的是原始的二进制数据
		//log.Println("Received message:", raw) // 添加日志

		var gameMsg proto2.GameMessage
		if err := proto.Unmarshal(raw, &gameMsg); err != nil {
			log.Println("Failed to unmarshal message:", err)
			continue
		}

		// 执行对应的命令
		registry.Handler(conn, gameMsg.Type, gameMsg.Data, sender)
	}
}

func main() {
	//初始化命令注册表
	registry := regis.NewCommandRegistry()
	for msgType, cmd := range command.CommandMap {
		// 将所有的协议注册进去
		registry.Register(msgType, cmd)
	}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(registry, w, r)
	})
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
