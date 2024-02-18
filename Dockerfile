FROM golang:alpine AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录：/build
WORKDIR /build

# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息
COPY go.mod .
COPY go.sum .

#设置代理
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod download

# 将代码复制到容器中
COPY . .

# 将我们的代码编译成二进制可执行文件 bluebell_app
RUN go build -o blog .

###################
# 接下来创建一个小镜像
###################
FROM debian:stretch-slim

COPY ./config /config
# 从builder镜像中把/dist/app 拷贝到当前目录
COPY --from=builder /build/blog /

# 声明服务端口
EXPOSE 8080
#这个不是真正暴露的端口

# 需要运行的命令
ENTRYPOINT [ "/blog","config/config.yaml"]
