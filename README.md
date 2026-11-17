# Bluebell

Go / Gin 社区项目：注册登录、社区、帖子、按时间或综合分数排序，以及 Redis 投票。

## 本次更新

- 修复注册页确认密码字段，提交 `re_password`。
- 详情页显示 Redis 中的真实赞成票数，移除固定的 `50.2k`。
- 列表与详情页投票成功后刷新数据，显示重复投票、投票过期等错误。
- 保存 Vue 源码、最新部署资源、Swagger 文档和部署文件。

## 本地运行

`go.mod` 声明 Go 1.27.1。准备 MySQL 和 Redis，复制配置示例后填写实际连接参数和私有密钥：

```powershell
Copy-Item conf/config.local.example.yaml conf/config.local.yaml
go run . .\conf\config.local.yaml
```

从项目根目录启动。首页为 `http://127.0.0.1:8080/`，Swagger 为 `http://127.0.0.1:8080/swagger/index.html`。

`init.sql` 用于初始化全新的数据库，其中包含 DROP TABLE，不要在已有业务数据库上重复执行。

实际配置、日志、数据库目录、IDE 文件、编译产物和前端依赖不提交到 Git。

## 前端构建

源码在 `frontend/`，Go 或 Nginx 使用根目录的 `static/` 与 `templates/index.html`。

```powershell
cd frontend
npm ci --ignore-scripts --legacy-peer-deps --registry=https://registry.npmjs.org
$env:NODE_OPTIONS="--openssl-legacy-provider"
npm run build
Copy-Item dist/static -Destination .. -Recurse -Force
Copy-Item dist/index.html -Destination ../templates/index.html -Force
```

兼容参数用于现有 Vue CLI 4 / Webpack 4。上传时先上传 static，再更新 HTML，保留旧哈希资源以兼容已打开的页面。

## 阿里云：Supervisor + Nginx

在 Windows PowerShell 编译 Linux x86-64 程序：

```powershell
$env:GOOS="linux"
$env:GOARCH="amd64"
$env:CGO_ENABLED="0"
go build -trimpath -ldflags="-s -w" -o bin/bluebell-linux-amd64 .
Remove-Item Env:GOOS,Env:GOARCH,Env:CGO_ENABLED
```

服务器示例目录为 `/data/app/bluebell`，包含 `bin/bluebell`、`conf/config.aliyun.yaml`、`static/` 和 `templates/index.html`。

根据 `conf/config.aliyun.example.yaml` 创建服务器私有配置，保留线上数据库参数。Supervisor 配置：

```ini
directory=/data/app/bluebell
command=/data/app/bluebell/bin/bluebell /data/app/bluebell/conf/config.aliyun.yaml
```

后端先上传为 `bin/bluebell.new`，备份旧程序，用 Supervisor 停止后替换，赋予执行权限并启动。不要直接覆盖运行中的可执行文件。Nginx 配置未变时不需要重载。

Vue 使用 history 路由。Nginx 在现有 server 中直接提供前端页面时，可使用：

```nginx
location / {
    root /data/app/bluebell/templates;
    index index.html;
    try_files $uri $uri/ /index.html;
}
location /static/ {
    alias /data/app/bluebell/static/;
}
location /api/ {
    proxy_pass http://127.0.0.1:8080;
}
location /swagger/ {
    proxy_pass http://127.0.0.1:8080;
}
```

按实际域名、HTTPS、文件权限和原配置合并；不要覆盖完整 Nginx 配置。直接访问 Go 的 8080 端口不经过 Nginx 的页面回退。

## 已有部署的密钥兼容

JWT 签名密钥和旧版密码哈希的盐值已移出源码。使用此提交重新构建后端前，在服务器私有 YAML 中配置：

```yaml
auth:
  jwt_expire: 8760
  jwt_secret: "填写服务器原来的 JWT 签名密钥"
  password_salt: "填写服务器原来的密码哈希盐值"
```

也可分别通过 `BLUEBELL_JWT_SECRET`、`BLUEBELL_PASSWORD_SALT` 环境变量提供。已有账号必须保留原密码盐值，否则旧密码无法校验；保留原 JWT 密钥可兼容已签发令牌。真实值只保存在私有配置。

之前上传的编译程序不受本次源码变化影响；以后重新编译部署时才需要以上配置。

## Docker Compose

```powershell
Copy-Item conf/config.docker.example.yaml conf/config.docker.yaml
```

填写私有密钥，然后启动：

```text
docker compose up -d --build
docker compose ps
docker compose logs --tail=100 bluebell_app
```

应用端口为 8888，MySQL 为 33061，Redis 为 63791。容器配置使用 `mysql8019:3306`、`redis507:6379`。配置通过只读挂载提供，不写入镜像。示例数据库密码用于本地开发，正式环境需同步调整 Compose 和私有 YAML。

仅更新应用：

```text
docker compose up -d --build --no-deps bluebell_app
```

Compose 数据库不会自动包含现有宿主机数据。不要删除已有 mysql/data。Redis 保存排序和投票数据，正式部署还需配置持久化和备份。

## 排序和票数

- New：发布时间倒序。
- Score：按“发布时间戳 + 432 ×（赞成票数 − 反对票数）”倒序。
- vote_num：赞成票数，不是综合排名分。
- 发帖超过一周不允许继续投票。

## 验证

```text
go test ./...
go vet ./...
```

默认测试不连接数据库。MySQL 集成写入测试需要显式设置指向隔离测试库的 BLUEBELL_TEST_CONFIG，并执行 `go test -tags=integration ./dao/mysql`。该测试会写入数据，不能指向线上数据库。
