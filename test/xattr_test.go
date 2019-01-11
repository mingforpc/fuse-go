package test

import (
	"strings"
	"syscall"
	"testing"

	"github.com/mingforpc/fuse-go/fuse"
)

func TestSetGetxattr(t *testing.T) {
	tempPoint, err := createTempPoint()

	if err != nil {
		t.Fatalf("TestSetxattr err: %+v \n", err)
	}

	opts := fuse.Opt{}
	opts.Getattr = &getattr
	opts.Lookup = &lookup
	opts.Setxattr = &setxattr
	opts.Getxattr = &getxattr

	se := NewTestFuse(tempPoint, opts)

	err = preTest(se)

	if err != nil {
		t.Fatalf("TestSetxattr err: %+v \n", err)
	}

	go se.FuseLoop()
	defer exitTest(se)

	wait.Wait()

	// setxattr
	path := tempPoint + "/" + rootFile.path
	err = syscall.Setxattr(path, "system.name", []byte("test"), 0)
	if err != nil {
		t.Fatalf("Failed to setxattr: %+v \n", err)
	}

	// getxattr
	buf := make([]byte, 128)
	n, err := syscall.Getxattr(path, "system.name", buf)
	if err != nil {
		t.Fatalf("Failed to getxattr: %+v \n", err)
	}
	// n-1 to avoid ‘\0'
	content := string(buf[:n-1])

	if content != "test" {
		t.Fatalf("xattr value shoud be %s \n", "test")
	}
}

func TestListxattr(t *testing.T) {
	tempPoint, err := createTempPoint()

	if err != nil {
		t.Fatalf("TestListxattr err: %+v \n", err)
	}

	opts := fuse.Opt{}
	opts.Getattr = &getattr
	opts.Lookup = &lookup
	opts.Setxattr = &setxattr
	opts.Listxattr = &listxattr

	se := NewTestFuse(tempPoint, opts)

	err = preTest(se)

	if err != nil {
		t.Fatalf("TestListxattr err: %+v \n", err)
	}

	go se.FuseLoop()
	defer exitTest(se)

	wait.Wait()

	// setxattr
	path := tempPoint + "/" + rootFile.path
	err = syscall.Setxattr(path, "system.name", []byte("test"), 0)
	if err != nil {
		t.Fatalf("Failed to setxattr: %+v \n", err)
	}
	err = syscall.Setxattr(path, "system.type", []byte("type"), 0)
	if err != nil {
		t.Fatalf("Failed to setxattr: %+v \n", err)
	}

	// listxattr
	buf := make([]byte, 128)
	n, err := syscall.Listxattr(path, buf)
	if err != nil {
		t.Fatalf("Failed to listxattr: %+v \n", err)
	}
	content := string(buf[:n])
	xattrs := strings.Split(content, string(byte(0)))

	for _, xattr := range xattrs {
		if xattr != "" && xattr != "system.name" && xattr != "system.type" {
			t.Errorf("listxattr not correct, [%s] should not exist \n", xattr)
		}
	}

}

func TestRemovexattr(t *testing.T) {
	tempPoint, err := createTempPoint()

	if err != nil {
		t.Fatalf("TestRemovexattr err: %+v \n", err)
	}

	opts := fuse.Opt{}
	opts.Getattr = &getattr
	opts.Lookup = &lookup
	opts.Readdir = &readdir
	opts.Setxattr = &setxattr
	opts.Getxattr = &getxattr
	opts.Removexattr = &removexattr

	se := NewTestFuse(tempPoint, opts)

	err = preTest(se)

	if err != nil {
		t.Fatalf("TestRemovexattr err: %+v \n", err)
	}

	go se.FuseLoop()
	defer exitTest(se)

	wait.Wait()

	// setxattr
	path := tempPoint + "/" + rootFile.path
	err = syscall.Setxattr(path, "system.name", []byte("test"), 0)
	if err != nil {
		t.Fatalf("Failed to setxattr: %+v \n", err)
	}

	// getxattr
	buf := make([]byte, 128)
	n, err := syscall.Getxattr(path, "system.name", buf)
	if err != nil {
		t.Fatalf("Failed to getxattr: %+v \n", err)
	}
	// n-1 to avoid ‘\0'
	content := string(buf[:n-1])

	if content != "test" {
		t.Errorf("xattr value shoud be %s \n", "test")
	}

	// removexattr
	err = syscall.Removexattr(path, "system.name")
	if err != nil {
		t.Fatal("Failed to removexattr")
	}

	// getxattr
	n, err = syscall.Getxattr(path, "system.name", buf)
	if err != syscall.ENODATA {
		t.Fatalf("Failed to getxattr: %+v \n", err)
	}
}
