.PHONY: all build run gotool clean help

# 定义生成的可执行文件名称
BINARY=bluebell.exe
CONFIG=conf/config.local.yaml

# 默认执行目标：先格式化和检查代码，再编译项目
all: gotool build

# 编译 Go 项目
# CGO_ENABLED=0：禁用 CGO
# GOOS=windows：目标操作系统为 Windows
# GOARCH=amd64：目标 CPU 架构为 64 位 x86
# -o：指定生成的可执行文件名称
build:
	set "CGO_ENABLED=0" && set "GOOS=linux" && set "GOARCH=amd64" && go build -ldflags="-s -w" -o "./bin/$(BINARY)"

# 直接运行当前目录下的 Go 项目
run:
	@go run ./ "$(CONFIG)"

# Go 代码格式化和检查
gotool:
	go fmt ./
	go vet ./

# 删除编译生成的可执行文件
clean:
	@if exist "$(BINARY)" del /Q "$(BINARY)"

# 显示帮助信息
help:
	@echo make - 格式化,检查并编译 Go 项目
	@echo make build - 编译 Go 项目,生成 Windows 可执行文件
	@echo make run - 直接运行 Go 项目
	@echo make clean - 删除编译生成的可执行文件
	@echo make gotool - 运行 Go 工具 fmt 和 vet

# Makefile 用于管理 Go 项目的常用操作：
# 1. 格式化 Go 代码
# 2. 检查代码问题
# 3. 编译 Go 项目
# 4. 运行 Go 项目
# 5. 删除编译生成的文件
