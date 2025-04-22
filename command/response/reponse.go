package response

import (
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"log"
	proto2 "wsprotGame/proto/gen"
)

// ResponseSender 用于发送响应的结构体
type ResponseSender struct{}

// Send 发送响应消息
func (rs *ResponseSender) Send(conn *websocket.Conn, msgType proto2.GameMessage_MessageType, response proto.Message) {
	//log.Printf("发送响应消息，消息类型: %v", msgType)

	data, err := proto.Marshal(response)
	if err != nil {
		log.Println("Failed to marshal response:", err)
		return
	}

	gameMsg := proto2.GameMessage{
		Type: msgType,
		Data: data,
	}

	raw, err := proto.Marshal(&gameMsg)
	if err != nil {
		log.Println("Failed to marshal game message:", err)
		return
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, raw); err != nil {
		log.Println("Failed to send response:", err)
	}
}
