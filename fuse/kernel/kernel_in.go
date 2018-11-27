package kernel

import (
	"bytes"

	"github.com/mingforpc/fuse-go/fuse/common"
)

// Each query starts with a FuseInHeader
type FuseInHeader struct {
	Len     uint32
	Opcode  uint32
	Unique  uint64
	Nodeid  uint64
	Uid     uint32
	Gid     uint32
	Pid     uint32
	Padding uint32
}

func (header *FuseInHeader) ParseBinary(bcontent []byte) error {
	err := common.ParseBinary(bcontent, header)

	return err
}

// init
type FuseInitIn struct {
	Major        uint32
	Minor        uint32
	MaxReadahead uint32
	Flags        uint32
}

func (init *FuseInitIn) ParseBinary(bcontent []byte) error {
	err := common.ParseBinary(bcontent, init)

	return err
}

// getattr
type FuseGetattrIn struct {
	GetattrFlags uint32
	Dummy        uint32
	Fh           uint64
}

func (getattr *FuseGetattrIn) ParseBinary(bcontent []byte) error {
	err := common.ParseBinary(bcontent, getattr)

	return err
}

// lookup
type FuseLookupIn struct {
	Name string
}

func (lookup *FuseLookupIn) ParseBinary(bcontent []byte) error {
	length := len(bcontent)

	if length > 0 {
		// avoid '\0'
		lookup.Name = string(bcontent[:length-1])
	}

	return nil
}

// forget (should not sne any reply)
type FuseForgetIn struct {
	Nlookup uint64
}

func (forget *FuseForgetIn) ParseBinary(bcontent []byte) error {
	err := common.ParseBinary(bcontent, forget)

	return err
}

// setattr
type FuseSetattrIn struct {
	Valid     uint32
	Padding   uint32
	Fh        uint64
	Size      uint64
	LockOwner uint64
	Atime     uint64
	Mtime     uint64
	Ctime     uint64
	AtimeNsec uint32
	MtimeNsec uint32
	CtimeNsec uint32
	Mode      uint32
	Unused4   uint32
	Uid       uint32
	Gid       uint32
	Unused5   uint32
}

func (setattr *FuseSetattrIn) ParseBinary(bcontent []byte) error {
	err := common.ParseBinary(bcontent, setattr)

	return err
}

// mknod
type FuseMknodIn struct {
	Mode    uint32
	Rdev    uint32
	Umask   uint32
	Padding uint32

	Name string
}

func (mknod *FuseMknodIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 16 {
		return DataLenErr
	}

	common.ParseBinary(bcontent[0:4], &mknod.Mode)
	common.ParseBinary(bcontent[4:8], &mknod.Rdev)
	common.ParseBinary(bcontent[8:12], &mknod.Umask)
	common.ParseBinary(bcontent[12:16], &mknod.Padding)

	mknod.Name = string(bcontent[16 : length-1])

	return nil
}

// mkdir
type FuseMkdirIn struct {
	Mode  uint32
	Umask uint32

	Name string
}

func (mkdir *FuseMkdirIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 8 {
		return DataLenErr
	}

	common.ParseBinary(bcontent[0:4], &mkdir.Mode)
	common.ParseBinary(bcontent[4:8], &mkdir.Umask)

	mkdir.Name = string(bcontent[8 : length-1])

	return nil
}

// rmdir
type FuseRmdirIn struct {
	Path string
}

func (rmdir *FuseRmdirIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)
	rmdir.Path = string(bcontent[:length-1])
	return nil
}

// unlink
type FuseUnlinkIn struct {
	Path string
}

func (unlink *FuseUnlinkIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)
	unlink.Path = string(bcontent[:length-1])
	return nil
}

// symlink
type FuseSymlinkIn struct {
	Name string

	LinkName string
}

func (symlink *FuseSymlinkIn) ParseBinary(bcontent []byte) error {

	array := bytes.Split(bcontent, []byte{0})

	if len(array) < 2 {
		return DataLenErr
	}

	symlink.Name = string(array[0])
	symlink.LinkName = string(array[1])

	return nil
}

// rename
type FuseRenameIn struct {
	NewDir uint64

	OldName string

	NewName string
}

func (rename *FuseRenameIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 8 {
		return DataLenErr
	}

	err := common.ParseBinary(bcontent[:8], &rename.NewDir)

	if err != nil {
		return err
	}

	array := bytes.Split(bcontent[8:], []byte{0})

	if len(array) < 2 {
		return DataLenErr
	}

	rename.OldName = string(array[0])
	rename.NewName = string(array[1])

	return nil
}

