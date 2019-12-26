package main

import (
	"flag"
	_ "go-container-process-survey/cgo"
	"go-container-process-survey/cgo_key"
	"go-container-process-survey/container"
	"log"
	"os"
	"os/exec"
)

var (
	runAsDaemon   = flag.Bool("d", false, "后台运行")
	containerName = flag.String("name", "", "容器的名称")
	imagePath     = flag.String("image", "/root/workplace_go/go-container-process-survey/build/busybox", "容器镜像文件系统的位置")
	volumes       = flag.String("v", "/tmp/just_for_test:/test_for_volume", "挂载的数据卷")
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
		if len(os.Args) > 1 && os.Args[1] == "start" { // start a container in child process
			fullCmd, _ := exec.LookPath(os.Args[2])
			childProcess(fullCmd, os.Args[2:])
			return
		}
		fullCmd, _ := exec.LookPath(defaultCmd[0])
		childProcess(fullCmd, defaultCmd)
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
	if os.Args[1] == "rm" { // delete a container
		if err := container.RemoveContainerByID(os.Args[2]); err != nil {
			log.Fatal(err)
		}
		return
	}
	if os.Args[1] == "start" { // start a stop container
		if err := container.StartContainerByID(os.Args[2]); err != nil {
			log.Fatal(err)
		}
		return
	}
	if os.Args[1] == "commit" { // pack a container into a image: commit <id> <path>
		if err := container.CommitContainer(os.Args[2], os.Args[3]); err != nil {
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

	if err := container.RunContainer(defaultCmd, *runAsDaemon, *containerName, *imagePath, []string{*volumes}); err != nil {
		log.Fatalf("Run container failed: %v", err)
	}
}
