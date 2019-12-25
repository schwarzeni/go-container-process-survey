package container

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// RecordContainerInfo 记录容器信息, 返回生成的容器的ID
func RecordContainerInfo(id string, containerPID int, cmdArr []string, containerName string) (err error) {
	var (
		cmd        = strings.Join(cmdArr, " ")
		createTime = time.Now().Format("2006-01-02 15:04:05") // 以当前的时间作为创建时间
		info       *containerInfo                             // 新建的容器的信息
		dirURL     string                                     // 存储数据的文件夹的路径
		jsonBytes  []byte
		file       *os.File
		filePath   string
	)
	// id = randStringBytes(IDLen)  // 生成容器的ID号
	if len(containerName) == 0 { // 未指定容器名，则使用ID作为其名字
		containerName = id
	}
	info = &containerInfo{
		ID:          id,
		Pid:         strconv.Itoa(containerPID),
		Command:     cmd,
		CreatedTime: createTime,
		Status:      RUNNING,
		Name:        containerName,
	}

	// 将容器信息的对象 json 序列化成字符串
	if jsonBytes, err = json.Marshal(info); err != nil {
		return fmt.Errorf("convert container data %#v into json error, %v", *info, err)
	}

	// 生成存储的文件夹以及文件
	dirURL = getContainerInfoDir(id)
	if err = os.MkdirAll(dirURL, 0755); err != nil {
		return fmt.Errorf("mkdir %s error, %v", dirURL, err)
	}
	filePath = path.Join(dirURL, ConfigName)
	if file, err = os.Create(filePath); err != nil {
		_ = os.RemoveAll(dirURL)
		return fmt.Errorf("create file %s error, %v", filePath, err)
	}

	// 写入文件
	if _, err := file.Write(jsonBytes); err != nil {
		_ = os.RemoveAll(dirURL)
		return fmt.Errorf("write to file %s error, %v", filePath, err)
	}

	return
}

// DeleteContainerInfo 删除容器的相关信息
func DeleteContainerInfo(containerID string) (err error) {
	if len(containerID) == 0 {
		return
	}
	dirURL := getContainerInfoDir(containerID)
	if err = os.RemoveAll(dirURL); err != nil {
		return fmt.Errorf("remove %s failed, %v", dirURL, err)
	}
	return
}

// ListContainers 列出container信息
func ListContainers() (err error) {
	var (
		dirURL     = fmt.Sprintf(defaultInfoLocation, "")
		files      []os.FileInfo
		containers []*containerInfo
	)
	if files, err = ioutil.ReadDir(dirURL); err != nil {
		return fmt.Errorf("read dir %s : %v", dirURL, err)
	}

	// 遍历该文件夹下的所有文件
	// TODO: 是否并行读取文件？
	for _, fileinfo := range files {
		var (
			filePath  = path.Join(fmt.Sprintf(defaultInfoLocation, fileinfo.Name()), ConfigName)
			byteData  []byte
			container containerInfo
		)
		if byteData, err = ioutil.ReadFile(filePath); err != nil {
			return fmt.Errorf("read file %s : %v", filePath, err)
		}
		if err = json.Unmarshal(byteData, &container); err != nil {
			return fmt.Errorf("unmarshal json data : %v", err)
		}
		containers = append(containers, &container)
	}

	// 此处打印信息就不和书上一样了，意思意思就行了
	for _, container := range containers {
		fmt.Printf("%#v\n", *container)
	}

	return
}
