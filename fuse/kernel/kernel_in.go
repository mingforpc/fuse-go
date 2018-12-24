package kernel

import (
	"bytes"

	"github.com/mingforpc/fuse-go/fuse/common"
)

// FuseInHeader : each query starts with a FuseInHeader
type FuseInHeader struct {
	Len     uint32
	Opcode  uint32
	Unique  uint64
	Nodeid  uint64
	UID     uint32
	Gid     uint32
	Pid     uint32
	Padding uint32
}

// ParseBinary : Parse binary to FuseInHeader
func (header *FuseInHeader) ParseBinary(bcontent []byte) error {
	err := common.ParseBinary(bcontent, header)

	return err
}

// FuseInitIn : init request
type FuseInitIn struct {
	Major        uint32
	Minor        uint32
	MaxReadahead uint32
	Flags        uint32
}

// ParseBinary : Parse binary to FuseInitIn
func (init *FuseInitIn) ParseBinary(bcontent []byte) error {
	err := common.ParseBinary(bcontent, init)

	return err
}

// FuseGetattrIn : getattr request
type FuseGetattrIn struct {
	GetattrFlags uint32
	Dummy        uint32
	Fh           uint64
}

// ParseBinary : Parse binary to FuseGetattrIn
func (getattr *FuseGetattrIn) ParseBinary(bcontent []byte) error {
	err := common.ParseBinary(bcontent, getattr)

	return err
}

// FuseLookupIn : lookup request
type FuseLookupIn struct {
	Name string
}

// ParseBinary : Parse binary to FuseLoopIn
func (lookup *FuseLookupIn) ParseBinary(bcontent []byte) error {
	length := len(bcontent)

	if length > 0 {
		// avoid '\0'
		lookup.Name = string(bcontent[:length-1])
	}

	return nil
}

// FuseForgetIn : forget request (should not send any reply)
type FuseForgetIn struct {
	Nlookup uint64
}

// ParseBinary : Parse binary to FuseForgetIn
func (forget *FuseForgetIn) ParseBinary(bcontent []byte) error {
	err := common.ParseBinary(bcontent, forget)

	return err
}

// FuseSetattrIn : setattr request
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
	UID       uint32
	Gid       uint32
	Unused5   uint32
}

// ParseBinary : Parse FuseInHeader to binary
func (setattr *FuseSetattrIn) ParseBinary(bcontent []byte) error {
	err := common.ParseBinary(bcontent, setattr)

	return err
}

// FuseMknodIn : mknod request
type FuseMknodIn struct {
	Mode    uint32
	Rdev    uint32
	Umask   uint32
	Padding uint32

	Name string
}

// ParseBinary : Parse binary to FuseMknodIn
func (mknod *FuseMknodIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 16 {
		return ErrDataLen
	}

	common.ParseBinary(bcontent[0:4], &mknod.Mode)
	common.ParseBinary(bcontent[4:8], &mknod.Rdev)
	common.ParseBinary(bcontent[8:12], &mknod.Umask)
	common.ParseBinary(bcontent[12:16], &mknod.Padding)

	mknod.Name = string(bcontent[16 : length-1])

	return nil
}

// FuseMkdirIn : mkdir request
type FuseMkdirIn struct {
	Mode  uint32
	Umask uint32

	Name string
}

// ParseBinary : Parse binary to FuseMkdirIn
func (mkdir *FuseMkdirIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 8 {
		return ErrDataLen
	}

	common.ParseBinary(bcontent[0:4], &mkdir.Mode)
	common.ParseBinary(bcontent[4:8], &mkdir.Umask)

	mkdir.Name = string(bcontent[8 : length-1])

	return nil
}

// FuseRmdirIn : rmdir request
type FuseRmdirIn struct {
	Path string
}

// ParseBinary : Parse binary to FuseRmdirIn
func (rmdir *FuseRmdirIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)
	rmdir.Path = string(bcontent[:length-1])
	return nil
}

// FuseUnlinkIn : unlink request
type FuseUnlinkIn struct {
	Path string
}

// ParseBinary : Parse binary to FuseUnlinkIn
func (unlink *FuseUnlinkIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)
	unlink.Path = string(bcontent[:length-1])
	return nil
}

// FuseSymlinkIn : symlink request
type FuseSymlinkIn struct {
	Name string

	LinkName string
}

// ParseBinary : Parse binary to FuseSymlinkIn
func (symlink *FuseSymlinkIn) ParseBinary(bcontent []byte) error {

	array := bytes.Split(bcontent, []byte{0})

	if len(array) < 2 {
		return ErrDataLen
	}

	symlink.Name = string(array[0])
	symlink.LinkName = string(array[1])

	return nil
}

// FuseRenameIn : rename request
type FuseRenameIn struct {
	NewDir uint64

	OldName string

	NewName string
}

