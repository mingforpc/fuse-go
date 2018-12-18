package kernel

import (
	"errors"
)

var DataLenErr = errors.New("Data len not enough")
var NoNeedReplyErr = errors.New("No need to reply")
var NotInit = errors.New("Fuse session not inited")

/** Version number of this interface */
const FUSE_KERNEL_VERSION = 7

/** Minor version number of this interface */
const FUSE_KERNEL_MINOR_VERSION = 26

/**
 * Lock flags
 */
const FUSE_LK_FLOCK = (1 << 0)

/**
 * Ioctl flags
 *
 * FUSE_IOCTL_COMPAT: 32bit compat ioctl on 64bit machine
 * FUSE_IOCTL_UNRESTRICTED: not restricted to well-formed ioctls, retry allowed
 * FUSE_IOCTL_RETRY: retry with new iovecs
 * FUSE_IOCTL_32BIT: 32bit ioctl
 * FUSE_IOCTL_DIR: is a directory
 *
 * FUSE_IOCTL_MAX_IOV: maximum of in_iovecs + out_iovecs
 */
const FUSE_IOCTL_COMPAT = (1 << 0)
const FUSE_IOCTL_UNRESTRICTED = (1 << 1)
const FUSE_IOCTL_RETRY = (1 << 2)
const FUSE_IOCTL_32BIT = (1 << 3)
const FUSE_IOCTL_DIR = (1 << 4)
const FUSE_IOCTL_MAX_IOV = 256

/**
 * Poll flags
 *
 * FUSE_POLL_SCHEDULE_NOTIFY: request poll notify
 */
const FUSE_POLL_SCHEDULE_NOTIFY = (1 << 0)

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
	Uid       uint32
	Gid       uint32
	Rdev      uint32
	Blksize   uint32
	Padding   uint32
}

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

type FuseFileLock struct {
	Start uint64
	End   uint64
	Type  uint32
	Pid   uint32 /* tgid */
}

type FluseIoctlIovec struct {
	base uint64
	len  uint64
}

const (
	FUSE_NOTIFY_POLL        = 1
	FUSE_NOTIFY_INVAL_INODE = 2
	FUSE_NOTIFY_INVAL_ENTRY = 3
	FUSE_NOTIFY_STORE       = 4
	FUSE_NOTIFY_RETRIEVE    = 5
	FUSE_NOTIFY_DELETE      = 6
	FUSE_NOTIFY_CODE_MAX
)

/* The read buffer is required to be at least 8k, but may be much larger */
const FUSE_MIN_READ_BUFFER = 8192

/**
 * INIT request/reply flags
 *
 * FUSE_ASYNC_READ: asynchronous read requests
 * FUSE_POSIX_LOCKS: remote locking for POSIX file locks
 * FUSE_FILE_OPS: kernel sends file handle for fstat, etc... (not yet supported)
 * FUSE_ATOMIC_O_TRUNC: handles the O_TRUNC open flag in the filesystem
 * FUSE_EXPORT_SUPPORT: filesystem handles lookups of "." and ".."
 * FUSE_BIG_WRITES: filesystem can handle write size larger than 4kB
 * FUSE_DONT_MASK: don't apply umask to file mode on create operations
 * FUSE_SPLICE_WRITE: kernel supports splice write on the device
 * FUSE_SPLICE_MOVE: kernel supports splice move on the device
 * FUSE_SPLICE_READ: kernel supports splice read on the device
 * FUSE_FLOCK_LOCKS: remote locking for BSD style file locks
 * FUSE_HAS_IOCTL_DIR: kernel supports ioctl on directories
 * FUSE_AUTO_INVAL_DATA: automatically invalidate cached pages
 * FUSE_DO_READDIRPLUS: do READDIRPLUS (READDIR+LOOKUP in one)
 * FUSE_READDIRPLUS_AUTO: adaptive readdirplus
 * FUSE_ASYNC_DIO: asynchronous direct I/O submission
 * FUSE_WRITEBACK_CACHE: use writeback cache for buffered writes
 * FUSE_NO_OPEN_SUPPORT: kernel supports zero-message opens
 * FUSE_PARALLEL_DIROPS: allow parallel lookups and readdir
 * FUSE_HANDLE_KILLPRIV: fs handles killing suid/sgid/cap on write/chown/trunc
 * FUSE_POSIX_ACL: filesystem supports posix acls
 */
const FUSE_ASYNC_READ = (1 << 0)
const FUSE_POSIX_LOCKS = (1 << 1)
const FUSE_FILE_OPS = (1 << 2)
const FUSE_ATOMIC_O_TRUNC = (1 << 3)
const FUSE_EXPORT_SUPPORT = (1 << 4)
const FUSE_BIG_WRITES = (1 << 5)
const FUSE_DONT_MASK = (1 << 6)
const FUSE_SPLICE_WRITE = (1 << 7)
const FUSE_SPLICE_MOVE = (1 << 8)
const FUSE_SPLICE_READ = (1 << 9)
const FUSE_FLOCK_LOCKS = (1 << 10)
const FUSE_HAS_IOCTL_DIR = (1 << 11)
const FUSE_AUTO_INVAL_DATA = (1 << 12)
const FUSE_DO_READDIRPLUS = (1 << 13)
const FUSE_READDIRPLUS_AUTO = (1 << 14)
const FUSE_ASYNC_DIO = (1 << 15)
const FUSE_WRITEBACK_CACHE = (1 << 16)
const FUSE_NO_OPEN_SUPPORT = (1 << 17)
const FUSE_PARALLEL_DIROPS = (1 << 18)
const FUSE_HANDLE_KILLPRIV = (1 << 19)
const FUSE_POSIX_ACL = (1 << 20)
