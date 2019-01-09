package test

import (
	"fmt"
	"syscall"
	"testing"

	"github.com/mingforpc/fuse-go/fuse"
	"github.com/mingforpc/fuse-go/fuse/errno"
)

func testStat(nodeid uint64) (stat *fuse.FileStat) {

	switch nodeid {
	case 1:
		// root dir
		stat = new(fuse.FileStat)
		stat.Nodeid = 1
		stat.Stat.Ino = 1
		stat.Stat.Mode = uint32(syscall.S_IFDIR) | uint32(0755)
		stat.Stat.Nlink = 2
	case 2:
		// getattr_test file
		stat = new(fuse.FileStat)
		stat.Nodeid = 2
		stat.Stat.Ino = 2
		stat.Stat.Mode = uint32(syscall.S_IFREG) | uint32(0444)
		stat.Stat.Nlink = 1
		stat.Stat.Size = int64(1024)
		stat.Stat.Ctim = syscall.Timespec{1547044000, 100}
		stat.Stat.Atim = syscall.Timespec{1547044001, 200}
		stat.Stat.Mtim = syscall.Timespec{1547044002, 300}
		stat.Stat.Uid = 1
		stat.Stat.Gid = 1

	case 3:
		// getattr_test_dir dir
		stat = new(fuse.FileStat)
		stat.Nodeid = 3
		stat.Stat.Ino = 3
		stat.Stat.Mode = uint32(syscall.S_IFDIR) | uint32(0755)
		stat.Stat.Nlink = 2
		stat.Stat.Uid = 2
		stat.Stat.Gid = 2

	case 4:
		// getattr_test_dir/test file
		stat = new(fuse.FileStat)
		stat.Nodeid = 4
		stat.Stat.Ino = 4
		stat.Stat.Mode = uint32(syscall.S_IFREG) | uint32(0444)
		stat.Stat.Nlink = 1
	default:
	}

	return stat
}

var getattr = func(req fuse.Req, nodeid uint64) (fsStat *fuse.FileStat, result int32) {

	fmt.Printf("Getattr: nodeid:%d \n", nodeid)

	fsStat = testStat(nodeid)
	if fsStat == nil {
		result = errno.ENOENT
	} else {
		result = errno.SUCCESS
	}

	return fsStat, result
}

var lookup = func(req fuse.Req, parentId uint64, name string) (fsStat *fuse.FileStat, result int32) {

	fmt.Printf("Looup: parentid:%d, name:%s \n", parentId, name)

	if (parentId != 1 && parentId != 3) || (name != "getattr_test" && name != "getattr_test_dir" && name != "test") {
		result = errno.ENOENT
	} else if name == "getattr_test" && parentId == 1 {
		fsStat = testStat(2)
	} else if name == "getattr_test_dir" && parentId == 1 {
		fsStat = testStat(3)
	} else if name == "test" && parentId == 3 {
		fsStat = testStat(4)
	}

	return fsStat, result
}

//TestLookup : test getattr() -> lookup file in fuse dir
func TestLookup(t *testing.T) {

	tempPoint, err := createTempPoint()

	if err != nil {
		t.Fatalf("TestLookup err: %+v \n", err)
	}

	opts := fuse.Opt{}
	opts.Getattr = &getattr
	opts.Lookup = &lookup

	se := NewTestFuse(tempPoint, opts)

	err = preTest(se)

	if err != nil {
		panic(err)
	}

	go se.FuseLoop()
	defer exitTest(se)

	wait.Wait()

	// {root}
	var rootStat syscall.Stat_t
	err = syscall.Stat(tempPoint, &rootStat)

	if rootStat.Ino != 1 {
		t.Errorf("{root} inode should be %d \n", 1)
	}
	if rootStat.Mode != uint32(syscall.S_IFDIR)|uint32(0755) {
		t.Errorf("{root} mode should be %x \n", uint32(syscall.S_IFDIR)|uint32(0755))
	}
	if rootStat.Nlink != 2 {
		t.Errorf("{root} Nlink should be %d \n", 1)
	}

	// {root}/getattr_test
	var rootFileStat syscall.Stat_t
	err = syscall.Stat(tempPoint+"/getattr_test", &rootFileStat)
	if err != nil {
		t.Errorf("TestLookup err: %+v \n", err)
	}
	if rootFileStat.Ino != 2 {
		t.Errorf("getattr_test inode should be %d \n", 2)
	}
	if rootFileStat.Mode != uint32(syscall.S_IFREG)|uint32(0444) {
		t.Errorf("getattr_test mode should be %x \n", uint32(syscall.S_IFREG)|uint32(0444))
	}
	if rootFileStat.Nlink != 1 {
		t.Errorf("getattr_test Nlink should be %d \n", 1)
	}
	if rootFileStat.Size != 1024 {
		t.Errorf("getattr_test Size should be %d \n", 1024)
	}
	if rootFileStat.Ctim.Sec != 1547044000 && rootFileStat.Ctim.Nsec != 100 {
		t.Errorf("getattr_test Ctim should be {%d, %d} \n", 1547044000, 100)
	}
	if rootFileStat.Atim.Sec != 1547044001 && rootFileStat.Atim.Nsec != 200 {
		t.Errorf("getattr_test Atim should be {%d, %d} \n", 1547044001, 200)
	}
	if rootFileStat.Mtim.Sec != 1547044002 && rootFileStat.Mtim.Nsec != 300 {
		t.Errorf("getattr_test Mtim should be {%d, %d} \n", 1547044002, 300)
	}
	if rootFileStat.Uid != 1 {
		t.Errorf("getattr_test Uid should be %d \n", 1)
	}
	if rootFileStat.Gid != 1 {
		t.Errorf("getattr_test Gid should be %d \n", 1)
	}

	// {root}/getattr_test_dir
	var rootDirStat syscall.Stat_t
	err = syscall.Stat(tempPoint+"/getattr_test_dir", &rootDirStat)
	if err != nil {
		t.Fatalf("TestLookup err: %+v \n", err)
	}

	if rootDirStat.Ino != 3 {
		t.Errorf("getattr_test_dir inode should be %d \n", 3)
	}
	if rootDirStat.Mode != uint32(syscall.S_IFDIR)|uint32(0755) {
		t.Errorf("getattr_test_dir mode should be %x \n", uint32(syscall.S_IFDIR)|uint32(0755))
	}
	if rootDirStat.Nlink != 2 {
		t.Errorf("getattr_test_dir Nlink should be %d \n", 1)
	}
	if rootDirStat.Uid != 2 {
		t.Errorf("getattr_test_dir Uid should be %d \n", 2)
	}
	if rootDirStat.Gid != 2 {
		t.Errorf("getattr_test_dir Gid should be %d \n", 2)
	}

	// {root}/getattr_test_dir/test
	var dirFileStat syscall.Stat_t
	err = syscall.Stat(tempPoint+"/getattr_test_dir/test", &dirFileStat)
	if err != nil {
		t.Fatalf("TestLookup err: %+v \n", err)
	}

	if dirFileStat.Ino != 4 {
		t.Errorf("getattr_test_dir/test inode should be %d \n", 4)
	}
	if dirFileStat.Mode != uint32(syscall.S_IFREG)|uint32(0444) {
		t.Errorf("getattr_test_dir/test mode should be %x \n", uint32(syscall.S_IFREG)|uint32(0444))
	}
}
