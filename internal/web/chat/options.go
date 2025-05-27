package chat

import (
	"wsprotGame/internal/web/chat/broadall"
	"wsprotGame/internal/web/chat/private"
	proto2 "wsprotGame/proto/gen"
	"wsprotGame/server/command"
	"wsprotGame/server/connection"
)

/**
 * @Description
 * @Date 2025/5/26 11:24
 **/

// WithChatCommands 选项函数，用于注册聊天命令
func WithChatCommands(connManager *connection.ConnectionManager) command.CommandOption {
	return func(cmdMap map[proto2.GameMessage_MessageType]command.Command) {
		// 注册聊天命令
		gloabalChat := broadall.NewChatCommand(connManager)
		cmdMap[proto2.GameMessage_CHAT_MESSAGE] = gloabalChat
		// 注册私聊命令
		privateChatCommand := private.NewPrivateChatCommand(connManager)
		cmdMap[proto2.GameMessage_PRIVATE_CHAT_MESSAGE] = privateChatCommand

	}
}
