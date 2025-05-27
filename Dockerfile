# 使用官方 Go 语言镜像作为构建阶段的基础镜像
FROM golang:1.21.11 AS builder

# 设置工作目录
WORKDIR /app

# 将服务器 /file 目录下的 go.mod 和 go.sum 复制到工作目录
COPY /wsprotGame/go.mod /wsprotGame/go.sum ./

# ... 已有代码 ...

# 设置 Go 模块代理
ENV GOPROXY=https://goproxy.cn,direct
# 下载依赖
RUN go mod download

# 设置 Go 模块校验和数据库
ENV GOSUMDB=off

# 下载依赖，设置超时时间为 300 秒
RUN timeout 300 go mod download

# 整理依赖
RUN go mod tidy

# 将服务器 /file 目录下的源代码复制到工作目录
COPY ./wsprotGame/ .

# 打印当前目录内容，方便调试
RUN ls -la

# 打印 Go 环境信息
RUN go env

# ... 已有代码 ...

# 编译 Go 应用，添加 -v 参数显示编译过程，添加 -x 参数显示执行的命令
RUN go env -w GOOS=linux
RUN go env -w GOARCH=amd64
# 修改 go build 命令
RUN CGO_ENABLED=0 go build -o main .

# 使用轻量级的镜像来运行应用
FROM alpine:latest

# 设置工作目录
WORKDIR /root/

# 从 builder 镜像中复制编译好的二进制文件
COPY --from=builder /app/main .

# 如果有其他需要的文件（如配置文件、静态资源等），可以在这里复制
# COPY --from=builder /app/config ./config

# 暴露应用运行的端口
EXPOSE 8080

# 运行应用
CMD ["./main"]