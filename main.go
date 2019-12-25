package main

import (
	"flag"
	"fmt"
	_ "go-container-process-survey/cgo"
	"go-container-process-survey/cgo_key"
	"go-container-process-survey/container"
	"log"
	"os"
	"os/exec"
	"syscall"
)

var (
	runAsDaemon   = flag.Bool("d", false, "后台运行")
	containerName = flag.String("name", "", "容器的名称")
	// defaultCmd    = []string{"sh", "-c", `while true ; do sleep 2; done`}
	defaultCmd = []string{"sh", "-c", `while true ; do sleep 2; echo \[$$\] $(date); done`}
	// defaultCmd = []string{"sh", "-c", `for i in $(seq 1 4);do echo "Welcome $i";sleep 1;done`}
)

func init() {
	flag.Parse()
}

func main() {
	if os.Args[0] == "/proc/self/exe" { // child process
		// note here, just a hack ... goto exec part (in child process)
		if len(os.Args) > 1 && os.Args[1] == "exec" {
			goto EXEC
		}
		childProcess()
		return
	}
	if os.Args[1] == "ps" { // show process
		if err := container.ListContainers(); err != nil {
			log.Fatal(err)
		}
		return
	}
	if os.Args[1] == "logs" { // show log, 注意这里就不判断用户是否提供了容器ID了，默认输入合法
		if err := container.GetLogContent(os.Args[2]); err != nil {
			log.Fatal(err)
		}
		return
	}
	if os.Args[1] == "stop" { // stop a container
		if err := container.StopContainerByID(os.Args[2]); err != nil {
			log.Fatal(err)
		}
		return
	}
EXEC:
	if os.Args[1] == "exec" { // 进入正在运行的容器内部，默认输入合法 exec <id> <cmd ...>
		if os.Getenv(cgo_key.EnvFlag) != "" { // using cgo
			return
		}
		id := os.Args[2]
		cmd := os.Args[3:]
		if err := container.Exec(id, cmd); err != nil {
			log.Fatal(err)
		}
		return
	}

	var (
		containerID = container.RandStringBytes(container.IDLen) //  // 生成容器的ID号
		err         error
		cmd         = exec.Command("/proc/self/exe")
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWIPC | syscall.CLONE_NEWUTS | syscall.CLONE_NEWNET}
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if *runAsDaemon { // 如果为后台运行模式，设置输出的文件
		var stdLog *os.File
		if stdLog, err = container.StdLog(containerID); err != nil {
			log.Fatalf("create log file failed: %v", err)
		}
		cmd.Stderr = stdLog
		cmd.Stdout = stdLog
	}
	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start() failed: %v", err)
	}

	if err = container.RecordContainerInfo(containerID, cmd.Process.Pid, defaultCmd, *containerName); err != nil { // 记录容器信息
		log.Fatalf("record container info failed: %v", err)
	}

	if !*runAsDaemon { // 前台运行模式
		defer container.DeleteContainerInfo(containerID)
		if err = cmd.Wait(); err != nil {
			log.Fatalf("cmd.Wait %s", err)
		}
	} else { // 后台运行模式
		log.Println(containerID)
	}
}

func childProcess() {
	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err := syscall.Mount("proc", "/proc", "proc", syscall.MS_NOEXEC|syscall.MS_NOSUID|syscall.MS_NODEV, ""); err != nil {
		fmt.Fprintf(os.Stderr, "mount proc error %v", err)
		return
	}
	if err := syscall.Exec("/bin/sh", defaultCmd, os.Environ()); err != nil {
		fmt.Fprintf(os.Stderr, "exec error %v", err)
		return
	}
}
