package test

import (
	"errors"
	"sync"
	"syscall"

	"github.com/mingforpc/fuse-go/fuse"
	"github.com/mingforpc/fuse-go/fuse/errno"
)

var helloName = "hello"
var helloContent = "hello world!\n"

//Se fuse session
var Se *fuse.Session

var wait sync.WaitGroup

var testInit = func(conn *fuse.ConnInfo) (userdata interface{}) {

	wait.Done()

	return nil
}

var readdir = func(req fuse.Req, nodeid uint64, size uint32, offset uint64, fi fuse.FileInfo) (fileList []fuse.Dirent, result int32) {

	if nodeid != 1 {
		result = errno.ENOTDIR
	} else if offset < 3 {

		current := fuse.Dirent{NameLen: uint32(len(".")), Ino: nodeid, Off: 0, Name: "."}
		prev := fuse.Dirent{NameLen: uint32(len("..")), Ino: nodeid, Off: 0, Name: ".."}
		hello := fuse.Dirent{NameLen: uint32(len(helloName)), Ino: 2, Off: 0, Name: helloName}

		fileList = make([]fuse.Dirent, 3)
		fileList[0] = current
		fileList[1] = prev
		fileList[2] = hello

	}
	result = errno.SUCCESS
	return fileList, result
}

var open = func(req fuse.Req, nodeid uint64, fi *fuse.FileInfo) (result int32) {

	if nodeid != 2 {
		result = errno.EISDIR
	} else if (fi.Flags & 3) != syscall.O_RDONLY {
		result = errno.EACCES
	} else {
		result = errno.SUCCESS
	}

	return result
}

var read = func(req fuse.Req, nodeid uint64, size uint32, offset uint64, fi fuse.FileInfo) (content []byte, result int32) {

	if nodeid != 2 {
		panic(errors.New("read error file"))
	}

	result = errno.SUCCESS

	if uint32(len(helloContent)) < size {

		contentBuf := []byte(helloContent)

		content = contentBuf[offset:len(helloContent)]

	} else {
		contentBuf := []byte(helloContent)
		content = contentBuf[offset:size]
	}

	return content, result
}

// NewTestFuse : create a fuse session for test
func NewTestFuse(mountpoint string, opts fuse.Opt) *fuse.Session {
	opts.Init = &testInit

	se := fuse.NewFuseSession(mountpoint, &opts, 1024)
	se.Debug = true
	se.FuseConfig.AttrTimeout = 1

	return se
}
