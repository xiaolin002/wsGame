package connection

import (
	"github.com/gorilla/websocket"
	"time"
)

// ConnInfo 连接包装结构体，包含连接及其属性
type ConnInfo struct {
	Conn *websocket.Conn
	// 可以根据需要添加更多属性
	CID         uint32
	Status      string    // 连接状态
	UserID      string    // 用户ID
	ConnectTime time.Time // 连接时间
}

// SendMessage 向该连接发送消息
func (ci *ConnInfo) SendMessage(messageType int, data []byte) error {
	if ci.Conn == nil {
		return nil
	}
	return ci.Conn.WriteMessage(messageType, data)
}

// Close 关闭该连接
func (ci *ConnInfo) Close() error {
	if ci.Conn == nil {
		return nil
	}
	return ci.Conn.Close()
}
