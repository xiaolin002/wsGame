package user

import (
	"wsprotGame/internal/api/user/login"
	"wsprotGame/internal/api/user/register"
	"wsprotGame/internal/service"
	proto2 "wsprotGame/proto/gen"
	"wsprotGame/server/command"
)

// WithUserCommands 用于初始化用户业务命令的选项函数

func WithUserCommands(userService service.UserService) command.CommandOption {
	return func(cmdMap map[proto2.GameMessage_MessageType]command.Command) {
		// 添加登录请求命令
		loginReqCmd := login.NewLoginRequestCommand(userService)
		cmdMap[proto2.GameMessage_LOGIN_REQUEST] = loginReqCmd

		// 添加注册请求命令
		registerCmd := register.NewRegisterRequestCommand(userService)
		cmdMap[proto2.GameMessage_REGISTER_REQUEST] = registerCmd

		loginRespCmd := &login.LoginResponseCommand{}
		cmdMap[proto2.GameMessage_LOGIN_RESPONSE] = loginRespCmd

		//添加注册响应命令
		registerRespCmd := &register.RegisterResponseCommand{}
		cmdMap[proto2.GameMessage_REGISTER_RESPONSE] = registerRespCmd
	}
}
