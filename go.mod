module wsprotGame

go 1.23.6

require (
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.3
	github.com/redis/go-redis/v9 v9.8.0
	google.golang.org/protobuf v1.28.0 // 替换为兼容的安全版本
	gorm.io/driver/mysql v1.5.7
	gorm.io/gorm v1.26.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-sql-driver/mysql v1.7.0
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/text v0.25.0 // indirect; 替换为兼容版本
)

require golang.org/x/crypto v0.38.0