// ParseBinary : Parse binary to FuseRenameIn
func (rename *FuseRenameIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 8 {
		return ErrDataLen
	}

	err := common.ParseBinary(bcontent[:8], &rename.NewDir)

	if err != nil {
		return err
	}

	array := bytes.Split(bcontent[8:], []byte{0})

	if len(array) < 2 {
		return ErrDataLen
	}

	rename.OldName = string(array[0])
	rename.NewName = string(array[1])

	return nil
}

// FuseLinkIn : link request
type FuseLinkIn struct {
	OldNodeid uint64

	NewName string
}

// ParseBinary : Parse binary to FuseLinkIn
func (link *FuseLinkIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 8 {
		return ErrDataLen
	}

	err := common.ParseBinary(bcontent[:8], &link.OldNodeid)

	if err != nil {
		return err
	}

	link.NewName = string(bcontent[8 : length-1])

	return nil
}

// FuseOpenIn : open, opendir request
type FuseOpenIn struct {
	Flags  uint32
	Unused uint32
}

// ParseBinary : Parse binary to FuseOpenIn
func (open *FuseOpenIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 8 {
		return ErrDataLen
	}

	err := common.ParseBinary(bcontent, open)

	return err
}

// FuseReadIn : read, readdir request
type FuseReadIn struct {
	Fh        uint64
	Offset    uint64
	Size      uint32
	ReadFlags uint32
	LockOwner uint64
	Flags     uint32
	Padding   uint32
}

// ParseBinary : Parse binary to FuseReadIn
func (read *FuseReadIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 40 {
		return ErrDataLen
	}

	err := common.ParseBinary(bcontent, read)

	return err
}

// FuseWriteIn : write request
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

