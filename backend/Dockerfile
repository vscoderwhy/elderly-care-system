# 后端 Dockerfile
FROM golang:1.21-alpine AS builder

# 安装必要的工具
RUN apk add --no-cache git make

# 设置工作目录
WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o elderly-care-api ./cmd/main.go

# 运行阶段
FROM alpine:latest

# 安装 ca 证书
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/elderly-care-api .
COPY --from=builder /app/config.yaml .

# 创建上传目录
RUN mkdir -p /uploads

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./elderly-care-api"]