// link
type FuseLinkIn struct {
	OldNodeid uint64

	NewName string
}

func (link *FuseLinkIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 8 {
		return DataLenErr
	}

	err := common.ParseBinary(bcontent[:8], &link.OldNodeid)

	if err != nil {
		return err
	}

	link.NewName = string(bcontent[8 : length-1])

	return nil
}

// open, opendir
type FuseOpenIn struct {
	Flags  uint32
	Unused uint32
}

func (open *FuseOpenIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 8 {
		return DataLenErr
	}

	err := common.ParseBinary(bcontent, open)

	return err
}

// read, readdir
type FuseReadIn struct {
	Fh        uint64
	Offset    uint64
	Size      uint32
	ReadFlags uint32
	LockOwner uint64
	Flags     uint32
	Padding   uint32
}

func (read *FuseReadIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 40 {
		return DataLenErr
	}

	err := common.ParseBinary(bcontent, read)

	return err
}

// write
type FuseWriteIn struct {
	Fh         uint64
	Offset     uint64
	Size       uint32
	WriteFlags uint32
	LockOwner  uint64
	Flags      uint32
	Padding    uint32

	Buf []byte
}

func (write *FuseWriteIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 40 {
		return DataLenErr
	}

	common.ParseBinary(bcontent[0:8], &write.Fh)
	common.ParseBinary(bcontent[8:16], &write.Offset)
	common.ParseBinary(bcontent[16:20], &write.Size)
	common.ParseBinary(bcontent[20:24], &write.WriteFlags)
	common.ParseBinary(bcontent[24:32], &write.LockOwner)
	common.ParseBinary(bcontent[32:36], &write.Flags)
	common.ParseBinary(bcontent[36:40], &write.Padding)

	write.Buf = bcontent[40:]

	return nil
}

// release, releasedir
type FuseReleaseIn struct {
	Fh           uint64
	Flags        uint32
	ReleaseFlags uint32
	LockOwner    uint64
}

func (release *FuseReleaseIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 24 {
		return DataLenErr
	}

	err := common.ParseBinary(bcontent, release)

	return err
}

// fsync, fsyncdir
type FuseFsyncIn struct {
	Fh         uint64
	FsyncFlags uint32
	Padding    uint32
}

func (fsync *FuseFsyncIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 16 {
		return DataLenErr
	}

	err := common.ParseBinary(bcontent, fsync)

	return err
}

// setxattr
type FuseSetxattrIn struct {
	Size  uint32
	Flags uint32

	Name  string
	Value string
}

func (setxattr *FuseSetxattrIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 8 {
		return DataLenErr
	}

	common.ParseBinary(bcontent[:4], &setxattr.Size)
	common.ParseBinary(bcontent[4:8], &setxattr.Flags)

	array := bytes.Split(bcontent[8:], []byte{0})

	if len(array) < 2 {
		return DataLenErr
	}

	setxattr.Name = string(array[0])
	setxattr.Value = string(array[1])

	return nil
}

// getxattr, listxattr
type FuseGetxattrIn struct {
	Size    uint32
	Padding uint32

	Name string
}

func (getxattr *FuseGetxattrIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 8 {
		return DataLenErr
	}

	common.ParseBinary(bcontent[:4], &getxattr.Size)
	common.ParseBinary(bcontent[4:8], &getxattr.Padding)

	getxattr.Name = string(bcontent[8 : length-1])

	return nil
}

// removexattr
type FuseRemovexattrIn struct {
	Name string
}

func (removexattr *FuseRemovexattrIn) ParseBinary(bcontent []byte) error {

	removexattr.Name = string(bcontent)

	return nil
}

// flush
type FuseFlushIn struct {
	Fh        uint64
	Unused    uint32
	Padding   uint32
	LockOwner uint64
}

func (flush *FuseFlushIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 24 {
		return DataLenErr
	}

	common.ParseBinary(bcontent, flush)

	return nil
}

// getlk, setlk, setlkw
type FuseLkIn struct {
	Fh      uint64
	Owner   uint64
	Lk      FuseFileLock
	LkFlags uint32
	Padding uint32
}

func (lk *FuseLkIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 48 {
		return DataLenErr
	}

	common.ParseBinary(bcontent, lk)

	return nil
}

// access
type FuseAccessIn struct {
	Mask    uint32
	Padding uint32
}

func (access *FuseAccessIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 8 {
		return DataLenErr
	}

	common.ParseBinary(bcontent, access)

	return nil
}

