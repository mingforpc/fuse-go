package test

import (
	"os"
	"syscall"
	"testing"

	"github.com/mingforpc/fuse-go/fuse"
)

func TestMknod(t *testing.T) {
	tempPoint, err := createTempPoint()

	if err != nil {
		t.Fatalf("TestMknod err: %+v \n", err)
	}

	opts := fuse.Opt{}
	opts.Getattr = &getattr
	opts.Lookup = &lookup
	opts.Mknod = &mknod

	se := NewTestFuse(tempPoint, opts)

	err = preTest(se)

	if err != nil {
		t.Fatalf("TestMknod err: %+v \n", err)
	}

	go se.FuseLoop()
	defer exitTest(se)

	// mknod
	newFile := tempPoint + "/" + "new_test"
	err = syscall.Mknod(newFile, uint32(syscall.S_IFREG)|uint32(0444), 0)
	if err != nil {
		t.Fatalf("Failed to call mknod: %+v \n", err)
	}

	// lookup
	fi, err := os.Stat(newFile)
	if err != nil {
		t.Fatalf("Failed to lookup new file: %+v \n", err)
	}
	if fi.Name() != "new_test" {
		t.Fatalf("new file name should be [%s]\n", "new_test")
	}
}

func TestUnlink(t *testing.T) {
	tempPoint, err := createTempPoint()

	if err != nil {
		t.Fatalf("TestUnlink err: %+v \n", err)
	}

	opts := fuse.Opt{}
	opts.Getattr = &getattr
	opts.Lookup = &lookup
	opts.Mknod = &mknod
	opts.Unlink = &unlink

	se := NewTestFuse(tempPoint, opts)

	err = preTest(se)

	if err != nil {
		t.Fatalf("TestUnlink err: %+v \n", err)
	}

	go se.FuseLoop()
	defer exitTest(se)

	// mknod
	newFile := tempPoint + "/" + "new_test"
	err = syscall.Mknod(newFile, uint32(syscall.S_IFREG)|uint32(0444), 0)
	if err != nil && !os.IsExist(err) {
		t.Fatalf("Failed to call mknod: %+v \n", err)
	}

	// lookup
	fi, err := os.Stat(newFile)
	if err != nil {
		t.Fatalf("Failed to lookup new file: %+v \n", err)
	}
	if fi.Name() != "new_test" {
		t.Fatalf("new file name should be [%s]\n", "new_test")
	}

	// unlink
	err = syscall.Unlink(newFile)
	if err != nil {
		t.Fatalf("Failed to unlink new file: %+v \n", err)
	}

	// lookup
	fi, err = os.Stat(newFile)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("new file should be unlinked: %+v \n", err)
	}
}

func TestRename(t *testing.T) {
	tempPoint, err := createTempPoint()

	if err != nil {
		t.Fatalf("TestRename err: %+v \n", err)
	}

	opts := fuse.Opt{}
	opts.Getattr = &getattr
	opts.Lookup = &lookup
	opts.Mknod = &mknod
	opts.Rename = &rename

	se := NewTestFuse(tempPoint, opts)

	err = preTest(se)

	if err != nil {
		t.Fatalf("TestRename err: %+v \n", err)
	}

	go se.FuseLoop()
	defer exitTest(se)

	// mknod
	newFile := tempPoint + "/" + "new_test"
	err = syscall.Mknod(newFile, uint32(syscall.S_IFREG)|uint32(0444), 0)
	if err != nil && !os.IsExist(err) {
		t.Fatalf("Failed to call mknod: %+v \n", err)
	}

	// lookup
	_, err = os.Stat(newFile)
	if err != nil {
		t.Fatalf("Failed to lookup new file: %+v \n", err)
	}

	// rename
	renamePath := tempPoint + "/" + "rename_test"
	err = os.Rename(newFile, renamePath)
	if err != nil {
		t.Fatalf("Failed to rename: %+v \n", err)
	}

	// lookup
	_, err = os.Stat(newFile)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("old name still exist: %+v \n", err)
	}

	_, err = os.Stat(renamePath)
	if err != nil {
		t.Fatalf("reanme file not exist: %+v \n", err)
	}
}

