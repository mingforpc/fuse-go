package kernel

// FuseKernelVersion version number of this interface
const FuseKernelVersion = 7

// FuseKernelMinorVersion Minor version number of this interface
const FuseKernelMinorVersion = 26

// FuseAttr : the attr struct, for getattr and setattr
type FuseAttr struct {
	Ino       uint64
	Size      uint64
	Blocks    uint64
	Atime     uint64
	Mtime     uint64
	Ctime     uint64
	AtimeNsec uint32
	MtimeNsec uint32
	CtimeNsec uint32
	Mode      uint32
	Nlink     uint32
	UID       uint32
	GID       uint32
	Rdev      uint32
	Blksize   uint32
	Padding   uint32
}

// FuseStatfs : the fuse statfs struct
type FuseStatfs struct {
	Blocks  uint64
	Bfree   uint64
	Bavail  uint64
	Files   uint64
	Ffree   uint64
	Bsize   uint32
	NameLen uint32
	Frsize  uint32
	Padding uint32
	Spare   [6]uint32
}

// FuseFileLock : the fuse file lock struct
type FuseFileLock struct {
	Start uint64
	End   uint64
	Type  uint32
	Pid   uint32 /* tgid */
}

// FluseIoctlIovec : the fuse ioctl iovec struct
type FluseIoctlIovec struct {
	base uint64
	len  uint64
}

// FuseMinReadBuffer : The read buffer is required to be at least 8k, but may be much larger
const FuseMinReadBuffer = 8192
