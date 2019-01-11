package test

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"syscall"

	"github.com/mingforpc/fuse-go/fuse"
	"github.com/mingforpc/fuse-go/fuse/errno"
)

//Se fuse session
var Se *fuse.Session

var wait sync.WaitGroup

type testFileStat struct {
	name    string
	content string
	path    string
	stat    fuse.FileStat
}

// root stat
var root testFileStat

// '/test' stat
var rootFile testFileStat

// '/dir' stat
var rootDir testFileStat

// '/dir/test' stat
var dirFile testFileStat

// the dir for create
var newDir testFileStat

// the map to save xattr
// key: inode id
// value: map[{xattr name}]{xattr value}
var xattrMap map[uint64]map[string]string

func init() {
	// root
	stat := fuse.FileStat{}
	stat.Nodeid = 1
	stat.Stat.Ino = 1
	stat.Stat.Mode = uint32(syscall.S_IFDIR) | uint32(0777)
	stat.Stat.Nlink = 2
	// root = testFileStat{}
	root.stat = stat

	// rootFile
	stat = fuse.FileStat{}
	stat.Nodeid = 2
	stat.Stat.Ino = 2
	stat.Stat.Mode = uint32(syscall.S_IFREG) | uint32(0444)
	stat.Stat.Nlink = 1
	stat.Stat.Size = int64(len("hello world!\n"))
	stat.Stat.Ctim = syscall.Timespec{1547044000, 100}
	stat.Stat.Atim = syscall.Timespec{1547044001, 200}
	stat.Stat.Mtim = syscall.Timespec{1547044002, 300}
	stat.Stat.Uid = 1
	stat.Stat.Gid = 1
	// rootFile = testFileStat{}
	rootFile.name = "test"
	rootFile.path = "test"
	rootFile.content = "hello world!\n"
	rootFile.stat = stat

	// rootDir
	stat = fuse.FileStat{}
	stat.Nodeid = 3
	stat.Stat.Ino = 3
	stat.Stat.Mode = uint32(syscall.S_IFDIR) | uint32(0755)
	stat.Stat.Nlink = 2
	stat.Stat.Uid = 2
	stat.Stat.Gid = 2
	// rootDir = testFileStat{}
	rootDir.name = "test_dir"
	rootDir.path = "test_dir"
	rootDir.stat = stat

	// dirFile
	stat = fuse.FileStat{}
	stat.Nodeid = 4
	stat.Stat.Ino = 4
	stat.Stat.Mode = uint32(syscall.S_IFREG) | uint32(0666)
	stat.Stat.Nlink = 1
	stat.Stat.Size = int64(len("hello world!\n"))
	dirFile.name = "test"
	dirFile.path = "test_dir/test"
	dirFile.content = "hello world!\n"
	dirFile.stat = stat

	// newDir
	stat = fuse.FileStat{}
	stat.Nodeid = 5
	stat.Stat.Ino = 5
	stat.Stat.Mode = uint32(syscall.S_IFDIR) | uint32(0755)
	stat.Stat.Nlink = 0
	stat.Stat.Uid = uint32(os.Getuid())
	stat.Stat.Gid = uint32(os.Getgid())
	newDir.stat = stat

	// init xattrmap
	xattrMap = make(map[uint64]map[string]string)
}

func getStat(nodeid uint64) (stat *fuse.FileStat) {

	switch nodeid {
	case 1:
		// root dir
		stat = &root.stat
	case 2:
		// /test file
		stat = &rootFile.stat

	case 3:
		// /test_dir dir
		stat = &rootDir.stat

	case 4:
		// test_dir/test file
		stat = &dirFile.stat

	default:
	}

	return stat
}

var testInit = func(conn *fuse.ConnInfo) (userdata interface{}) {

	wait.Done()

	return nil
}

var getattr = func(req fuse.Req, nodeid uint64) (fsStat *fuse.FileStat, result int32) {

	fmt.Printf("Getattr: nodeid:%d \n", nodeid)

	fsStat = getStat(nodeid)
	if fsStat == nil {
		result = errno.ENOENT
	} else {
		result = errno.SUCCESS
	}

	return fsStat, result
}

var lookup = func(req fuse.Req, parentId uint64, name string) (fsStat *fuse.FileStat, result int32) {

	fmt.Printf("Looup: parentid:%d, name:%s \n", parentId, name)

	// parent id not root folder or rootDir forlder
	if parentId != root.stat.Nodeid && parentId != rootDir.stat.Nodeid {
		result = errno.ENOENT
	} else if parentId == root.stat.Nodeid {
		// root folder

		if name == rootFile.name {
			fsStat = &rootFile.stat
		} else if name == rootDir.name {
			fsStat = &rootDir.stat
		} else {
			result = errno.ENOENT
		}

	} else {
		// rootDir folder

		if name == dirFile.name {
			fsStat = &dirFile.stat
		} else {
			result = errno.ENOENT
		}
	}

	return fsStat, result
}

