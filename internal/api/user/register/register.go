package register

import (
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"log"
	"wsprotGame/internal/service"
	proto2 "wsprotGame/proto/gen"
	"wsprotGame/server/command/response"
)

// RegisterRequestCommand 注册请求命令
type RegisterRequestCommand struct {
	us service.UserService
}

func NewRegisterRequestCommand(userService service.UserService) *RegisterRequestCommand {
	return &RegisterRequestCommand{us: userService}
}

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
	err := c.us.HandleRegister(req)
	if err != nil {
		sender.Send(conn, proto2.GameMessage_REGISTER_RESPONSE, &proto2.RegisterResponse{
			Success: false,
			Message: "注册失败",
		})
	}
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
