package aufs

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// NewWorkSpace 创建新的容器层
func NewWorkSpace(imageURL string, mntURL string, writerLayerURL string, volume string) (err error) {
	if err = createReadOnlyLayer(imageURL); err != nil {
		goto ERR
	}
	if err = createWriteLayer(writerLayerURL); err != nil {
		goto ERR
	}
	if err = CreateMountPoint(mntURL, imageURL, writerLayerURL); err != nil {
		goto ERR
	}
	if len(volume) != 0 {
		var volumeURLs []string
		if volumeURLs, err = volumeURLExtract(volume); err != nil {
			goto ERR
		}
		// if err = mountVolume(rootURL, mntURL, volumeURLs); err != nil {
		// goto ERR
		// }
		log.Printf("%s mount on %s in container", volumeURLs[0], volumeURLs[1])
	}
	return
ERR:
	deleteWriteLayer(writerLayerURL)
	DeleteMountPoint(mntURL)
	return
}

// createReadOnlyLayer 暂时不支持解压进行镜像
func createReadOnlyLayer(imageURL string) (err error) {
	// var exist bool
	// busyboxURL := path.Join(rootURL, "busybox")
	// busyboxTarURL := path.Join(rootURL, "busybox.rar")
	// exist, err = pathExists(busyboxURL)
	// if err != nil {
	// 	return fmt.Errorf("Failed to judge whether dir %s exists. %v", busyboxURL, err)
	// }
	// if exist == false {
	// 	if err = os.Mkdir(busyboxURL, 0777); err != nil {
	// 		return fmt.Errorf("Midir dir %s error. %v", busyboxURL, err)
	// 	}
	// 	if _, err = exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
	// 		return fmt.Errorf("unTar dir %s error %v", busyboxTarURL, err)
	// 	}
	// }
	return
}

// createWriteLayer 创建容器可写层
func createWriteLayer(writerLayerURL string) (err error) {
	if err = os.MkdirAll(writerLayerURL, 0777); err != nil {
		return fmt.Errorf("Mkdir dir %s error. %v", writerLayerURL, err)
	}
	return
}

// CreateMountPoint 将只读层和可写层都mount到一处
func CreateMountPoint(mntURL string, imageURL string, writerLayerURL string) (err error) {
	if err = os.MkdirAll(mntURL, 0777); err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("Mkdir dir %s error. %v", mntURL, err)
		}
		err = nil
	}
	dirs := fmt.Sprintf("dirs=%s:%s", writerLayerURL, imageURL)
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)
	cmd.Stdout = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cmd %s %v error %v", cmd.Path, cmd.Args, err)
	}
	return
}
