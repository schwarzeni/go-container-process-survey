package container

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// RunContainer 启动一个新容器
func RunContainer(defaultCmd []string, daemon bool, name string) (err error) {
	var (
		containerID = RandStringBytes(IDLen) //  // 生成容器的ID号
		cmd         = exec.Command("/proc/self/exe")
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWIPC | syscall.CLONE_NEWUTS | syscall.CLONE_NEWNET}
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if daemon { // 如果为后台运行模式，设置输出的文件
		var stdLog *os.File
		if stdLog, err = StdLog(containerID); err != nil {
			return fmt.Errorf("create log file failed: %v", err)
		}
		cmd.Stderr = stdLog
		cmd.Stdout = stdLog
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cmd.Start() failed: %v", err)
	}

	if err = RecordContainerInfo(containerID, cmd.Process.Pid, defaultCmd, name); err != nil { // 记录容器信息
		return fmt.Errorf("record container info failed: %v", err)
	}

	if !daemon { // 前台运行模式
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt) // 监听 ctrl+c 事件
		go func() {
			for range c {
				DeleteContainerInfo(containerID)
				if err := cmd.Process.Signal(os.Interrupt); err != nil {
					log.Printf("send signal to child process failed, %v", err)
				}
				break
			}
			close(c)
		}()
		if err = cmd.Wait(); err != nil {
			if cmd.ProcessState.ExitCode() == 130 { // ignore ctrl+c errror
				return nil
			}
			return fmt.Errorf("cmd.Wait %v", err)
		}
	} else { // 后台运行模式
		log.Println(containerID)
	}
	return
}
