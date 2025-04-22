package main

import (
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
	"wsprotGame/internal/api/user"
	"wsprotGame/internal/repository"
	"wsprotGame/internal/repository/dao"
	"wsprotGame/internal/service"
	"wsprotGame/ioc"
	proto2 "wsprotGame/proto/gen"
	"wsprotGame/server/command"
	"wsprotGame/server/command/register"
	"wsprotGame/server/command/response"
	"wsprotGame/server/connection"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func handleWebSocket(registry *regis.CommandRegistry, connManager *connection.ConnectionManager, w http.ResponseWriter, r *http.Request) {
	// 添加跨域头
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()
	// 添加连接到管理器

	// 创建连接信息实例
	connInfo := &connection.ConnInfo{
		Conn:        conn,
		Status:      "connected",
		ConnectTime: time.Now(),
	}

	// 添加连接到管理器，自动分配 cid
	cid := connManager.Add(connInfo)
	defer connManager.Remove(cid)
	// 心跳机制
	go func() {
		for {
			time.Sleep(30 * time.Second)
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("Failed to send ping:", err)
				break
			}
		}
	}()

	// 创建 ResponseSender 实例
	sender := &response.ResponseSender{}

	for {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second)) // 设置读取超时
		_, raw, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		// 这里打印出来的是原始的二进制数据
		//log.Println("Received message:", raw) // 添加日志

		var gameMsg proto2.GameMessage
		if err := proto.Unmarshal(raw, &gameMsg); err != nil {
			log.Println("Failed to unmarshal message:", err)
			continue
		}

		// 执行对应的命令
		registry.Handler(conn, gameMsg.Type, gameMsg.Data, sender)
	}
}

func main() {

	db := ioc.InitDB()
	userService := InitUserComponents(db)

	// 初始化命令注册表，传入用户业务的选项函数
	registry := InitCommandRegistry(
		user.WithUserCommands(userService),
		// 可以添加其他业务的选项函数
		// otherpackage.WithOtherCommands(otherService),
	)

	// 创建连接管理器
	connManager := connection.NewConnectionManager()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(registry, connManager, w, r)
	})
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// InitUserComponents 初始化用户相关的数据业务组件
func InitUserComponents(db *gorm.DB) service.UserService {
	gromDao := dao.NewUserGromDao(db)
	cacheRepository := repository.NewUserCacheRepository(gromDao)
	userService := service.NewUserService(cacheRepository)
	return userService
}

// InitCommandRegistry 初始化命令映射并注册命令
func InitCommandRegistry(options ...command.CommandOption) *regis.CommandRegistry {
	cmdMap := make(map[proto2.GameMessage_MessageType]command.Command)

	// 应用所有选项函数
	for _, opt := range options {
		opt(cmdMap)
	}

	registry := regis.NewCommandRegistry()
	for msgType, cmd := range cmdMap {
		registry.Register(msgType, cmd)
	}
	return registry
}
