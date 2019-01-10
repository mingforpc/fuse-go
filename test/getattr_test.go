package test

import (
	"syscall"
	"testing"

	"github.com/mingforpc/fuse-go/fuse"
)

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

	// {root}/test
	var rootFileStat syscall.Stat_t
	err = syscall.Stat(tempPoint+"/"+rootFile.name, &rootFileStat)
	if err != nil {
		t.Errorf("TestLookup err: %+v \n", err)
	}
	stat := rootFile.stat.Stat
	if rootFileStat.Ino != stat.Ino {
		t.Errorf("getattr_test inode should be %d \n", stat.Ino)
	}
	if rootFileStat.Mode != stat.Mode {
		t.Errorf("getattr_test mode should be %x \n", stat.Mode)
	}
	if rootFileStat.Nlink != stat.Nlink {
		t.Errorf("getattr_test Nlink should be %d \n", 1)
	}
	if rootFileStat.Size != stat.Size {
		t.Errorf("getattr_test Size should be %d \n", stat.Size)
	}
	if rootFileStat.Ctim.Sec != stat.Ctim.Sec && rootFileStat.Ctim.Nsec != stat.Ctim.Nsec {
		t.Errorf("getattr_test Ctim should be {%d, %d} \n", stat.Ctim.Sec, stat.Ctim.Nsec)
	}
	if rootFileStat.Atim.Sec != stat.Atim.Sec && rootFileStat.Atim.Nsec != stat.Atim.Nsec {
		t.Errorf("getattr_test Atim should be {%d, %d} \n", stat.Atim.Sec, stat.Atim.Nsec)
	}
	if rootFileStat.Mtim.Sec != stat.Mtim.Sec && rootFileStat.Mtim.Nsec != stat.Mtim.Nsec {
		t.Errorf("getattr_test Mtim should be {%d, %d} \n", stat.Mtim.Sec, stat.Mtim.Nsec)
	}
	if rootFileStat.Uid != stat.Uid {
		t.Errorf("getattr_test Uid should be %d \n", stat.Uid)
	}
	if rootFileStat.Gid != stat.Gid {
		t.Errorf("getattr_test Gid should be %d \n", stat.Gid)
	}

	// {root}/test_dir
	var rootDirStat syscall.Stat_t
	err = syscall.Stat(tempPoint+"/"+rootDir.name, &rootDirStat)
	if err != nil {
		t.Fatalf("TestLookup err: %+v \n", err)
	}
	stat = rootDir.stat.Stat
	if rootDirStat.Ino != stat.Ino {
		t.Errorf("getattr_test_dir inode should be %d \n", stat.Ino)
	}
	if rootDirStat.Mode != stat.Mode {
		t.Errorf("getattr_test_dir mode should be %x \n", stat.Mode)
	}
	if rootDirStat.Nlink != stat.Nlink {
		t.Errorf("getattr_test_dir Nlink should be %d \n", stat.Nlink)
	}
	if rootDirStat.Uid != stat.Uid {
		t.Errorf("getattr_test_dir Uid should be %d \n", stat.Uid)
	}
	if rootDirStat.Gid != stat.Gid {
		t.Errorf("getattr_test_dir Gid should be %d \n", stat.Gid)
	}

	// {root}/test_dir/test
	var dirFileStat syscall.Stat_t
	err = syscall.Stat(tempPoint+"/"+dirFile.path, &dirFileStat)
	if err != nil {
		t.Fatalf("TestLookup err: %+v \n", err)
	}
	stat = dirFile.stat.Stat
	if dirFileStat.Ino != stat.Ino {
		t.Errorf("getattr_test_dir/test inode should be %d \n", stat.Ino)
	}
	if dirFileStat.Mode != stat.Mode {
		t.Errorf("getattr_test_dir/test mode should be %x \n", stat.Mode)
	}
}

func TestGetattr(t *testing.T) {
	tempPoint, err := createTempPoint()

	if err != nil {
		t.Fatalf("TestGetattr err: %+v \n", err)
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
	stat := root.stat.Stat
	if rootStat.Ino != stat.Ino {
		t.Errorf("{root} inode should be %d \n", stat.Ino)
	}
	if rootStat.Mode != stat.Mode {
		t.Errorf("{root} mode should be %x \n", stat.Mode)
	}
	if rootStat.Nlink != stat.Nlink {
		t.Errorf("{root} Nlink should be %d \n", stat.Nlink)
	}

}
