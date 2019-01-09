package test

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mingforpc/fuse-go/fuse"

	"github.com/mingforpc/fuse-go/fuse/mount"
)

func readFromMtab() ([]string, error) {

	content, err := ioutil.ReadFile("/etc/mtab")

	if err != nil {
		return nil, err
	}

	contentStr := string(content)

	lines := strings.Split(contentStr, "\n")

	return lines, err
}

// TestMount : Test fuse mount and umount, check it in /etc/mtab
func TestMount(t *testing.T) {

	tempPoint, err := createTempPoint()

	if err != nil {
		t.Errorf("create temp point error: %+v \n", err)
	}

	se := NewTestFuse(tempPoint, fuse.Opt{})

	// call mount
	err = mount.Mount(se)

	if err != nil {
		t.Errorf("Mount error: %+v \n", err)
	}

	// check if mount success in system,
	// get mountpoint in /etc/mtab
	var lines []string

	lines, err = readFromMtab()

	if err != nil {
		t.Errorf("Failed to read /etc/mtab: %+v \n", err)
	}

	abspath := tempPoint

	if !filepath.IsAbs(abspath) {
		abspath, _ = filepath.Abs(abspath)
	}

	var isMount bool

	for _, line := range lines {

		if strings.HasPrefix(line, "/dev/fuse") && strings.Contains(line, abspath) {
			isMount = true
		}

	}

	if isMount == false {
		t.Errorf("Mount Failed, cannot find in /etc/mtab \n")
	}

	// call umount
	err = mount.Unmount(se.Mountpoint)

	if err != nil {
		t.Errorf("Unmount error: %+v \n", err)
	}

	// check if unmount success in system,
	// get mountpoint in /etc/mtab
	lines, err = readFromMtab()

	if err != nil {
		t.Errorf("Failed to read /etc/mtab: %+v \n", err)
	}

	isUnMount := true

	for _, line := range lines {

		if strings.HasPrefix(line, "/dev/fuse") && strings.Contains(line, abspath) {
			isUnMount = false
		}

	}

	if isUnMount == false {
		t.Errorf("UnMount Failed, find it in /etc/mtab \n")
	}

}
