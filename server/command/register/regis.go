package regis

import (
	"context"
	"log"
	proto2 "wsprotGame/proto/gen"
	"wsprotGame/server/command"
	"wsprotGame/server/command/response"
	"wsprotGame/server/connection"
)

// CommandRegistry 命令注册表
type CommandRegistry struct {
	commands map[proto2.GameMessage_MessageType]command.Command
}

// NewCommandRegistry 创建新的命令注册表
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[proto2.GameMessage_MessageType]command.Command),
	}
}

// Register 注册命令
func (cr *CommandRegistry) Register(messageType proto2.GameMessage_MessageType, cmd command.Command) {
	cr.commands[messageType] = cmd
}

func (cr *CommandRegistry) Handler(conn *connection.ConnInfo, messageType proto2.GameMessage_MessageType, data []byte, sender *response.ResponseSender, ctx context.Context) {
	if cmd, exists := cr.commands[messageType]; exists {
		err := cmd.Execute(conn, data, sender, ctx)
		if err != nil {
			log.Printf("Failed to execute command for message type %v: %v", messageType, err)
		}
	} else {
		log.Printf("Unknown message type: %v", messageType)
	}
}