var readdir = func(req fuse.Req, nodeid uint64, size uint32, offset uint64, fi fuse.FileInfo) (fileList []fuse.Dirent, result int32) {

	fmt.Printf("Readdir: nodeid:%d, size:%d offset:%d, fi:[%+v] \n", nodeid, size, offset, fi)

	if nodeid != root.stat.Nodeid && nodeid != rootDir.stat.Nodeid {
		result = errno.ENOTDIR
	}

	current := fuse.Dirent{NameLen: uint32(len(".")), Ino: nodeid, Off: 0, Name: "."}
	prev := fuse.Dirent{NameLen: uint32(len("..")), Ino: nodeid, Off: 0, Name: ".."}
	if nodeid == root.stat.Nodeid && offset < 4 {

		test := fuse.Dirent{NameLen: uint32(len(rootFile.name)), Ino: rootFile.stat.Stat.Ino, Off: 0, Name: rootFile.name}
		dir := fuse.Dirent{NameLen: uint32(len(rootDir.name)), Ino: rootDir.stat.Stat.Ino, Off: 0, Name: rootDir.name}

		fileList = make([]fuse.Dirent, 4)
		fileList[0] = current
		fileList[1] = prev
		fileList[2] = test
		fileList[3] = dir
	} else if nodeid == rootDir.stat.Nodeid && offset < 3 {

		test := fuse.Dirent{NameLen: uint32(len(dirFile.name)), Ino: dirFile.stat.Stat.Ino, Off: 0, Name: dirFile.name}

		fileList = make([]fuse.Dirent, 3)
		fileList[0] = current
		fileList[1] = prev
		fileList[2] = test
	}

	result = errno.SUCCESS
	return fileList, result
}

var open = func(req fuse.Req, nodeid uint64, fi *fuse.FileInfo) (result int32) {

	fmt.Printf("Open: nodeid:%d,  fi:[%+v] \n", nodeid, fi)

	return result
}

var read = func(req fuse.Req, nodeid uint64, size uint32, offset uint64, fi fuse.FileInfo) (content []byte, result int32) {

	if nodeid != 2 {
		panic(errors.New("read error file"))
	}

	result = errno.SUCCESS

	return content, result
}

var fsyncdir = func(req fuse.Req, nodeid uint64, datasync uint32, fi fuse.FileInfo) (result int32) {

	fmt.Printf("Open: nodeid:[%d], datasync:[%d], fi:[%+b] \n", nodeid, datasync, fi)

	return result
}

var mkdir = func(req fuse.Req, parentid uint64, name string, mode uint32) (fsStat *fuse.FileStat, result int32) {

	fmt.Printf("Mkdir: parentid:[%d], name:[%s], mode:[%d] \n", parentid, name, mode)

	fsStat = &newDir.stat
	newDir.name = name

	result = errno.SUCCESS

	return fsStat, result
}

var rmdir = func(req fuse.Req, parentid uint64, name string) (res int32) {
	fmt.Printf("Rmdir: parentid:[%d], name:[%s] \n", parentid, name)

	if parentid == root.stat.Nodeid && name == rootDir.name {
		return errno.ENOTEMPTY
	}

	return errno.SUCCESS
}

var setxattr = func(req fuse.Req, nodeid uint64, name string, value string, flags uint32) (result int32) {

	fmt.Printf("Setxattr: nodeid:[%d], name:[%s], value:[%s] \n", nodeid, name, value)

	kvMap := xattrMap[nodeid]

	if kvMap == nil {
		xattrMap[nodeid] = make(map[string]string)
		kvMap = xattrMap[nodeid]
	}

	kvMap[name] = value

	result = errno.SUCCESS

	return result
}

var getxattr = func(req fuse.Req, nodeid uint64, name string, size uint32) (value string, result int32) {

	fmt.Printf("Getxattr: nodeid:[%d], name:[%s] \n", nodeid, name)

	kvMap := xattrMap[nodeid]

	if kvMap == nil {
		return "", errno.ENOATTR
	}

	value = kvMap[name]

	if value == "" {
		result = errno.ENOATTR
	} else {
		result = errno.SUCCESS
	}

	return value, result
}

var listxattr = func(req fuse.Req, nodeid uint64, size uint32) (listVal string, result int32) {
	fmt.Printf("Listxattr: nodeid:[%d] \n", nodeid)
	kvMap := xattrMap[nodeid]

	if kvMap == nil {
		return "", errno.ENOATTR
	}

	for k := range kvMap {
		listVal += k + string(byte(0))
	}
	// remove lastest '\0'
	listVal = listVal[:len(listVal)-1]

	result = errno.SUCCESS

	return listVal, result
}

var removexattr = func(req fuse.Req, nodeid uint64, name string) (result int32) {
	fmt.Printf("Removexattr: nodeid:[%d], name[%s] \n", nodeid, name)

	kvMap := xattrMap[nodeid]

	if kvMap == nil {
		return errno.ENOATTR
	}

	if kvMap[name] == "" {
		return errno.ENOATTR
	}

	delete(kvMap, name)

	result = errno.SUCCESS
	return result
}

// NewTestFuse : create a fuse session for test
func NewTestFuse(mountpoint string, opts fuse.Opt) *fuse.Session {
	opts.Init = &testInit

	se := fuse.NewFuseSession(mountpoint, &opts, 1024)
	se.Debug = false
	se.FuseConfig.AttrTimeout = 1

	return se
}
