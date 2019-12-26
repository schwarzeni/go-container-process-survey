玩具项目，0分健壮性，0分测试覆盖


pstree -apl | grep 'sh' -A3 -B3

[文件句柄何时关闭？](./container/log.go)

exec 执行 sh 会报错（影响不大）：

目前只能使用 main.go 中的命令，因为处理传入的 arg 比较麻烦，而且这个代码的重点也不是在这里

目前也不支持开机自动启动容器，也没有隔离文件系统，也没有使用镜像

同时对传入程序的 arg 没有做验证，默认合法

## v0.1 功能

```bash
# 前台运行
go run . -name "sh-1"

# 后台运行
go run . -name "sh-1" -d

# 查看容器运行状态
go run . log

# 停止容器
go run . stop <id>

# 启动停止容器
go run . start <id>

# 登入运行的容器
go run . exec <id> bash
```

## v0.2 功能

```bash
# 提供独立的文件系统，注意，不能是压缩包
go run . -name "sh-1" -d -image "/path/to/fs"

# 支持用户自己挂载 volume，
# 由于 golang 原生 flag 库局限性，只能挂载一个目录
# src 为外部需要挂载的
# dst 为容器内部的路径，挂载目的地
go run . -name "sh-1" -d -image "/path/to/fs" -v "src:dst"

# 导出正在运行的容器，注意，path为导出文件的完整路径，比如 /tmp/image.tar
go run . commit <container_id> <path>

# 添加环境变量
# 由于 golang 原生 flag 库局限性，只能传入一个环境变量
go run . -name "sh-1" -d -image "/path/to/fs" -e HELLO=WORLD
```
