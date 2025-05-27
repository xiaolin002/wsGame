package connection

import (
	"context"
	"github.com/gorilla/websocket"
	"time"
)

// ConnInfo 连接包装结构体，包含连接及其属性
type ConnInfo struct {
	Conn *websocket.Conn
	// 可以根据需要添加更多属性
	Status        string          // 连接状态
	Uid           uint64          // 用户ID
	UserName      string          // 用户名称
	ConnectTime   time.Time       // 连接时间
	Ctx           context.Context // 存储 context 数据
	Authenticated bool            // 用户是否认证
	Cid           uint32          // 连接ID

}

func NewConnInfo(conn *websocket.Conn, ctx context.Context) *ConnInfo {
	return &ConnInfo{
		Conn:        conn,
		Status:      "connected",
		ConnectTime: time.Now(),
		Ctx:         ctx,
	}

}

func (ci *ConnInfo) SetAuthenticated(authenticated bool) {
	ci.Authenticated = authenticated
}
func (ci *ConnInfo) SetCid(cid uint32) {
	ci.Cid = cid
}

func (ci *ConnInfo) SetUid(uid uint64) {
	ci.Uid = uid
}

func (ci *ConnInfo) SetUserName(name string) {
	ci.UserName = name

}
func (ci *ConnInfo) SetStatus(status string) {
	ci.Status = status
}

// SendMessage 向该连接发送消息
func (ci *ConnInfo) SendMessage(messageType int, data []byte) error {
	if ci.Conn == nil {
		return nil
	}
	return ci.Conn.WriteMessage(messageType, data)
}
func (ci *ConnInfo) GetFormattedConnectTime() string {
	return ci.ConnectTime.Format("2006-01-02 15:04:05")
}

// Close 关闭该连接
func (ci *ConnInfo) Close() error {
	if ci.Conn == nil {
		return nil
	}
	return ci.Conn.Close()
}

func (ci *ConnInfo) GetStatus() string {
	return ci.Status
}

func (ci *ConnInfo) GetUserID() uint64 {
	return ci.Uid
}

func (ci *ConnInfo) GetConnectTime() time.Time {
	return ci.ConnectTime
}

func (ci *ConnInfo) GetUserName() string {
	return ci.UserName
}
