package container

import "fmt"

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
)

// 存储容器相关信息的json结构
type containerInfo struct {
	Pid         string `json:"pid"`        // 容器的init的进程在宿主机上的PID
	ID          string `json:"id"`         // 容器ID
	Name        string `json:"name"`       // 容器名称
	Command     string `json:"command"`    // 容器内 init 进程的运行命令
	CreatedTime string `json:"createTime"` // 创建时间
	Status      string `json:"status"`     // 容器状态
}

// Info 存储容器相关信息的json结构
type Info containerInfo

func getContainerInfoDir(id string) string {
	return fmt.Sprintf(defaultInfoLocation, id)
}
