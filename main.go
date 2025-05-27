package main

import (
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
	"wsprotGame/internal/repository"
	"wsprotGame/internal/repository/dao"
	"wsprotGame/internal/service"
	"wsprotGame/internal/web/chat"
	"wsprotGame/internal/web/user"
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

	//
	ctx := r.Context()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	// 创建连接信息实例
	connInfo := connection.NewConnInfo(conn, ctx)
	// 添加连接到管理器，自动分配 cid
	cid := connManager.Add(connInfo)

	defer conn.Close()
	defer connManager.Remove(cid)
	// 创建 ResponseSender 实
	sender := &response.ResponseSender{}
	// 心跳机制
	go func() {
		for {
			time.Sleep(30 * time.Second)
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("心跳机制ping失败:", err)
				break
			}
		}
	}()

	for {
		conn.SetReadDeadline(time.Now().Add(60 * time.Minute)) // 设置读取超时
		_, raw, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected close error: %v", err)
			} else if closeErr, ok := err.(*websocket.CloseError); ok {
				log.Printf("Connection closed with code %d and message %q", closeErr.Code, closeErr.Text)
			} else if err.Error() == "EOF" {
				log.Println("Connection closed: EOF")
			} else {
				log.Printf("Error reading message: %v", err)
			}
			// 显式关闭连接
			if closeErr := conn.Close(); closeErr != nil {
				log.Println("Error closing connection:", closeErr)
			}
			break
		}

		var gameMsg proto2.GameMessage
		if er := proto.Unmarshal(raw, &gameMsg); er != nil {
			log.Println("Failed to unmarshal message:", err)
			continue
		}
		// 如果用户未认证，只处理注册和登录请求
		if !connInfo.Authenticated && gameMsg.Type != proto2.GameMessage_REGISTER_REQUEST && gameMsg.Type != proto2.GameMessage_LOGIN_REQUEST {
			sender.Send(connInfo.Conn, proto2.GameMessage_LOGIN_RESPONSE, &proto2.LoginResponse{
				Success: false,
				Message: "请先进行注册或登录",
			})
			continue
		}
		// 启动协程处理消息
		go func(localGameMsg proto2.GameMessage, localConnInfo *connection.ConnInfo) {
			registry.Handler(localConnInfo, localGameMsg.Type, localGameMsg.Data, sender, ctx)
		}(gameMsg, connInfo)
	}

	// 执行对应的命令
	//registry.Handler(connInfo, gameMsg.Type, gameMsg.Data, sender, ctx)

}

func main() {

	db := ioc.InitDB()
	rdb := ioc.InitRedis()
	// 初始化用户业务组件
	userService := InitUserComponents(db, rdb)
	// 创建连接管理器
	connManager := connection.NewConnectionManager()

	// 初始化命令注册表，传入用户业务的选项函数
	registry := InitCommandRegistry(
		user.WithUserCommands(userService, connManager),
		chat.WithChatCommands(connManager),
	)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(registry, connManager, w, r)
	})
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// InitUserComponents 初始化用户相关的数据业务组件
func InitUserComponents(db *gorm.DB, rd redis.Cmdable) service.UserService {
	gromDao := dao.NewUserGromDao(db, rd)
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
	// 这里是将cmdMao中的命令注册到CommandRegistry 命令注册表
	registry := regis.NewCommandRegistry()
	for msgType, cmd := range cmdMap {
		registry.Register(msgType, cmd)
	}
	return registry
}
