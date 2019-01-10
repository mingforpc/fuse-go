package test

import (
	"io/ioutil"
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
