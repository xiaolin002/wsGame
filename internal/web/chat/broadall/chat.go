package broadall

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	proto2 "wsprotGame/proto/gen"
	"wsprotGame/server/command/response"
	"wsprotGame/server/connection"
)

/**
 * @Description
 * @Date 2025/5/26 11:23
 **/
// ChatCommand 处理聊天消息的命令

type ChatCommand struct {
	ConnManager *connection.ConnectionManager
}

func NewChatCommand(connManager *connection.ConnectionManager) *ChatCommand {
	return &ChatCommand{ConnManager: connManager}
}

// Execute 执行聊天消息处理逻辑
func (cc *ChatCommand) Execute(connInfo *connection.ConnInfo, data []byte, sender *response.ResponseSender, ctx context.Context) error {
	// 将data 反序列化为 proto2.ChatMessage
	var chatMsg proto2.ChatMessage
	if err := proto.Unmarshal(data, &chatMsg); err != nil {
		log.Println("Failed to unmarshal chat message:", err)
		return err
	}
	// 广播聊天消息给所有连接的用户
	cc.ConnManager.Broadcast(chatMsg, connInfo.UserName)
	return nil
}
