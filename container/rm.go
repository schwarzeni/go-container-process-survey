package container

import (
	"fmt"
	"os"
)

// RemoveContainerByID 删除容器
func RemoveContainerByID(containerID string) (err error) {
	var (
		containerInfo *Info
	)
	if containerInfo, err = GetContainerInfoByID(containerID); err != nil {
		return fmt.Errorf("get container %s info failed, %v", containerID, err)
	}
	if !canRemove(containerInfo) {
		return fmt.Errorf("cannot remove container %s with status %s", containerID, containerInfo.Status)
	}

	dirURL := getContainerInfoDir(containerID)
	if err = os.RemoveAll(dirURL); err != nil {
		return fmt.Errorf("remove container dir %s failed, %v", dirURL, err)
	}

	return
}

func canRemove(containerInfo *Info) bool {
	return containerInfo.Status == STOP
}
