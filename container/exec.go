package container

import (
	"fmt"
	"go-container-process-survey/cgo_key"
	"os"
	"os/exec"
	"strings"
)

// Exec 进入到容器内部
func Exec(containerID string, cmdArr []string) (err error) {
	var (
		containerInfo *Info
		cmdStr        = strings.Join(cmdArr, " ")
		cmd           *exec.Cmd
		envs          []string
	)

	if containerInfo, err = GetContainerInfoByID(containerID); err != nil {
		return fmt.Errorf("get container %s info failed : %v", containerID, err)
	}

	if envs, err = getEnvsByPid(containerInfo.Pid); err != nil {
		return fmt.Errorf("get container %s env failed, %v", containerInfo.Pid, err)
	}

	cmd = exec.Command("/proc/self/exe", "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = os.Setenv(cgo_key.EnvPID, containerInfo.Pid)
	_ = os.Setenv(cgo_key.EnvCMD, cmdStr)
	_ = os.Setenv(cgo_key.EnvFlag, "true")

	cmd.Env = append(os.Environ(), envs...)

	if err = cmd.Run(); err != nil {
		return fmt.Errorf("Exec container %s error %v", containerID, err)
	}
	return
}
