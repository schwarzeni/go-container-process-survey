pstree -apl | grep 'sh' -A3 -B3

[文件句柄何时关闭？](./container/log.go)

exec 执行 sh 会报错（影响不大）：

目前只能使用 main.go 中的命令，因为处理传入的 arg 比较麻烦，而且这个代码的重点也不是在这里

目前也不支持开机自动启动容器，也没有隔离文件系统，也没有使用镜像

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
