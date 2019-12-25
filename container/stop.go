package container

import (
	"fmt"
	"strconv"
	"syscall"
)

// StopContainerByID 根据容器的ID停止容器
func StopContainerByID(id string) (err error) {
	var (
		containerInfo *Info
		pidInt        int
	)
	if containerInfo, err = GetContainerInfoByID(id); err != nil {
		return fmt.Errorf("get container %s info failed, %v", id, err)
	}

	if pidInt, err = strconv.Atoi(containerInfo.Pid); err != nil {
		return fmt.Errorf("parse container pid %s failed, %v", containerInfo.Pid, err)
	}

	// send SIGTERM, and then SIGKILL after grace period
	if err = syscall.Kill(pidInt, syscall.SIGKILL); err != nil {
		return fmt.Errorf("stop container %s failed, %v", containerInfo.Pid, err)
	}

	containerInfo.Status = STOP
	containerInfo.Pid = ""

	if err = UpdateContainerInfo(containerInfo); err != nil {
		return fmt.Errorf("update local file of container %s failed, %v", containerInfo.ID, err)
	}
	return
}
