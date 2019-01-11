package test

import (
	"os"
	"testing"

	"github.com/mingforpc/fuse-go/fuse"
)

// create symlink and read link
func TestSymlink(t *testing.T) {
	tempPoint, err := createTempPoint()

	if err != nil {
		t.Fatalf("TestReaddir err: %+v \n", err)
	}

	opts := fuse.Opt{}
	opts.Getattr = &getattr
	opts.Lookup = &lookup
	opts.Symlink = &symlink
	opts.Readlink = &readlink

	se := NewTestFuse(tempPoint, opts)

	err = preTest(se)

	if err != nil {
		t.Fatalf("TestReaddir err: %+v \n", err)
	}

	go se.FuseLoop()
	defer exitTest(se)

	wait.Wait()

	// create symlink
	oldPath := tempPoint + "/" + rootFile.path
	newPath := tempPoint + "/" + "test_symlink"
	err = os.Symlink(oldPath, newPath)

	if err != nil {
		t.Fatalf("Failed to create symlink: %+v \n", err)
	}

	// readlink
	path, err := os.Readlink(newPath)
	if err != nil {
		t.Fatalf("Failed to readlink: %+v \n", err)
	}

	if path != oldPath {
		t.Fatalf("the link in stmlink not correct, [%s] != [%s]", oldPath, path)
	}

}
