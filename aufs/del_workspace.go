package aufs

import (
	"fmt"
	"os"
	"os/exec"
)

// DeleteWorkSpace 删除容器层
func DeleteWorkSpace(mntURL string, writerLayerURL string, volumes []string) (err error) {
	if err = DeleteVolumes(mntURL, volumes); err != nil {
		goto ERR
	}
	if err = DeleteMountPoint(mntURL); err != nil {
		goto ERR
	}
	if err = deleteWriteLayer(writerLayerURL); err != nil {
		goto ERR
	}
	return
ERR:
	// TODO: handle error here
	return
}

// DeleteVolumes 删除挂载的数据卷
func DeleteVolumes(mntURL string, volumes []string) (err error) {
	for _, volume := range volumes {
		if len(volume) != 0 {
			var volumeURLs []string
			if volumeURLs, err = volumeURLExtract(volume); err != nil {
				return
			}
			if err = umountVolume(mntURL, volumeURLs); err != nil {
				return
			}
		}
	}
	return
}

// DeleteMountPoint 删除挂载点
func DeleteMountPoint(mntURL string) (err error) {
	cmd := exec.Command("umount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("unmount %s error %v", mntURL, err)
	}
	if err = os.RemoveAll(mntURL); err != nil {
		return fmt.Errorf("removeAll %s error %v", mntURL, err)
	}
	return
}

// deleteWriteLayer 删除可写层
func deleteWriteLayer(writerLayerURL string) (err error) {
	if err = os.RemoveAll(writerLayerURL); err != nil {
		return fmt.Errorf("removeAll %s error %v", writerLayerURL, err)
	}
	return
}
