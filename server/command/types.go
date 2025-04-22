package command

import (
	"github.com/gorilla/websocket"
	proto2 "wsprotGame/proto/gen"
	"wsprotGame/server/command/response"
)

// Command 命令接口
type Command interface {
	Execute(conn *websocket.Conn, data []byte, sender *response.ResponseSender)
}

// CommandOption 选项函数类型
type CommandOption func(map[proto2.GameMessage_MessageType]Command)
