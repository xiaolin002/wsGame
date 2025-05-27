package command

import (
	"context"
	proto2 "wsprotGame/proto/gen"
	"wsprotGame/server/command/response"
	"wsprotGame/server/connection"
)

// Command 命令接口
type Command interface {
	// Execute 每个具体的命令都需要实现这个接口，用于执行具体的业务逻辑。
	Execute(conn *connection.ConnInfo, data []byte, sender *response.ResponseSender, ctx context.Context) error
}

// CommandOption 选项函数类型
type CommandOption func(map[proto2.GameMessage_MessageType]Command)
