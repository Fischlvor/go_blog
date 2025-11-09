# 第一阶段：构建阶段
FROM golang:1.23.0-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的系统依赖
RUN apk add --no-cache git ca-certificates tzdata

# 设置Go代理（加速下载）
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.org

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 设置环境变量
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# 构建应用
RUN go build -a -installsuffix cgo -o main .

# 第二阶段：运行阶段
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 创建必要的目录
RUN mkdir -p /app/configs /app/log

# 复制配置文件
COPY --from=builder /app/configs/config.yaml /app/configs/

# 复制密钥文件
COPY --from=builder /app/keys /app/keys

# 暴露端口
EXPOSE 8001

# 启动应用
CMD ["./main"]

