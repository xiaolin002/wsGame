package connection

import (
	"github.com/google/uuid"
	"log"
	"net"
	"sync"
)

// ConnectionManager 连接管理结构体

type ConnectionManager struct {
	connections map[uint32]*ConnInfo // 使用自增作为键
	mutex       sync.Mutex
}

// / NewConnectionManager 创建一个新的连接管理器

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[uint32]*ConnInfo),
	}
}

// Add 添加一个新的连接
func (cm *ConnectionManager) Add(connInfo *ConnInfo) uint32 {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	// 分配当前的 cid
	cid := generateUniqueID() // 生成唯一ID
	cm.connections[cid] = connInfo
	// 获取远程地址并转换为 *net.TCPAddr 类型
	tcpAddr, ok := connInfo.Conn.RemoteAddr().(*net.TCPAddr)
	if !ok {
		log.Println("Failed to convert RemoteAddr to *net.TCPAddr")
	}

	log.Println("Client connected") // 添加日志

	log.Printf("New connection added with CID: %d, remote Addr is %s", cid, tcpAddr.IP.String())
	return cid
}

// Remove 移除一个连接
func (cm *ConnectionManager) Remove(cid uint32) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	if connInfo, ok := cm.connections[cid]; ok {
		if err := connInfo.Close(); err != nil {
			log.Printf("Failed to close connection for user %d: %v", cid, err)
		}
		delete(cm.connections, cid)
	}
}

// Get 获取指定cid的连接信息
func (cm *ConnectionManager) Get(cid uint32) (*ConnInfo, bool) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	connInfo, ok := cm.connections[cid]
	return connInfo, ok
}

// Broadcast 向所有连接广播消息
func (cm *ConnectionManager) Broadcast(messageType int, message []byte) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	for _, connInfo := range cm.connections {
		if err := connInfo.SendMessage(messageType, message); err != nil {
			// 处理发送失败的情况，移除连接
			if err := connInfo.Close(); err != nil {
				log.Printf("Failed to close connection for user %d: %v", connInfo.CID, err)
			}
			delete(cm.connections, connInfo.CID)
		}
	}
}

// BroadcastToSpecificUsers 向指定用户列表广播消息
func (cm *ConnectionManager) BroadcastToSpecificUsers(cids []uint32, messageType int, message []byte) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	for _, cid := range cids {
		if connInfo, ok := cm.connections[cid]; ok {
			if err := connInfo.SendMessage(messageType, message); err != nil {
				if err := connInfo.Close(); err != nil {
					log.Printf("Failed to close connection for user %d: %v", cid, err)
				}
				delete(cm.connections, cid)
			}
		}
	}
}
func generateUniqueID() uint32 {
	id := uuid.New()
	// 将 UUID 的前 4 个字节转换为 uint32
	var result uint32
	for i := 0; i < 4; i++ {
		result = (result << 8) | uint32(id[i])
	}
	return result // 将 UUID 转换为 uint32
}
