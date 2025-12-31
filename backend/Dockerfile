# 多阶段构建 - 构建阶段
FROM golang:1.23-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache ca-certificates tzdata

# 设置Go代理和环境变量
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

# 复制所有源代码
COPY . .

# 构建应用程序（不需要单独下载依赖）
RUN go mod tidy && go build -a -installsuffix cgo -o stock-api .

# 运行阶段
FROM alpine:latest

# 安装ca-certificates（HTTPS请求需要）
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# 从构建阶段复制可执行文件
COPY --from=builder /app/stock-api .

# 暴露端口
EXPOSE 8080

# 运行应用程序
CMD ["./stock-api"]