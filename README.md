# 《自己动手写Docker》第五章相关技术调研

玩具项目，0分代码健壮性，0分测试覆盖率，没有做 cgroup 资源限制

---

## 主要功能

### 支持后台运行

这个目前的简单实现是将容器进程挂到 `pid` 为1的进程下。这样的话就需要设置一下后台运行容器进程的输出，程序中将 `stdout` 和 `stderr` 定向到一个 log 文件中，这样确保输出不会丢失，如下：

```go
var stdLog *os.File
if stdLog, err = StdLog(containerID); err != nil {
  return fmt.Errorf("create log file failed: %v", err)
}
cmd.Stderr = stdLog
cmd.Stdout = stdLog
```

详细的内容在 [container/run.go](container/run.go) 中

---

### 查看运行容器信息

对容器进行操作时，在本地的文件中做相关的记录。那么首先需要定义一个记录容器信息的结构体，如下([container/container_info.go](./container/container_info.go))，将其存放在一个文件中

```go
// 存储容器相关信息的json结构
type containerInfo struct {
	Pid         string   `json:"pid"`         // 容器的init的进程在宿主机上的PID
	ID          string   `json:"id"`          // 容器ID
	Name        string   `json:"name"`        // 容器名称
	Command     string   `json:"command"`     // 容器内 init 进程的运行命令
	FullCommand []string `json:"fullCommand"` // command and args
	CreatedTime string   `json:"createTime"`  // 创建时间
	Status      string   `json:"status"`      // 容器状态
	ImageURL    string   `json:"imageURL"`    // 镜像的存储位置
	Volumes     []string `json:"volumes"`     // 挂载的数据卷
	Envs        []string `json:"envs"`        // 用户传入的环境变量
}
```

大概的文件目录结构如下：

```txt
|-- 0244200882
|   |-- aufs
|   |   |-- mnt
|   |   `-- writelayer
|   |-- config.json
|   `-- container.log
`-- 0291874589
    |-- aufs
    |   |-- mnt
    |   `-- writelayer
    |-- config.json
    `-- container.log
```

- 那些数字为每个容器进程特有的ID
- `aufs` 为容器挂载文件系统的相关信息
- `config.json` 就是记录容器相关信息的文件
- `container.log` 为记录容器进程相关输出的文件

而需要查看相关的信息，不论是容器的相关信息还是容器的输出，只需要读取这些文件即可

---

### 进入运行容器内部

书中提到，进入进程运行的命名空间中，需要执行系统调用 `setns` ，但是，其无法被多线程的程序调用。此时，就需要使用 CGO 技术，使用go语言调用c程序，c程序的大体结构如下([cgo/exec.h](./cgo/exec.h))

```c
__attribute__((constructor)) void enter_namespace(void)
{
  const char *PID_ENV = "mydocker_pid";
  const char *CMD_ENV = "mydocker_cmd";
  char *mydocker_pid = getenv(PID_ENV); // 进程需要进入的pid
  char *mydocker_cmd = getenv(CMD_ENV); // 需要执行的命令

  // check env
  if (!mydocker_pid)
  {
    return;
  }

  char nspath[1024];
  char nslist[][5] = {"ipc", "uts", "net", "pid", "mnt"}; // 需要进入的五种namespace

  for (size_t i = 0; i < 5; i++)
  {
    sprintf(nspath, "/proc/%s/ns/%s", mydocker_pid, nslist[i]);
    int fd = open(nspath, O_RDONLY);
    if (setns(fd, 0) == -1)
    {
      fprintf(stderr, "enter ns %s failed", nspath);
      exit(EXIT_FAILURE);
    }
    if (close(fd) == -1)
    {
      fprintf(stderr, "close fd for %s failed", nspath);
      exit(EXIT_FAILURE);
    }
  }

  int res = system(mydocker_cmd);
  exit(EXIT_SUCCESS);
}
```

其中，通过 gcc 的扩展  `__attribute__((constructor))` 来实现程序启动前执行特定代码。在go语言代码中，使用 `exec.Command("/proc/self/exe")` 来调用自身来使用以上的c代码，同时设置相关的环境变量，代码如下([container/exec.go](./container/exec.go))

```go
cmd = exec.Command("/proc/self/exe", "exec")
cmd.Stdin = os.Stdin
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
_ = os.Setenv("mydocker_pid", containerInfo.Pid)
_ = os.Setenv("mydocker_cmd", cmdStr)
```
---

### 容器的停止、再启动、删除

关于如何停止容器，书中使用的是 `syscall.Kill(pidInt, syscall.SIGTERM)` ，也就是向容器进程发送 `15` 信号，但是，对某些进程无效，所以这里就直接发送 `9` 信号，也就是 `syscall.Kill(pidInt, syscall.SIGKILL)`

启动容器和一开始的运行一个新容器非常类似，就是略去了一些步骤，比如创建 aufs 的可写层

删除容器就是把一个停止容器的相关目录删掉就行了，就是上面提到的那个目录结构中的目录

---

### 传入环境变量

对于启动新容器和启动停止的容器，这个很简单，直接做如下设置即可

```go
cmd.Env = append(os.Environ(), envs...)
```

但是对于进入运行容器的命名空间，需要读取 `/proc/<pid>/environ` 中的值，其使用 `\u0000` 分割，读取信息的函数如下([container/util.go](./container/util.go))

```go
// getEnvsByPid 获取进程的环境变量
func getEnvsByPid(pid string) (envs []string, err error) {
	var (
		path         = fmt.Sprintf("/proc/%s/environ", pid)
		contentBytes []byte
	)

	if contentBytes, err = ioutil.ReadFile(path); err != nil {
		return nil, fmt.Errorf("read file %s failed, %v", path, err)
	}
	envs = strings.Split(string(contentBytes), "\u0000")
	return
}
```

---

## 问题

一个问题，在以daemon运行时，[文件句柄何时关闭？](./container/log.go)

exec 执行 sh 会报错（影响不大）：

目前只能使用 main.go 中的命令，因为处理传入的 arg 比较麻烦，而且这个代码的重点也不是在这里

目前也不支持开机自动启动容器

不支持传入自定义命令，而是在 [main.go](main.go) 中硬编码，默认值为 `sh -c while true ; do sleep 2; echo [$$] $(date); done` ，也就是每个两秒输出当前进程的 pid （当然，为1）以及当前的日期

同时对传入程序的 arg 没有做验证，默认合法

---

## v0.1 功能

请 cd 到 `build/` 文件夹中自行解压 busybox.rar 至同级目录的 busybox 文件夹中（build/busybox）

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

---

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
