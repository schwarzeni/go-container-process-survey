package container

import (
	"fmt"
	"io"
	"os"
	"path"
)

// StdLog 如果容器是以后台运行形式，则将其log输出至指定文件
func StdLog(containerID string) (stdLog *os.File, err error) {
	var (
		containerInfoDir = getContainerInfoDir(containerID)
		logFilePath      = path.Join(containerInfoDir, LogFileName)
	)
	// 生成存储的文件夹以及文件
	if err = os.MkdirAll(containerInfoDir, 0755); err != nil {
		return nil, fmt.Errorf("mkdir %s error, %v", containerInfoDir, err)
	}
	if stdLog, err = os.Create(logFilePath); err != nil {
		return nil, fmt.Errorf("create file %s error, %v", logFilePath, err)
	}
	return
}

// GetLogContent 读取后台运行程序的日志输出
func GetLogContent(containerID string) (err error) {
	var (
		containerInfoDir = getContainerInfoDir(containerID)
		logFilePath      = path.Join(containerInfoDir, LogFileName)
		file             *os.File
	)
	if file, err = os.Open(logFilePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no such container with id %s or it has no log file", containerID)
		}
		return fmt.Errorf("open file %s error, %v", logFilePath, err)
	}
	defer file.Close() // TODO: handle error here

	if _, err = io.Copy(os.Stdout, file); err != nil {
		return fmt.Errorf("read file %s error, %v", logFilePath, err)
	}

	return
}
