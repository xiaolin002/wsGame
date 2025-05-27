package login

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

/**
 * @Description
 * @Date 2025/4/20 16:22
 **/

// LoginRequestCommand 登录请求命令
type LoginRequestCommand struct {
	us service.UserService
	cm *connection.ConnectionManager
}

func NewLoginRequestCommand(userService service.UserService, connManager *connection.ConnectionManager) *LoginRequestCommand {
	return &LoginRequestCommand{
		us: userService,
		cm: connManager,
	}
}

func (c *LoginRequestCommand) Execute(conn *connection.ConnInfo, data []byte, sender *response.ResponseSender, ctx context.Context) error {
	var req proto2.LoginRequest
	if err := proto.Unmarshal(data, &req); err != nil {
		sender.Send(conn.Conn, proto2.GameMessage_LOGIN_RESPONSE, &proto2.LoginResponse{
			Success: false,
			Message: "解析请求失败",
		})
		return err
	}
	log.Printf("Received LoginRequest: Username = %s, Password = %s", req.Username, req.Password)

	if err := c.us.HandleLogin(ctx, conn, domain.User{
		NickName: req.Username,
		Password: req.Password,
	}); err != nil {
		sender.Send(conn.Conn, proto2.GameMessage_LOGIN_RESPONSE, &proto2.LoginResponse{
			Success: false,
			Message: "账户或者密码错误",
		})
		return err
	}
	// 登录成功，设置认证状态
	conn.SetAuthenticated(true)
	conn.SetUserName(req.Username)
	conn.SetStatus("login")
	err := c.cm.UpdateUserName(conn.Cid, conn.UserName)
	if err != nil {
		log.Printf("更新用户状态失败：%v", err)
	}

	// 登录成功，发送成功响应
	sender.Send(conn.Conn, proto2.GameMessage_LOGIN_RESPONSE, &proto2.LoginResponse{
		Success: true,
		Message: "登录成功",
		Token:   "mock_token", // 实际应用中应生成安全的令牌
	})
	log.Printf("连接属性：状态is ：%s,用户名 ：%s, 连接时间 ：%s,用户uid ：%020d\n", conn.Status, conn.UserName, conn.GetFormattedConnectTime(), conn.Uid)
	return nil

}

// LoginResponseCommand 登录响应命令
type LoginResponseCommand struct{}

func (c *LoginResponseCommand) Execute(conn *connection.ConnInfo, data []byte, sender *response.ResponseSender, ctx context.Context) error {

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

	sender.Send(conn.Conn, proto2.GameMessage_LOGIN_RESPONSE, &proto2.LoginResponse{
		Success: true,
		Message: "登录成功",
		Token:   "mock_token", // 实际应用中应生成安全的令牌
	})
	return nil
}
