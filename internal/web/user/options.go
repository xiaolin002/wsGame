package user

import (
	"wsprotGame/internal/service"
	"wsprotGame/internal/web/user/login"
	"wsprotGame/internal/web/user/register"
	proto2 "wsprotGame/proto/gen"
	"wsprotGame/server/command"
	"wsprotGame/server/connection"
)

// WithUserCommands 用于初始化用户业务命令的选项函数

func WithUserCommands(userService service.UserService, connManager *connection.ConnectionManager) command.CommandOption {
	return func(cmdMap map[proto2.GameMessage_MessageType]command.Command) {
		// 添加登录请求命令
		loginReqCmd := login.NewLoginRequestCommand(userService, connManager)
		cmdMap[proto2.GameMessage_LOGIN_REQUEST] = loginReqCmd

		// 添加注册请求命令
		registerCmd := register.NewRegisterRequestCommand(userService)
		cmdMap[proto2.GameMessage_REGISTER_REQUEST] = registerCmd

		// 添加登录返回命令
		loginRespCmd := &login.LoginResponseCommand{}
		cmdMap[proto2.GameMessage_LOGIN_RESPONSE] = loginRespCmd

		//添加注册响应命令
		registerRespCmd := &register.RegisterResponseCommand{}
		cmdMap[proto2.GameMessage_REGISTER_RESPONSE] = registerRespCmd
	}
}
