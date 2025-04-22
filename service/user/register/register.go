package register

import (
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"log"
	"sync"
	"wsprotGame/command/response"
	proto2 "wsprotGame/proto/gen"
)

// 模拟用户数据库
var users = make(map[string]string)
var usersMutex sync.Mutex

// RegisterRequestCommand 注册请求命令
type RegisterRequestCommand struct{}

func (c *RegisterRequestCommand) Execute(conn *websocket.Conn, data []byte, sender *response.ResponseSender) {
	var req proto2.RegisterRequest
	if err := proto.Unmarshal(data, &req); err != nil {
		sender.Send(conn, proto2.GameMessage_REGISTER_RESPONSE, &proto2.RegisterResponse{
			Success: false,
			Message: "解析请求失败",
		})
		return
	}
	log.Printf("发送注册成功响应，消息类型: %v", proto2.GameMessage_REGISTER_RESPONSE)
	log.Printf("Received RegisterRequest: Username = %s, Password = %s", req.Username, req.Password)
	usersMutex.Lock()
	defer usersMutex.Unlock()

	if _, exists := users[req.Username]; exists {
		sender.Send(conn, proto2.GameMessage_REGISTER_RESPONSE, &proto2.RegisterResponse{
			Success: false,
			Message: "用户名已存在",
		})
		return
	}

	users[req.Username] = req.Password
	sender.Send(conn, proto2.GameMessage_REGISTER_RESPONSE, &proto2.RegisterResponse{
		Success: true,
		Message: "注册成功",
	})

}

// RegisterResponseCommand 注册响应命令
type RegisterResponseCommand struct{}

func (c *RegisterResponseCommand) Execute(conn *websocket.Conn, data []byte, sender *response.ResponseSender) {
	var resp proto2.RegisterResponse
	if err := proto.Unmarshal(data, &resp); err != nil {
		log.Println("Failed to unmarshal RegisterResponse:", err)
		return
	}
	log.Printf("Received RegisterResponse: Success = %v, Message = %s", resp.Success, resp.Message)
	// 这里可以添加处理注册响应的逻辑

}
