package test

import (
	"os"

	"github.com/mingforpc/fuse-go/fuse"
	"github.com/mingforpc/fuse-go/fuse/mount"
)

func createTempPoint() (string, error) {
	var tempMount = "/home/ming/Downloads/test"
	err := os.Mkdir(tempMount, os.ModePerm)

	if err != nil && !os.IsExist(err) {
		return "", err
	}

	return tempMount, nil
}

func preTest(se *fuse.Session) error {
	wait.Add(1)

	err := mount.Mount(se)

	return err
}

func exitTest(se *fuse.Session) {
	se.Close()
	mount.Unmount(se.Mountpoint)

	os.Remove(se.Mountpoint)
}
