FROM golang:1.21.5 AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

## 为我们的镜像设置必要的环境变量
#ENV GO111MODULE=on \
#    GOPROXY=https://goproxy.cn,direct \
#    CGO_ENABLED=0 \
#    GOOS=darwin \
#    GOARCH=arm64

# RUN apk update --no-cache && apk add --no-cache tzdata

# 移动到工作目录：/build
WORKDIR /www

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .

# 将我们的代码编译成二进制可执行文件 /.app
RUN go build -o app .

# 测试一下
#RUN echo `pwd` && ls -al && sleep 10s

###################
# 接下来创建一个小镜像
###################
FROM alpine:3.19.0
#FROM alpine:latest

COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /www

# 从builder镜像中把静态文件拷贝到当前目录
COPY ./conf /www/conf
COPY ./storage /www/storage

# 从builder镜像中把/拷贝到当前目录
COPY --from=builder /www/app .

# 暴露应用端口 需要和config.yaml中的port一致
EXPOSE 8200

# 需要运行的命令
ENTRYPOINT ["/www/app", "/www/conf/config.yaml"]