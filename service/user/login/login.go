package login

import (
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"log"
	"sync"
	"wsprotGame/command/response"
	proto2 "wsprotGame/proto/gen"
)

/**
 * @Description
 * @Date 2025/4/20 16:22
 **/
var users = make(map[string]string)
var usersMutex sync.Mutex

// LoginRequestCommand 登录请求命令
type LoginRequestCommand struct{}

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
	// 这里可以添加具体的登录处理逻辑
	usersMutex.Lock()
	defer usersMutex.Unlock()
	users["testuser"] = "testpassword"

	password, exists := users[req.Username]
	if !exists || password != req.Password {
		sender.Send(conn, proto2.GameMessage_LOGIN_RESPONSE, &proto2.LoginResponse{
			Success: false,
			Message: "用户名或密码错误",
		})
		return
	}

	// 假设有用 authResponseData使处理完其他业务逻辑返回的数据

	authResponseData := []byte{123}
	// 直接使用 LoginResponseCommand 处理响应
	loginResponseCmd := &LoginResponseCommand{}
	loginResponseCmd.Execute(conn, authResponseData, sender)

}

// LoginResponseCommand 登录响应命令
type LoginResponseCommand struct{}

func (c *LoginResponseCommand) Execute(conn *websocket.Conn, data []byte, sender *response.ResponseSender) {
	//var resp proto2.LoginResponse
	//if err := proto.Unmarshal(data, &resp); err != nil {
	//	log.Println("Failed to unmarshal LoginResponse:", err)
	//	return
	//}
	//log.Printf("Received LoginResponse: Success = %v, Message = %s, Token = %s", resp.Success, resp.Message, resp.Token)
	// 这里可以添加处理登录响应的逻辑

	sender.Send(conn, proto2.GameMessage_LOGIN_RESPONSE, &proto2.LoginResponse{
		Success: true,
		Message: "登录成功",
		Token:   "mock_token", // 实际应用中应生成安全的令牌
	})
}