func TestLink(t *testing.T) {
	tempPoint, err := createTempPoint()

	if err != nil {
		t.Fatalf("TestLink err: %+v \n", err)
	}

	opts := fuse.Opt{}
	opts.Getattr = &getattr
	opts.Lookup = &lookup
	opts.Link = &link

	se := NewTestFuse(tempPoint, opts)

	err = preTest(se)

	if err != nil {
		t.Fatalf("TestLink err: %+v \n", err)
	}

	go se.FuseLoop()
	defer exitTest(se)

	// link
	err = os.Link(tempPoint+"/"+rootFile.path, tempPoint+"/"+"hardlink")
	if err != nil {
		t.Fatalf("create hardlink err: %+v \n", err)
	}

	// lookup
	_, err = os.Stat(tempPoint + "/" + "hardlink")
	if err != nil {
		t.Fatal("Failed to lookup new hard link")
	}

}

// open(), read(), write(), fsync(), flush(), release()
func TestOpenReadWrite(t *testing.T) {
	tempPoint, err := createTempPoint()

	if err != nil {
		t.Fatalf("TestLink err: %+v \n", err)
	}

	opts := fuse.Opt{}
	opts.Getattr = &getattr
	opts.Lookup = &lookup
	opts.Open = &open
	opts.Read = &read
	opts.Write = &write
	opts.Fsync = &fsync
	opts.Flush = &flush
	opts.Release = &release

	se := NewTestFuse(tempPoint, opts)

	err = preTest(se)

	if err != nil {
		t.Fatalf("TestLink err: %+v \n", err)
	}

	go se.FuseLoop()
	defer exitTest(se)

	//open
	path := tempPoint + "/" + rootFile.path
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		t.Fatalf("Failed to open file: %+v \n", err)
	}

	// read
	buf := make([]byte, 1024)
	n, err := file.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read file: %+v \n", err)
	}
	content := string(buf[:n])
	if content != rootFile.content {
		t.Fatal("read content correct!")
	}

	// write
	orginContent := rootFile.content

	n, err = file.Write([]byte("123"))
	if err != nil {
		t.Fatalf("Failed to write file: %+v \n", err)
	}

	// sync
	err = file.Sync()
	if err != nil {
		t.Fatalf("Failed to sync file: %+v \n", err)
	}

	if orginContent+"123" != rootFile.content {
		t.Fatal("file content not correct ")
	}

	// flush, release
	err = file.Close()
	if err != nil {
		t.Fatal("Failed to close file")
	}
}

func TestStatfs(t *testing.T) {
	tempPoint, err := createTempPoint()

	if err != nil {
		t.Fatalf("TestLink err: %+v \n", err)
	}

	opts := fuse.Opt{}
	opts.Getattr = &getattr
	opts.Lookup = &lookup
	opts.Statfs = &statfs

	se := NewTestFuse(tempPoint, opts)

	err = preTest(se)

	if err != nil {
		t.Fatalf("TestLink err: %+v \n", err)
	}

	go se.FuseLoop()
	defer exitTest(se)

	// statfs
	buf := syscall.Statfs_t{}
	err = syscall.Statfs(tempPoint, &buf)
	if err != nil {
		t.Fatalf("Failed to call statfs: %+v \n", err)
	}

	if !compareStatfs(buf) {
		t.Fatal("statfs data not correct")
	}
}

func compareStatfs(statfs syscall.Statfs_t) bool {
	if statfs.Bsize != int64(rootStatfs.Bsize) {
		return false
	}
	if statfs.Blocks != rootStatfs.Blocks {
		return false
	}
	if statfs.Bfree != rootStatfs.Bfree {
		return false
	}

	if statfs.Bavail != rootStatfs.Bavail {
		return false
	}
	if statfs.Files != rootStatfs.Files {
		return false
	}
	if statfs.Ffree != rootStatfs.Ffree {
		return false
	}
	if statfs.Namelen != int64(rootStatfs.NameLen) {
		return false
	}

	return true
}
