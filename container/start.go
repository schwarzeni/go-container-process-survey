package container

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

// StartContainerByID 启动已经停止的容器
func StartContainerByID(id string) (err error) {
	var (
		containerInfo *Info
		cmd           *exec.Cmd
		stdLog        *os.File
		fullcmd       = []string{"start"}
	)
	if containerInfo, err = GetContainerInfoByID(id); err != nil {
		return fmt.Errorf("get container %s from local failed, %v", id, err)
	}
	if !canStart(containerInfo) {
		return fmt.Errorf("cannot start current container %s with status %s", id, containerInfo.Status)
	}

	for _, arg := range containerInfo.FullCommand {
		fullcmd = append(fullcmd, arg)
	}
	cmd = exec.Command("/proc/self/exe", fullcmd...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWIPC | syscall.CLONE_NEWUTS | syscall.CLONE_NEWNET}
	if stdLog, err = StdAppendLog(id); err != nil {
		return fmt.Errorf("get stdlog for container %s failed, %v", id, err)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = stdLog
	cmd.Stderr = stdLog

	if err = cmd.Start(); err != nil {
		containerInfo.Status = EXIT
		if err := UpdateContainerInfo(containerInfo); err != nil {
			return fmt.Errorf("update container %s info failed, %v", id, err)
		}
		return fmt.Errorf("start container %s failed, %v", id, err)
	}
	containerInfo.Pid = strconv.Itoa(cmd.Process.Pid)
	containerInfo.Status = RUNNING
	if err = UpdateContainerInfo(containerInfo); err != nil {
		return fmt.Errorf("update container %s info failed, %v", id, err)
	}
	return
}

func canStart(info *Info) bool {
	return info.Status == STOP || info.Status == EXIT
}
