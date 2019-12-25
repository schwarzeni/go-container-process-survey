package main

import (
	"flag"
	"fmt"
	"go-container-process-survey/container"
	"log"
	"os"
	"os/exec"
	"syscall"
)

var (
	runAsDaemon   = flag.Bool("d", false, "后台运行")
	containerName = flag.String("name", "", "容器的名称")
	defaultCmd    = []string{"sh", "-c", `while true ; do sleep 2; done`}
	// defaultCmd    = []string{"sh", "-c", `while true ; do sleep 2; echo $(date); done`}
	// defaultCmd = []string{"sh", "-c", `for i in $(seq 1 4);do echo "Welcome $i";sleep 1;done`}
)

func init() {
	flag.Parse()
}

func main() {
	if os.Args[0] == "/proc/self/exe" { // child process
		childProcess()
		return
	}
	if os.Args[1] == "ps" { // show process
		if err := container.ListContainers(); err != nil {
			log.Fatal(err)
		}
		return
	}
	var (
		containerID string
		err         error
		cmd         = exec.Command("/proc/self/exe")
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWPID}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start() failed: %v", err)
	}

	if containerID, err = container.RecordContainerInfo(cmd.Process.Pid, defaultCmd, *containerName); err != nil {
		log.Fatalf("record container info failed: %v", err)
	}

	if !*runAsDaemon {
		defer container.DeleteContainerInfo(containerID)
		if err = cmd.Wait(); err != nil {
			log.Fatalf("cmd.Wait %s", err)
		}
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
