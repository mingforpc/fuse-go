package test

import (
	"errors"
	"syscall"

	"github.com/mingforpc/fuse-go/fuse"
	"github.com/mingforpc/fuse-go/fuse/errno"
)

var se *fuse.Session

func init() {
	se = NewTestFuse("./test_mountpoint")
}

var helloName = "hello"
var helloContent = "hello world!\n"

func helloStat(nodeid uint64) (stat *fuse.FileStat) {

	switch nodeid {
	case 1:
		stat = new(fuse.FileStat)
		stat.Nodeid = 1
		stat.Stat.Mode = uint32(syscall.S_IFDIR) | uint32(0755)
		stat.Stat.Nlink = 2
	case 2:
		stat = new(fuse.FileStat)
		stat.Nodeid = 2
		stat.Stat.Mode = uint32(syscall.S_IFREG) | uint32(0444)
		stat.Stat.Nlink = 1
		stat.Stat.Size = int64(len(helloContent))
	default:
	}

	return stat
}

var getattr = func(req fuse.Req, nodeid uint64) (fsStat *fuse.FileStat, result int32) {

	fsStat = helloStat(nodeid)
	if fsStat == nil {
		result = errno.ENOENT
	} else {
		result = errno.SUCCESS
	}

	return fsStat, result
}

var lookup = func(req fuse.Req, parentId uint64, name string) (fsStat *fuse.FileStat, result int32) {

	if parentId != 1 || name != helloName {
		result = errno.ENOENT
	} else {
		fsStat = helloStat(2)
	}

	return fsStat, result
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
func NewTestFuse(mountpoint string) *fuse.Session {
	opts := fuse.Opt{}
	opts.Getattr = &getattr
	opts.Readdir = &readdir
	opts.Lookup = &lookup
	opts.Open = &open
	opts.Read = &read

	se := fuse.NewFuseSession(mountpoint, &opts, 1024)
	se.Debug = true
	se.FuseConfig.AttrTimeout = 1

	return se
}
