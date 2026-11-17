
#在 Go 程序启动之前，先检查 MySQL、Redis 等依赖服务是否已经启动并且可以连接。
# 如果还不能连接，就一直等待；直到连接成功，再启动 Go 程序。

#!/bin/sh

# 遇到错误立即退出
set -e

# 依次等待传入的 host:port
while [ "$1" != "--" ]; do
    host="${1%:*}"
    port="${1##*:}"

    echo "Waiting for $host:$port ..."

    # 使用 nc 检查端口是否可以连接
    until nc -z "$host" "$port"; do
        sleep 1
    done

    echo "$host:$port is ready!"

    shift
done

# 跳过 --
shift

# 执行后面的 Go 程序
exec "$@"
