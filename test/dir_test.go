package test

import (
	"io/ioutil"
	"os"
	"syscall"
	"testing"

	"github.com/mingforpc/fuse-go/fuse"
)

func TestReaddir(t *testing.T) {
	tempPoint, err := createTempPoint()

	if err != nil {
		t.Fatalf("TestReaddir err: %+v \n", err)
	}

	opts := fuse.Opt{}
	opts.Getattr = &getattr
	opts.Lookup = &lookup
	opts.Readdir = &readdir

	se := NewTestFuse(tempPoint, opts)

	err = preTest(se)

	if err != nil {
		t.Fatalf("TestReaddir err: %+v \n", err)
	}

	go se.FuseLoop()
	defer exitTest(se)

	wait.Wait()

	// readdir

	// read root folder
	fis, err := ioutil.ReadDir(tempPoint)

	if err != nil {
		t.Errorf("Failed to read dir: %+v \n", err)
	}

	if len(fis) != 2 {
		t.Errorf("count of files under root[%s] should be %d \n", tempPoint, 2)
	}
	for _, fi := range fis {
		name := fi.Name()
		if name != rootFile.name && name != rootDir.name {
			t.Errorf("not exist file [%s] \n", fi.Name())
		}
	}

}

func TestFsyncdir(t *testing.T) {
	tempPoint, err := createTempPoint()

	if err != nil {
		t.Fatalf("TestReaddir err: %+v \n", err)
	}

	opts := fuse.Opt{}
	opts.Getattr = &getattr
	opts.Lookup = &lookup
	opts.Fsyncdir = &fsyncdir

	se := NewTestFuse(tempPoint, opts)

	err = preTest(se)

	if err != nil {
		t.Fatalf("TestReaddir err: %+v \n", err)
	}

	go se.FuseLoop()
	defer exitTest(se)

	wait.Wait()

	// fsyncdir
	f, err := os.Open(tempPoint + "/" + rootDir.name)

	if err != nil {
		t.Fatalf("Failed to open [%s]: %+v \n", rootDir.path, err)
	}

	err = f.Sync()
	if err != nil {
		t.Fatalf("Failed to call syncdir: %+v \n", err)
	}
}

func TestMkdir(t *testing.T) {
	tempPoint, err := createTempPoint()

	if err != nil {
		t.Fatalf("TestReaddir err: %+v \n", err)
	}

	opts := fuse.Opt{}
	opts.Getattr = &getattr
	opts.Lookup = &lookup
	opts.Mkdir = &mkdir

	se := NewTestFuse(tempPoint, opts)

	err = preTest(se)

	if err != nil {
		t.Fatalf("TestReaddir err: %+v \n", err)
	}

	go se.FuseLoop()
	defer exitTest(se)

	wait.Wait()

	dirPath := tempPoint + "/" + "newdir"
	err = os.Mkdir(dirPath, os.ModeDir)

	if err != nil {
		t.Fatalf("Failed to mkdir [%s]: %+v \n", dirPath, err)
	}
}

func TestRmdir(t *testing.T) {
	tempPoint, err := createTempPoint()

	if err != nil {
		t.Fatalf("TestReaddir err: %+v \n", err)
	}

	opts := fuse.Opt{}
	opts.Getattr = &getattr
	opts.Lookup = &lookup
	opts.Mkdir = &mkdir
	opts.Rmdir = &rmdir

	se := NewTestFuse(tempPoint, opts)

	err = preTest(se)

	if err != nil {
		t.Fatalf("TestReaddir err: %+v \n", err)
	}

	go se.FuseLoop()
	defer exitTest(se)

	wait.Wait()

	// rm not empty dir
	rmPath := tempPoint + "/" + rootDir.path
	// avoid use os.Remove(), due to os.Remove() will call unlink() first
	err = syscall.Rmdir(rmPath)
	if err != syscall.ENOTEMPTY {
		t.Errorf("remove not empty dir should raise ENOTEMPTY, but err: %+v \n", err)
	}

	// create empty dir
	dirPath := tempPoint + "/" + "newdir"
	err = os.Mkdir(dirPath, os.ModeDir)

	if err != nil {
		t.Errorf("Failed to mkdir [%s]: %+v \n", dirPath, err)
	}
	// rm empty dir
	rmPath = dirPath
	err = syscall.Rmdir(rmPath)
	if err != nil {
		t.Errorf("Failed to rmove [%s]: %+v \n", dirPath, err)
	}
}
