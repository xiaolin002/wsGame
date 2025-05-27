package private

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	proto2 "wsprotGame/proto/gen"
	"wsprotGame/server/command/response"
	"wsprotGame/server/connection"
)

type PrivateChatCommand struct {
	ConnManager *connection.ConnectionManager
}

func NewPrivateChatCommand(connManager *connection.ConnectionManager) *PrivateChatCommand {
	return &PrivateChatCommand{ConnManager: connManager}
}

func (pcc *PrivateChatCommand) Execute(connInfo *connection.ConnInfo, data []byte, sender *response.ResponseSender, ctx context.Context) error {
	var privateMsg proto2.PrivateChatMessage
	if err := proto.Unmarshal(data, &privateMsg); err != nil {
		log.Println("Failed to unmarshal private chat message:", err)
		return err
	}
	// 发送私聊消息给指定用户
	pcc.ConnManager.SendPrivateMessage(privateMsg, connInfo.UserName)
	return nil
}