// create
type FuseCreateIn struct {
	Flags   uint32
	Mode    uint32
	Umask   uint32
	Padding uint32

	Name string
}

func (create *FuseCreateIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 16 {
		return DataLenErr
	}

	common.ParseBinary(bcontent[:4], &create.Flags)
	common.ParseBinary(bcontent[4:8], &create.Mode)
	common.ParseBinary(bcontent[8:12], &create.Umask)
	common.ParseBinary(bcontent[12:16], &create.Padding)

	// length-1 是为了避开最后一个\0字符
	create.Name = string(bcontent[16 : length-1])

	return nil
}

// interrupt
type FuseInterruptIn struct {
	Unique uint64
}

func (interrupt *FuseInterruptIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 8 {
		return DataLenErr
	}

	err := common.ParseBinary(bcontent, interrupt)

	return err
}

// bmap
type FuseBmapIn struct {
	Block     uint64
	BlockSize uint32
	Padding   uint32
}

func (bmap *FuseBmapIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 16 {
		return DataLenErr
	}

	err := common.ParseBinary(bcontent, bmap)

	return err
}

// ioctl
type FuseIoctlIn struct {
	Fh      uint64
	Flags   uint32
	Cmd     uint32
	Arg     uint64
	InSize  uint32
	OutSize uint32

	InBuf []byte
}

func (ioctl *FuseIoctlIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 32 {
		return DataLenErr
	}

	common.ParseBinary(bcontent[:8], &ioctl.Fh)
	common.ParseBinary(bcontent[8:12], &ioctl.Flags)
	common.ParseBinary(bcontent[12:16], &ioctl.Cmd)
	common.ParseBinary(bcontent[16:24], &ioctl.Arg)
	common.ParseBinary(bcontent[24:28], &ioctl.InSize)
	common.ParseBinary(bcontent[28:32], &ioctl.OutSize)

	ioctl.InBuf = bcontent[32:]

	return nil
}

// poll
type FusePollIn struct {
	Fh     uint64
	Kh     uint64
	Flags  uint32
	Events uint32
}

func (poll *FusePollIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 24 {
		return DataLenErr
	}

	err := common.ParseBinary(bcontent, poll)

	return err
}

type FuseForgetOne struct {
	Nodeid  uint64
	Nlookup uint64
}

// batch_forget
type FuseBatchForgetIn struct {
	Count uint32
	Dummy uint32

	NodeList []FuseForgetOne
}

func (forget *FuseBatchForgetIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 8 {
		return DataLenErr
	}

	common.ParseBinary(bcontent[4:8], &forget.Dummy)

	forget.NodeList = make([]FuseForgetOne, forget.Count)
	var i uint32
	for i = 0; i < forget.Count; i++ {

		var temp = FuseForgetOne{}
		common.ParseBinary(bcontent[8*(i+1):16*(i+1)], &temp)
		forget.NodeList[i] = temp
	}

	return nil
}

// fallocate
type FuseFallocateIn struct {
	Fh      uint64
	Offset  uint64
	Length  uint64
	Mode    uint32
	Padding uint32
}

func (fallocate *FuseFallocateIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 32 {
		return DataLenErr
	}

	err := common.ParseBinary(bcontent, fallocate)

	return err
}

// rename2
type FuseRename2In struct {
	NewDir  uint64
	Flags   uint32
	Padding uint32

	OldName string
	NewName string
}

func (rename *FuseRename2In) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 16 {
		return DataLenErr
	}

	common.ParseBinary(bcontent[:8], &rename.NewDir)
	common.ParseBinary(bcontent[8:12], &rename.Flags)
	common.ParseBinary(bcontent[12:16], &rename.Padding)

	array := bytes.Split(bcontent[16:], []byte{0})

	if len(array) < 2 {
		return DataLenErr
	}

	rename.OldName = string(array[0])
	rename.NewName = string(array[1])

	return nil
}

// lseek
type FuseLseekIn struct {
	Fh      uint64
	Offset  uint64
	Whence  uint32
	Padding uint32
}

func (lseek *FuseLseekIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 24 {
		return DataLenErr
	}

	err := common.ParseBinary(bcontent, lseek)

	return err
}

// cuse_init
type CuseInitIn struct {
	Major  uint32
	Minor  uint32
	Unused uint32
	Flags  uint32
}

func (cuseInit *CuseInitIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 16 {
		return DataLenErr
	}

	err := common.ParseBinary(bcontent, cuseInit)

	return err
}
