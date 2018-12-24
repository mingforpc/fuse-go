package fuse

import (
	"os"
	"syscall"
	"time"

	"github.com/mingforpc/fuse-go/fuse/kernel"
)

const KERNEL_BUF_PAGES = 32

/* room needed in buffer to accommodate header */
const HEADER_SIZE = 0x1000

/**
 * Flags returned by the OPEN request
 *
 * FOPEN_DIRECT_IO: bypass page cache for this open file
 * FOPEN_KEEP_CACHE: don't invalidate the data cache on open
 * FOPEN_NONSEEKABLE: the file is not seekable
 */
const FOPEN_DIRECT_IO = (1 << 0)
const FOPEN_KEEP_CACHE = (1 << 1)
const FOPEN_NONSEEKABLE = (1 << 2)

/* 'to_set' flags in setattr */
const FUSE_SET_ATTR_MODE = (1 << 0)
const FUSE_SET_ATTR_UID = (1 << 1)
const FUSE_SET_ATTR_GID = (1 << 2)
const FUSE_SET_ATTR_SIZE = (1 << 3)
const FUSE_SET_ATTR_ATIME = (1 << 4)
const FUSE_SET_ATTR_MTIME = (1 << 5)
const FUSE_SET_ATTR_ATIME_NOW = (1 << 7)
const FUSE_SET_ATTR_MTIME_NOW = (1 << 8)
const FUSE_SET_ATTR_CTIME = (1 << 10)

/* XATTR set flag  */
const XATTR_CREATE = 0x1
const XATTR_REPLACE = 0x2

type FuseConnInfo struct {
	Major        uint32
	Minor        uint32
	MaxReadahead uint32

	/**
	 * Capability flags that the kernel supports (read-only)
	 */
	Capable uint32

	/**
	 * Capability flags that the filesystem wants to enable.
	 *
	 * libfuse attempts to initialize this field with
	 * reasonable default values before calling the init() handler.
	 */
	Want uint32

	MaxWrite            uint32
	MaxBackground       uint16
	CongestionThreshold uint16
	TimeGran            uint32
}

type FuseSession struct {
	Mountpoint string

	dev *os.File // "dev/fuse"

	inited bool // is inited or not

	bufsize int // read buffser size (/dev/fuse)

	maxGoro int // max goroutine num

	connInfo *FuseConnInfo

	FuseConfig *FuseConfig

	Opts *FuseOpt

	Debug bool

	Running bool

	readChan  chan []byte
	writeChan chan []byte

	closeCh chan interface{}

	userdata interface{} // user data
}

func NewFuseSession(mountpoint string, opts *FuseOpt, maxGoro int) *FuseSession {

	se := &FuseSession{}
	se.Init(mountpoint, opts, maxGoro)

	return se
}

func (se *FuseSession) Init(mountpoint string, opts *FuseOpt, maxGoro int) {

	se.Mountpoint = mountpoint

	se.bufsize = KERNEL_BUF_PAGES*syscall.Getpagesize() + HEADER_SIZE
	se.Opts = opts
	se.maxGoro = maxGoro

	se.connInfo = &FuseConnInfo{}

	se.connInfo.TimeGran = 1

	se.FuseConfig = &FuseConfig{}
	se.FuseConfig.Init()

	se.inited = true
}

func (se *FuseSession) IsInited() bool {
	return se.inited
}

func (se *FuseSession) SetDev(fd uintptr) {
	se.dev = os.NewFile(fd, "/dev/fuse")
}

type FuseReq struct {
	session *FuseSession

	Unique  uint64
	Uid     uint32
	Gid     uint32
	Pid     uint32
	Padding uint32

	Arg *interface{}
}

func (req *FuseReq) Init(se *FuseSession, inheader kernel.FuseInHeader) {
	req.session = se
	req.Unique = inheader.Unique
	req.Uid = inheader.UID
	req.Gid = inheader.Gid
	req.Pid = inheader.Pid
	req.Padding = inheader.Padding
}

func (req *FuseReq) GetFuseConfig() FuseConfig {
	return *req.session.FuseConfig
}

type FuseConfig struct {
	/* Fuse的开始时间的时间戳 */
	FuseStartTime int64
	/**
	 * The timeout in seconds for which file/directory attributes
	 * (as returned by e.g. the `getattr` handler) are cached.
	 */
	AttrTimeout float64
}

func (config *FuseConfig) Init() {
	config.FuseStartTime = time.Now().UnixNano()
	config.AttrTimeout = 2
}

type FuseFileInfo struct {

	/** Open flags.	 Available in open() and release() */
	Flags uint32

	/** In case of a write operation indicates if this was caused by a
	  writepage */
	Writepage uint

	/** Can be filled in by open, to use direct I/O on this file. */
	DirectIo uint

	/** Can be filled in by open, to indicate that currently
	  cached file data (that the filesystem provided the last
	  time the file was open) need not be invalidated. */
	KeepCache uint

	/** Indicates a flush operation.  Set in flush operation, also
	  maybe set in highlevel lock operation and lowlevel release
	  operation. */
	Flush uint

	/** Can be filled in by open, to indicate that the file is not
	  seekable. */
	Nonseekable uint

	/* Indicates that flock locks for this file should be
	   released.  If set, lock_owner shall contain a valid value.
	   May only be set in ->release(). */
	FlockRelease uint

	/** Padding.  Do not use*/
	Padding uint

	/** File handle.  May be filled in by filesystem in open().
	  Available in all other file operations */
	Fh uint64

	/** Lock owner id.  Available in locking operations and flush */
	LockOwner uint64

	/** Requested poll events.  Available in ->poll.  Only set on kernels
	  which support it.  If unsupported, this field is set to zero. */
	PollEvent uint32
}

func NewFuseFileInfo() FuseFileInfo {
	info := FuseFileInfo{}

	info.Writepage = 1
	info.DirectIo = 1
	info.KeepCache = 1
	info.Flush = 1
	info.Nonseekable = 1
	info.FlockRelease = 1
	info.Padding = 27

	return info
}

type FusePollhandle struct {
	Kh uint64
	Se FuseSession
}

type FuseStat struct {
	Nodeid     uint64
	Generation uint64
	Stat       syscall.Stat_t
}

// The Dirent sturct provide to outside
type FuseDirent kernel.FuseDirent

// The FuseStatfs stuct provide to outside
type FuseStatfs kernel.FuseStatfs

// The FuseIoctlOut struct provide to outside
type FuseIoctl kernel.FuseIoctlOut

// The FuseForgetOne struct provide to outside
type FuseForgetOne kernel.FuseForgetOne

// The syscall.Flock_t struct provide to outside
// This will make it easy if I want to change the struct of FuseFLock
type FuseFlock syscall.Flock_t
