package test

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mingforpc/fuse-go/fuse"
	"github.com/mingforpc/fuse-go/fuse/mount"
)

func createTempPoint() (string, error) {

	var tempMount = "./test_point" + strconv.Itoa(rand.Intn(100))
	tempMount, _ = filepath.Abs(tempMount)
	fmt.Println(tempMount)
	err := os.Mkdir(tempMount, os.ModePerm)

	if err != nil && !os.IsExist(err) {
		return "", err
	}

	return tempMount, nil
}

func preTest(se *fuse.Session) error {
	wait.Add(1)

	err := mount.Mount(se, nil)

	return err
}

func preTestArgs(se *fuse.Session, args []string) error {
	wait.Add(1)

	err := mount.Mount(se, args)

	return err
}

func exitTest(se *fuse.Session) {
	se.Close()
	mount.Unmount(se.Mountpoint)

	os.Remove(se.Mountpoint)
}