// ParseBinary : Parse binary to
func (write *FuseWriteIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 40 {
		return ErrDataLen
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

// FuseReleaseIn : release, releasedir request
type FuseReleaseIn struct {
	Fh           uint64
	Flags        uint32
	ReleaseFlags uint32
	LockOwner    uint64
}

// ParseBinary : Parse binary to FuseReleaseIn
func (release *FuseReleaseIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 24 {
		return ErrDataLen
	}

	err := common.ParseBinary(bcontent, release)

	return err
}

// FuseFsyncIn : fsync, fsyncdir request
type FuseFsyncIn struct {
	Fh         uint64
	FsyncFlags uint32
	Padding    uint32
}

// ParseBinary : Parse binary to FuseFsyncIn
func (fsync *FuseFsyncIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 16 {
		return ErrDataLen
	}

	err := common.ParseBinary(bcontent, fsync)

	return err
}

// FuseSetxattrIn : setxattr
type FuseSetxattrIn struct {
	Size  uint32
	Flags uint32

	Name  string
	Value string
}

// ParseBinary : Parse binary to FuseSetxattrIn
func (setxattr *FuseSetxattrIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 8 {
		return ErrDataLen
	}

	common.ParseBinary(bcontent[:4], &setxattr.Size)
	common.ParseBinary(bcontent[4:8], &setxattr.Flags)

	array := bytes.Split(bcontent[8:], []byte{0})

	if len(array) < 2 {
		return ErrDataLen
	}

	setxattr.Name = string(array[0])
	setxattr.Value = string(array[1])

	return nil
}

// FuseGetxattrIn : getxattr, listxattr request
type FuseGetxattrIn struct {
	Size    uint32
	Padding uint32

	Name string
}

// ParseBinary : Parse binary to FuseGetxattrIn
func (getxattr *FuseGetxattrIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 8 {
		return ErrDataLen
	}

	common.ParseBinary(bcontent[:4], &getxattr.Size)
	common.ParseBinary(bcontent[4:8], &getxattr.Padding)

	if length > 8 {
		getxattr.Name = string(bcontent[8 : length-1])
	}

	return nil
}

// FuseRemovexattrIn : removexattr request
type FuseRemovexattrIn struct {
	Name string
}

// ParseBinary : Parse binary to FuseRemovexattrIn
func (removexattr *FuseRemovexattrIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)
	if length > 0 {
		removexattr.Name = string(bcontent[:length-1])
	}

	return nil
}

// FuseFlushIn : flush request
type FuseFlushIn struct {
	Fh        uint64
	Unused    uint32
	Padding   uint32
	LockOwner uint64
}

// ParseBinary : Parse binary to FuseFlushIn
func (flush *FuseFlushIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 24 {
		return ErrDataLen
	}

	common.ParseBinary(bcontent, flush)

	return nil
}

// FuseLkIn : getlk, setlk, setlkw request
type FuseLkIn struct {
	Fh      uint64
	Owner   uint64
	Lk      FuseFileLock
	LkFlags uint32
	Padding uint32
}

// ParseBinary : Parse binary to FuseLkIn
func (lk *FuseLkIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 48 {
		return ErrDataLen
	}

	common.ParseBinary(bcontent, lk)

	return nil
}

// FuseAccessIn : access request
type FuseAccessIn struct {
	Mask    uint32
	Padding uint32
}

// ParseBinary : Parse binary to FuseAccessIn
func (access *FuseAccessIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 8 {
		return ErrDataLen
	}

	common.ParseBinary(bcontent, access)

	return nil
}

// FuseCreateIn : create request
type FuseCreateIn struct {
	Flags   uint32
	Mode    uint32
	Umask   uint32
	Padding uint32

	Name string
}

// ParseBinary : Parse binary to FuseCreateIn
func (create *FuseCreateIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 16 {
		return ErrDataLen
	}

	common.ParseBinary(bcontent[:4], &create.Flags)
	common.ParseBinary(bcontent[4:8], &create.Mode)
	common.ParseBinary(bcontent[8:12], &create.Umask)
	common.ParseBinary(bcontent[12:16], &create.Padding)

	// length-1 是为了避开最后一个'\0'字符
	create.Name = string(bcontent[16 : length-1])

	return nil
}

// FuseInterruptIn : interrupt request
type FuseInterruptIn struct {
	Unique uint64
}

// ParseBinary : Parse binary to FuseInterruptIn
func (interrupt *FuseInterruptIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 8 {
		return ErrDataLen
	}

	err := common.ParseBinary(bcontent, interrupt)

	return err
}

// FuseBmapIn : bmap request
type FuseBmapIn struct {
	Block     uint64
	BlockSize uint32
	Padding   uint32
}

// ParseBinary : Parse binary to FuseBmapIn
func (bmap *FuseBmapIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 16 {
		return ErrDataLen
	}

	err := common.ParseBinary(bcontent, bmap)

	return err
}

// FuseIoctlIn : ioctl request
type FuseIoctlIn struct {
	Fh      uint64
	Flags   uint32
	Cmd     uint32
	Arg     uint64
	InSize  uint32
	OutSize uint32

	InBuf []byte
}

// ParseBinary : Parse binary to FuseIoctlIn
func (ioctl *FuseIoctlIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 32 {
		return ErrDataLen
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

// FusePollIn : poll request
type FusePollIn struct {
	Fh     uint64
	Kh     uint64
	Flags  uint32
	Events uint32
}

// ParseBinary : Parse binary to FusePollIn
func (poll *FusePollIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 24 {
		return ErrDataLen
	}

	err := common.ParseBinary(bcontent, poll)

	return err
}

// FuseForgetOne : the inode to fotget
type FuseForgetOne struct {
	Nodeid  uint64
	Nlookup uint64
}

// FuseBatchForgetIn : batch_forget request
type FuseBatchForgetIn struct {
	Count uint32
	Dummy uint32

	NodeList []FuseForgetOne
}

// ParseBinary : Parse binary to FuseBatchForgetIn
func (forget *FuseBatchForgetIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 8 {
		return ErrDataLen
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

// FuseFallocateIn : fallocate request
type FuseFallocateIn struct {
	Fh      uint64
	Offset  uint64
	Length  uint64
	Mode    uint32
	Padding uint32
}

// ParseBinary : Parse binary to FuseFallocateIn
func (fallocate *FuseFallocateIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 32 {
		return ErrDataLen
	}

	err := common.ParseBinary(bcontent, fallocate)

	return err
}

// FuseRename2In : rename2 request
type FuseRename2In struct {
	NewDir  uint64
	Flags   uint32
	Padding uint32

	OldName string
	NewName string
}

// ParseBinary : Parse binary to FuseRename2In
func (rename *FuseRename2In) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 16 {
		return ErrDataLen
	}

	common.ParseBinary(bcontent[:8], &rename.NewDir)
	common.ParseBinary(bcontent[8:12], &rename.Flags)
	common.ParseBinary(bcontent[12:16], &rename.Padding)

	array := bytes.Split(bcontent[16:], []byte{0})

	if len(array) < 2 {
		return ErrDataLen
	}

	rename.OldName = string(array[0])
	rename.NewName = string(array[1])

	return nil
}

// FuseLseekIn : lseek request
type FuseLseekIn struct {
	Fh      uint64
	Offset  uint64
	Whence  uint32
	Padding uint32
}

// ParseBinary : Parse binary to FuseLseekIn
func (lseek *FuseLseekIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 24 {
		return ErrDataLen
	}

	err := common.ParseBinary(bcontent, lseek)

	return err
}

// CuseInitIn : cuse_init request
type CuseInitIn struct {
	Major  uint32
	Minor  uint32
	Unused uint32
	Flags  uint32
}

// ParseBinary : Parse binary to CuseInitIn
func (cuseInit *CuseInitIn) ParseBinary(bcontent []byte) error {

	length := len(bcontent)

	if length < 16 {
		return ErrDataLen
	}

	err := common.ParseBinary(bcontent, cuseInit)

	return err
}
