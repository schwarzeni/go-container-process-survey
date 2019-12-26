package container

import (
	"fmt"
	"path"
)

var (
	// RUNNING 运行状态
	RUNNING string = "running"
	// STOP 停止状态
	STOP string = "stopped"
	// EXIT 退出状态
	EXIT string = "exited"
	// defaultInfoLocation 默认存储位置
	defaultInfoLocation string = "/var/run/mydocker/%s/"
	// ConfigName 配置文件标题
	ConfigName string = "config.json"
	// LogFileName 容器的日志文件
	LogFileName string = "container.log"
	// IDLen ID长度
	IDLen int = 10
	// writerLayerMntPoint 可写层挂载位置
	writerLayerMntPoint = "aufs/writelayer"
	// mntPoint 独立文件系统挂载位置
	mntPoint = "aufs/mnt"
	// 默认导出文件的路径
	defaultCommitDIR = "/tmp"
)

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
}

// Info 存储容器相关信息的json结构
type Info containerInfo

func getContainerInfoDir(id string) string {
	return fmt.Sprintf(defaultInfoLocation, id)
}

func getContainerMntPoint(id string) string {
	return path.Join(getContainerInfoDir(id), mntPoint)
}

func getContainerWriterLayerDir(id string) string {
	return path.Join(getContainerInfoDir(id), writerLayerMntPoint)
}
