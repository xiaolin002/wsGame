package regis

import (
	"github.com/gorilla/websocket"
	"log"
	"wsprotGame/command"
	"wsprotGame/command/response"
	proto2 "wsprotGame/proto/gen"
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

func (cr *CommandRegistry) Handler(conn *websocket.Conn, messageType proto2.GameMessage_MessageType, data []byte, sender *response.ResponseSender) {
	if cmd, exists := cr.commands[messageType]; exists {
		cmd.Execute(conn, data, sender)
	} else {
		log.Printf("Unknown message type: %v", messageType)
	}
}
