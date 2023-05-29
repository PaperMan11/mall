FROM golang:alpine AS builder 
# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOPROXY=https://go.proxy.kylinsec.net
ENV GOPRIVATE=*.kylinsec.net

# 移动到工作目录
WORKDIR /build

# 复制项目并下载依赖信息
COPY . .
RUN go mod tidy

# 将我们的代码编译成二进制可执行文件
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o mall .


# 分阶段构建
FROM scratch
COPY ./config /config
# 从builder镜像中拷贝到当前目录
COPY --from=builder /build/mall /
# 声明端口
EXPOSE 8000
# 需要运行的命令
ENTRYPOINT ["/mall"]
CMD ["-f", "/config/config.yml"]