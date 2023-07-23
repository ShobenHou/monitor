# 使用Golang基础镜像
FROM golang:1.19.2 AS builder

# 设置工作目录
WORKDIR /go/src/app

# 拷贝项目代码到镜像中
COPY . .

#依赖下载
RUN go mod init github.com/ShobenHou/monitor
RUN go mod tidy
RUN go mod download

# ENV CGO_ENABLED=0
# 编译Golang代码
RUN go build -o agent /go/src/app/cmd/agent/main.go

# 设置启动命令
CMD ["/go/src/app/agent"]