FROM golang:alpine AS builder
#AS builder：给当前构建阶段起名为 builder。

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录：/build
WORKDIR /build

# 复制项目中的 go.mod 和 go.sum文件并下载项目所需要的 Go 依赖包。
COPY go.mod .
COPY go.sum .
RUN go mod download

# 将当前 Docker 构建上下文中的项目文件复制到容器的 /build 目录
COPY . .

# 将我们的代码编译成二进制可执行文件 bubble_app,.：编译当前目录下的 Go 项目。
RUN go build -o bluebell_app .

###################
# 接下来创建一个小镜像
###################
FROM debian:bookworm-slim
#把项目中的 templates 文件夹复制到最终镜像的根目录下。
COPY ./wait-for.sh /
COPY ./templates /templates
COPY ./static /static
COPY ./conf /conf

# 从builder镜像中把/dist/app 拷贝到当前目录
COPY --from=builder /build/bluebell_app /

RUN set -eux; \
    sed -i 's|http://deb.debian.org|http://mirrors.tuna.tsinghua.edu.cn|g' /etc/apt/sources.list.d/debian.sources; \
    apt-get -o Acquire::Retries=5 update; \
    apt-get install -y --no-install-recommends ca-certificates netcat-openbsd; \
    chmod 755 /wait-for.sh; \
    rm -rf /var/lib/apt/lists/*

#声明服务端口,表示应用预计监听容器的 8080 端口。
EXPOSE 8080
#需要运行的命令
#ENTRYPOINT ["/bluebell_app","conf/config.yaml"]
