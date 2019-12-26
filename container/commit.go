package container

import (
	"fmt"
	"log"
	"os/exec"
	"path"
)

// CommitContainer 导出容器为镜像，tar打包
func CommitContainer(id string, outputPath string) (err error) {
	var (
		mntPath        = getContainerMntPoint(id)
		outputFilePath = path.Join(defaultCommitDIR, id+".tar")
	)
	if len(outputPath) != 0 {
		outputFilePath = outputPath
	}
	if _, err = exec.Command("tar", "-czf", outputFilePath, "-C", mntPath, ".").CombinedOutput(); err != nil {
		return fmt.Errorf("pack container using tar command failed, %v", err)
	}
	log.Println("save to file: " + outputFilePath)
	return
}
