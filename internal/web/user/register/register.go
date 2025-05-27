package register

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"wsprotGame/internal/domain"
	"wsprotGame/internal/service"
	proto2 "wsprotGame/proto/gen"
	"wsprotGame/server/command/response"
	"wsprotGame/server/connection"
)

// RegisterRequestCommand 注册请求命令
type RegisterRequestCommand struct {
	us service.UserService
}

func NewRegisterRequestCommand(userService service.UserService) *RegisterRequestCommand {
	return &RegisterRequestCommand{us: userService}
}

func (c *RegisterRequestCommand) Execute(conn *connection.ConnInfo, data []byte, sender *response.ResponseSender, ctx context.Context) error {
	var req proto2.RegisterRequest
	if err := proto.Unmarshal(data, &req); err != nil {
		sender.Send(conn.Conn, proto2.GameMessage_REGISTER_RESPONSE, &proto2.RegisterResponse{
			Success: false,
			Message: "解析请求失败",
		})
		return err
	}
	log.Printf("发送注册成功响应，消息类型: %v", proto2.GameMessage_REGISTER_RESPONSE)
	err := c.us.HandleRegister(ctx, domain.User{
		Password: req.Password,
		NickName: req.Username,
		Phone:    req.Phone,
	})
	if err != nil {
		sender.Send(conn.Conn, proto2.GameMessage_REGISTER_RESPONSE, &proto2.RegisterResponse{
			Success: false,
			Message: "邮箱冲突",
		})
		return err
	}
	sender.Send(conn.Conn, proto2.GameMessage_REGISTER_RESPONSE, &proto2.RegisterResponse{
		Success: true,
		Message: "注册成功",
	})
	return nil

}

// RegisterResponseCommand 注册响应命令
type RegisterResponseCommand struct{}

func (c *RegisterResponseCommand) Execute(conn *connection.ConnInfo, data []byte, sender *response.ResponseSender, ctx context.Context) error {
	var resp proto2.RegisterResponse
	if err := proto.Unmarshal(data, &resp); err != nil {
		log.Println("Failed to unmarshal RegisterResponse:", err)
		return err
	}
	log.Printf("Received RegisterResponse: Success = %v, Message = %s", resp.Success, resp.Message)
	// 这里可以添加处理注册响应的逻辑
	return nil

}
