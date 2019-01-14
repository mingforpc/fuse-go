package test

import (
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

	// only for symlink
	link string
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

// the file for create
var newFile testFileStat

// the map to save xattr
// key: inode id
// value: map[{xattr name}]{xattr value}
var xattrMap map[uint64]map[string]string

// the symlink File
var symlinkFile testFileStat

// the hard link file
var hardlinkFile testFileStat

// the statfs of root
var rootStatfs fuse.Statfs

func init() {
	// root
	stat := fuse.FileStat{}
	stat.Nodeid = 1
	stat.Stat.Ino = 1
	stat.Stat.Mode = uint32(syscall.S_IFDIR) | uint32(0777)
	stat.Stat.Nlink = 2
	root.stat = stat

	// rootFile
	stat = fuse.FileStat{}
	stat.Nodeid = 2
	stat.Stat.Ino = 2
	stat.Stat.Mode = uint32(syscall.S_IFREG) | uint32(0777)
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

	// rootStatfs
	rootStatfs.Bsize = 1024
	rootStatfs.Blocks = 2048
	rootStatfs.Bfree = 2048
	rootStatfs.Bavail = 256
	rootStatfs.Files = 8
	rootStatfs.Ffree = 1024
	rootStatfs.NameLen = 256
	rootStatfs.Frsize = 256
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
	case 5:
		// newDir
		stat = &newDir.stat
	case 6:
		stat = &symlinkFile.stat
	case 7:
		stat = &newFile.stat
	case 8:
		stat = &hardlinkFile.stat

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
		} else if name == newFile.name {
			fsStat = &newFile.stat
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

	if nodeid != rootFile.stat.Nodeid {
		result = errno.EACCES
	} else {
		result = errno.SUCCESS
	}

	return result
}

var read = func(req fuse.Req, nodeid uint64, size uint32, offset uint64, fi fuse.FileInfo) (content []byte, result int32) {
	fmt.Printf("Read: nodeid:%d, size[%d], offset[%d], fi[%+v] \n", nodeid, size, offset, fi)

	if nodeid != rootFile.stat.Nodeid {
		result = errno.EACCES
	} else {
		content = []byte(rootFile.content)
		result = errno.SUCCESS
	}

	return content, result
}

var write = func(req fuse.Req, nodeid uint64, buf []byte, offset uint64, fi fuse.FileInfo) (size uint32, result int32) {
	fmt.Printf("Write: nodeid:%d, len of buf[%d], offset[%d], fi[%+v] \n", nodeid, len(buf), offset, fi)

	if nodeid != rootFile.stat.Nodeid {
		result = errno.EACCES
	} else {

		data := string(buf)
		dataLen := uint64(len(data))
		contentLen := uint64(len(rootFile.content))
		if offset >= contentLen {
			rootFile.content += data
		} else {

			if offset+dataLen < contentLen {
				rootFile.content = rootFile.content[:offset] + data + rootFile.content[offset+dataLen:]
			} else {
				rootFile.content = rootFile.content[:offset] + data
			}
		}

		size = uint32(dataLen)
		result = errno.SUCCESS
	}

	return size, result
}

var fsync = func(req fuse.Req, nodeid uint64, datasync uint32, fi fuse.FileInfo) (result int32) {

	fmt.Printf("Fsync: nodeid:%d, datasync[%d],fi[%+v] \n", nodeid, datasync, fi)

	if nodeid != rootFile.stat.Nodeid {
		result = errno.EACCES
	} else {
		result = errno.SUCCESS
	}

	return result
}

var flush = func(req fuse.Req, nodeid uint64, fi fuse.FileInfo) (result int32) {
	fmt.Printf("Flush: nodeid:%d,  fi[%+v] \n", nodeid, fi)
	return errno.SUCCESS
}
var release = func(req fuse.Req, nodeid uint64, fi fuse.FileInfo) (result int32) {
	fmt.Printf("Release: nodeid:[%d],  fi:[%+b] \n", nodeid, fi)
	return result
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

var symlink = func(req fuse.Req, parentid uint64, link string, name string) (fsStat *fuse.FileStat, result int32) {

	fmt.Printf("Symlink: parentid[%d], link[%s], name[%s] \n", parentid, link, name)

	symlinkFile.name = name
	symlinkFile.link = link
	symlinkFile.stat.Nodeid = 6
	symlinkFile.stat.Stat.Ino = 6
	symlinkFile.stat.Stat.Mode = uint32(syscall.S_IFLNK) | uint32(0444)

	result = errno.SUCCESS

	return &symlinkFile.stat, result
}

var readlink = func(req fuse.Req, nodeid uint64) (path string, result int32) {

	fmt.Printf("Readlink: nodeid[%d] \n", nodeid)

	if symlinkFile.stat.Nodeid != nodeid {
		return "", errno.EINVAL
	}

	path = symlinkFile.link
	result = errno.SUCCESS

	return path, result
}

var mknod = func(req fuse.Req, parentid uint64, name string, mode uint32, rdev uint32) (fsStat *fuse.FileStat, result int32) {

	fmt.Printf("Mknod: parentid:%d, name:%s, mode:%x, rdev:%d \n", parentid, name, mode, rdev)

	newFile.name = name
	newFile.stat.Nodeid = 7
	fsStat = &newFile.stat

	fsStat.Nodeid = 7
	fsStat.Stat.Ino = 7
	fsStat.Stat.Mode = mode
	fsStat.Stat.Rdev = uint64(rdev)

	return fsStat, result
}

var unlink = func(req fuse.Req, parentid uint64, name string) (result int32) {
	fmt.Printf("Unlink: parentid:%d, name:%s \n", parentid, name)

	result = errno.SUCCESS
	if parentid == root.stat.Nodeid {

		if name == newFile.name {
			newFile.name = ""
			newFile.stat = fuse.FileStat{}
		} else {
			result = errno.EACCES
		}

	} else if parentid == rootDir.stat.Nodeid {

		result = errno.EACCES
	} else {
		result = errno.ENOENT
	}

	return result
}

var rename = func(req fuse.Req, parentid uint64, name string, newparentid uint64, newname string) (result int32) {
	fmt.Printf("Rename: parentid[%d], name[%s], newparentid[%d], newname[%s] \n", parentid, name, newparentid, newname)

	if parentid == root.stat.Nodeid && name == newFile.name {

		if newparentid == root.stat.Nodeid {

			newFile.name = newname

		} else {
			result = errno.EACCES
		}

	} else {
		result = errno.EACCES
	}

	return result
}

var link = func(req fuse.Req, oldnodeid uint64, newparentid uint64, newname string) (fsStat *fuse.FileStat, result int32) {

	fmt.Printf("Link: oldnodeid[%d], newparentid[%d], newname[%s] \n", oldnodeid, newparentid, newname)

	if oldnodeid != rootFile.stat.Nodeid && newparentid != root.stat.Nodeid {
		result = errno.EACCES
	} else {

		hardlinkFile.name = newname
		hardlinkFile.stat = rootFile.stat
		fsStat = &hardlinkFile.stat
		result = errno.SUCCESS
	}

	return fsStat, result
}

var statfs = func(req fuse.Req, nodeid uint64) (statfs *fuse.Statfs, result int32) {

	fmt.Printf("Statfs: nodeid[%d] \n", nodeid)

	if nodeid == root.stat.Nodeid {
		statfs = &rootStatfs
		result = errno.SUCCESS
	} else {
		result = errno.EACCES
	}

	return statfs, result
}

var access = func(req fuse.Req, nodeid uint64, mask uint32) (result int32) {
	fmt.Printf("Access: nodeid[%d], mask[%d] \n", nodeid, mask)

	if nodeid != rootFile.stat.Nodeid {
		result = errno.SUCCESS
	} else {
		result = errno.EACCES
	}

	return result
}

var create = func(req fuse.Req, parentid uint64, name string, mode uint32, fi *fuse.FileInfo) (fsStat *fuse.FileStat, result int32) {
	fmt.Printf("Create: parentid:%d, name:%s, mode:%x, fi:%+v \n", parentid, name, mode, fi)

	newFile.name = name
	newFile.stat.Nodeid = 7
	fsStat = &newFile.stat

	fsStat.Nodeid = 7
	fsStat.Stat.Ino = 7
	fsStat.Stat.Mode = mode

	return fsStat, result
}

var getlk = func(req fuse.Req, nodeid uint64, fi fuse.FileInfo, lock *fuse.Flock) (result int32) {

	fmt.Printf("Getlk: nodeid:%d, fi:%+v, lock:%+v \n", nodeid, fi, lock)

	if nodeid != rootFile.stat.Nodeid {
		return errno.EACCES
	}

	return errno.SUCCESS
}

var setlk = func(req fuse.Req, nodeid uint64, fi fuse.FileInfo, lock fuse.Flock, lksleep int) (result int32) {

	fmt.Printf("Setlk: nodeid:%d, fi:%+v, lock:%+v, lksleep:%d \n", nodeid, fi, lock, lksleep)

	if nodeid != rootFile.stat.Nodeid {
		return errno.EACCES
	}

	return errno.SUCCESS

}

// NewTestFuse : create a fuse session for test
func NewTestFuse(mountpoint string, opts fuse.Opt) *fuse.Session {
	if opts.Init == nil {
		opts.Init = &testInit
	}

	se := fuse.NewFuseSession(mountpoint, &opts, 1024)
	se.Debug = true
	se.FuseConfig.AttrTimeout = 1

	return se
}
