FROM golang:alpine AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct \
    GO111MODULE=on

# 移动到工作目录：/build
WORKDIR /build

# 将代码复制到容器中
COPY . .

# 下载依赖 将我们的代码编译成二进制可执行文件 bluebell_app
RUN go mod download && \
    go build -o blog .

###################
# 接下来创建一个小镜像
###################
FROM scratch

COPY ./config /config
COPY --from=builder /build/blog /

# 声明服务端口
EXPOSE 8080

# 需要运行的命令
ENTRYPOINT [ "/blog","config/config.yaml"]

