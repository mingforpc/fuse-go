package fuse

import (
	"os"
	"syscall"
	"time"

	"github.com/mingforpc/fuse-go/fuse/kernel"
)

// KernelBufPages : the buffer pages of kernel
const KernelBufPages = 32

// HeaderSize : room needed in buffer to accommodate header
const HeaderSize = 0x1000

/**
 * Flags returned by the OPEN request
 *
 * FOpenDirectIO: bypass page cache for this open file
 * FOpenKeepCache: don't invalidate the data cache on open
 * FOpenNonSeekable: the file is not seekable
 */
const (
	FOpenDirectIO    = (1 << 0)
	FOpenKeepCache   = (1 << 1)
	FOpenNonSeekable = (1 << 2)
)

/* 'toSet' flags in setattr */
const (
	FuseSetAttrMode     = (1 << 0)
	FuseSetAttrUID      = (1 << 1)
	FuseSetAttrGID      = (1 << 2)
	FuseSetAttrSize     = (1 << 3)
	FuseSetAttrAtime    = (1 << 4)
	FuseSetAttrMtime    = (1 << 5)
	FuseSetAttrAtimeNow = (1 << 7)
	FuseSetAttrMtimeNow = (1 << 8)
	FuseSetAttrCtime    = (1 << 10)
)

/* XATTR set flag  */
const (
	XattrCreate  = 0x1
	XattrReplace = 0x2
)

// ConnInfo : Fuse Connection Info
type ConnInfo struct {
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

// Session : The main session to control fuse application
type Session struct {
	Mountpoint string

	dev *os.File // "dev/fuse"

	inited bool // is inited or not

	bufsize int // read buffser size (/dev/fuse)

	maxGoro int // max goroutine num

	connInfo *ConnInfo

	FuseConfig *FuseConfig

	Opts *FuseOpt

	Debug bool

	Running bool

	readChan  chan []byte
	writeChan chan []byte

	closeCh chan interface{}

	userdata interface{} // user data
}

func NewFuseSession(mountpoint string, opts *FuseOpt, maxGoro int) *Session {

	se := &Session{}
	se.Init(mountpoint, opts, maxGoro)

	return se
}

func (se *Session) Init(mountpoint string, opts *FuseOpt, maxGoro int) {

	se.Mountpoint = mountpoint

	se.bufsize = KernelBufPages*syscall.Getpagesize() + HeaderSize
	se.Opts = opts
	se.maxGoro = maxGoro

	se.connInfo = &ConnInfo{}

	se.connInfo.TimeGran = 1

	se.FuseConfig = &FuseConfig{}
	se.FuseConfig.Init()

	se.inited = true
}

func (se *Session) IsInited() bool {
	return se.inited
}

func (se *Session) SetDev(fd uintptr) {
	se.dev = os.NewFile(fd, "/dev/fuse")
}

type FuseReq struct {
	session *Session

	Unique  uint64
	Uid     uint32
	Gid     uint32
	Pid     uint32
	Padding uint32

	Arg *interface{}
}

func (req *FuseReq) Init(se *Session, inheader kernel.FuseInHeader) {
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
	Se Session
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
