package connection

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"sync"
	"time"
	proto2 "wsprotGame/proto/gen"
)

// ConnectionManager 连接管理结构体

type ConnectionManager struct {
	connections   map[uint32]*ConnInfo // 使用自增作为键
	mutex         sync.Mutex
	userNameToCID map[string]uint32 // 新增：用户名到 cid 的映射
}

// / NewConnectionManager 创建一个新的连接管理器

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections:   make(map[uint32]*ConnInfo),
		userNameToCID: make(map[string]uint32),
	}
}

// Add 添加一个新的连接
func (cm *ConnectionManager) Add(connInfo *ConnInfo) uint32 {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	// 分配当前的 cid
	cid := generateUniqueID() // 生成唯一ID
	cm.connections[cid] = connInfo
	connInfo.SetCid(cid)
	// 获取远程地址并转换为 *net.TCPAddr 类型
	tcpAddr, ok := connInfo.Conn.RemoteAddr().(*net.TCPAddr)
	if !ok {
		log.Println("Failed to convert RemoteAddr to *net.TCPAddr")
	}

	log.Printf("Client connected ,connection CID  is: %d, remote Addr is %s", cid, tcpAddr.IP.String())
	// 广播用户上线消息

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
		// 若连接信息包含用户名，从映射中删除记录
		if connInfo.UserName != "" {
			delete(cm.userNameToCID, connInfo.UserName)
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
func (cm *ConnectionManager) Broadcast(chatMsg proto2.ChatMessage, sender string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	// 安全验证和统一管理 所以重新赋值
	chatMsg.Sender = sender
	chatMsg.Timestamp = time.Now().Unix()

	chatData, err := proto.Marshal(&chatMsg)
	if err != nil {
		log.Println("Failed to marshal chat message:", err)
		return
	}

	gameMsg := proto2.GameMessage{
		Type: proto2.GameMessage_CHAT_MESSAGE,
		Data: chatData,
	}

	raw, err := proto.Marshal(&gameMsg)
	if err != nil {
		log.Println("Failed to marshal game message:", err)
		return
	}

	for _, connInfo := range cm.connections {
		if connInfo.Authenticated {
			if err := connInfo.Conn.WriteMessage(websocket.BinaryMessage, raw); err != nil {
				log.Println("Failed to send chat message:", err)
			}
		}
	}
}

func (cm *ConnectionManager) SendPrivateMessage(privateMsg proto2.PrivateChatMessage, sender string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	// 消息发送人 安全性和统一管理
	privateMsg.Sender = sender
	privateMsg.Timestamp = time.Now().Unix()
	chatData, err := proto.Marshal(&privateMsg)
	if err != nil {
		log.Println("Failed to marshal chat message:", err)
		return
	}

	gameMsg := proto2.GameMessage{
		Type: proto2.GameMessage_PRIVATE_CHAT_MESSAGE,
		Data: chatData,
	}
	raw, err := proto.Marshal(&gameMsg)
	if err != nil {
		log.Println("Failed to marshal game message:", err)
		return
	}

	// 获取收件人信息
	receiver := privateMsg.Receiver
	// 查找接收者的连接
	// 查找接收者的 cid
	if cid, ok := cm.userNameToCID[receiver]; ok {
		if connInfo, exists := cm.connections[cid]; exists && connInfo.Authenticated {
			if err := connInfo.Conn.WriteMessage(websocket.BinaryMessage, raw); err != nil {
				log.Printf("Failed to send private chat message to %s: %v", receiver, err)
				if closeErr := connInfo.Conn.Close(); closeErr != nil {
					log.Printf("Failed to close connection for %s: %v", receiver, closeErr)
				}
				// 移除断开的连接
				cm.Remove(cid)
			}
		}
	} else {
		log.Printf("User %s not found", receiver)
	}
}

// UpdateUserName 在用户登录成功后更新用户名到 cid 的映射
func (cm *ConnectionManager) UpdateUserName(cid uint32, newUserName string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if connInfo, ok := cm.connections[cid]; ok {
		// 如果之前有用户名，先从映射中删除
		if oldUserName := connInfo.UserName; oldUserName != "" {
			delete(cm.userNameToCID, oldUserName)
		}
		// 更新连接信息中的用户名
		connInfo.UserName = newUserName
		// 将新的用户名和 cid 存入映射
		cm.userNameToCID[newUserName] = cid
	}
	return nil
}

// GetOnlinePlayers 获取所有在线玩家的用户名
func (cm *ConnectionManager) GetOnlinePlayers() []string {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	var playerNames []string
	for _, connInfo := range cm.connections {
		if connInfo.Authenticated {
			playerNames = append(playerNames, connInfo.UserName)
		}
	}
	return playerNames
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
