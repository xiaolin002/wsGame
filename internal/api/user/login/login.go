package login

import (
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"log"
	"sync"
	"wsprotGame/internal/service"
	proto2 "wsprotGame/proto/gen"
	"wsprotGame/server/command/response"
)

/**
 * @Description
 * @Date 2025/4/20 16:22
 **/
var users = make(map[string]string)
var usersMutex sync.Mutex

// LoginRequestCommand 登录请求命令
type LoginRequestCommand struct {
	us service.UserService
}

func NewLoginRequestCommand(userService service.UserService) *LoginRequestCommand {
	return &LoginRequestCommand{us: userService}
}

func (c *LoginRequestCommand) Execute(conn *websocket.Conn, data []byte, sender *response.ResponseSender) {
	var req proto2.LoginRequest
	if err := proto.Unmarshal(data, &req); err != nil {
		sender.Send(conn, proto2.GameMessage_LOGIN_RESPONSE, &proto2.LoginResponse{
			Success: false,
			Message: "解析请求失败",
		})
		return
	}
	log.Printf("发送登录·成功响应，消息类型: %v", proto2.GameMessage_LOGIN_RESPONSE)
	log.Printf("Received LoginRequest: Username = %s, Password = %s", req.Username, req.Password)

	/*
		 假设有用 authResponseData使处理完其他业务逻辑返回的数据
		authResponseData := []byte{123}
		从注册表获取 LoginResponseCommand 并执行
		if registry != nil { // 假设 registry 是全局的命令注册表
			if cmd, ok := registry.Get(proto2.GameMessage_LOGIN_RESPONSE); ok {
				if loginRespCmd, ok := cmd.(*LoginResponseCommand); ok {
					authResponseData := []byte{123} // 模拟处理完业务逻辑返回的数据
					loginRespCmd.Execute(conn, authResponseData, sender)
				}
			}
		}*/

	if err := c.us.HandleLogin(req); err != nil {
		sender.Send(conn, proto2.GameMessage_LOGIN_RESPONSE, &proto2.LoginResponse{
			Success: false,
			Message: "账户或者密码错误",
		})
	}
	// 登录成功，发送成功响应
	sender.Send(conn, proto2.GameMessage_LOGIN_RESPONSE, &proto2.LoginResponse{
		Success: true,
		Message: "登录成功",
		Token:   "mock_token", // 实际应用中应生成安全的令牌
	})
}

// LoginResponseCommand 登录响应命令
type LoginResponseCommand struct{}

func (c *LoginResponseCommand) Execute(conn *websocket.Conn, data []byte, sender *response.ResponseSender) {

	sender.Send(conn, proto2.GameMessage_LOGIN_RESPONSE, &proto2.LoginResponse{
		Success: true,
		Message: "登录成功",
		Token:   "mock_token", // 实际应用中应生成安全的令牌
	})
}
