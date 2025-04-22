package command

import (
	"github.com/gorilla/websocket"
	"wsprotGame/command/response"
	proto2 "wsprotGame/proto/gen"
	"wsprotGame/service/user/login"
	"wsprotGame/service/user/register"
)

// Command 命令接口
type Command interface {
	Execute(conn *websocket.Conn, data []byte, sender *response.ResponseSender)
}

var CommandMap = map[proto2.GameMessage_MessageType]Command{
	proto2.GameMessage_REGISTER_REQUEST:  &register.RegisterRequestCommand{},
	proto2.GameMessage_REGISTER_RESPONSE: &register.RegisterResponseCommand{},
	proto2.GameMessage_LOGIN_REQUEST:     &login.LoginRequestCommand{},
	proto2.GameMessage_LOGIN_RESPONSE:    &login.LoginResponseCommand{},
	// 后续新增的协议命令可以继续添加到这个映射表中
	// proto2.GameMessage_NEW_COMMAND: &package.NewCommand{},

}
