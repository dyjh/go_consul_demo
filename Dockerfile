# 使用 golang:1.20.12 版本的 alpine 镜像作为基础镜像
FROM golang:1.20.12-alpine as builder




# 设置工作目录
WORKDIR /app

# 复制整个项目目录到工作目录
COPY . .

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go env -w CGO_ENABLED=0
RUN go env
RUN go mod tidy
# 下载依赖
RUN go mod download

# 复制服务端代码到工作目录
COPY main.go .

# 构建服务端应用
RUN go build -o server .

# 创建最终的镜像
FROM alpine:latest

# 安装依赖
RUN apk add --no-cache ca-certificates

# 设置工作目录
WORKDIR /app

# 从构建镜像中复制二进制文件到最终镜像
COPY --from=builder /app/server .

# 暴露服务端口
EXPOSE 1234

# 启动服务
CMD ["./server"]
