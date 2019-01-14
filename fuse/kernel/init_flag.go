package kernel

/**
 * INIT request/reply flags
 *
 * FuseAsyncRead : asynchronous read requests
 * FusePosixLocks : remote locking for POSIX file locks
 * FuseFileOps : kernel sends file handle for fstat, etc... (not yet supported)
 * FuseAtomicOTrunc : handles the O_TRUNC open flag in the filesystem
 * FuseExportSupport : filesystem handles lookups of "." and ".."
 * FuseBigWrites : filesystem can handle write size larger than 4kB
 * FuseDontMask : don't apply umask to file mode on create operations
 * FuseSpliceWrite : kernel supports splice write on the device
 * FuseSpliceMove : kernel supports splice move on the device
 * FuseSpliceRead : kernel supports splice read on the device
 * FuseFlockLocks: remote locking for BSD style file locks
 * FuseHasIoCtlDir: kernel supports ioctl on directories
 * FuseAutoInvalData: automatically invalidate cached pages
 * FuseDoReaddirplus: do READDIRPLUS (READDIR+LOOKUP in one)
 * FuseReaddirplusAuto: adaptive readdirplus
 * FuseAsyncDio: asynchronous direct I/O submission
 * FuseWritebackCache: use writeback cache for buffered writes
 * FuseNoOpenSupport: kernel supports zero-message opens
 * FuseParallelDirops: allow parallel lookups and readdir
 * FuseCapPosixACL: filesystem supports posix acls
 * FuseHandleKillPriv: fs handles killing suid/sgid/cap on write/chown/trunc
 */
const (
	FuseAsyncRead       = (1 << 0)
	FusePosixLocks      = (1 << 1)
	FuseFileOps         = (1 << 2)
	FuseAtomicOTrunc    = (1 << 3)
	FuseExportSupport   = (1 << 4)
	FuseBigWrites       = (1 << 5)
	FuseDontMask        = (1 << 6)
	FuseSpliceWrite     = (1 << 7)
	FuseSpliceMove      = (1 << 8)
	FuseSpliceRead      = (1 << 9)
	FuseFlockLocks      = (1 << 10)
	FuseHasIoCtlDir     = (1 << 11)
	FuseAutoInvalData   = (1 << 12)
	FuseDoReaddirplus   = (1 << 13)
	FuseReaddirplusAuto = (1 << 14)
	FuseAsyncDio        = (1 << 15)
	FuseWritebackCache  = (1 << 16)
	FuseNoOpenSupport   = (1 << 17)
	FuseParallelDirops  = (1 << 18)
	FuseCapPosixACL     = (1 << 19)
	FuseHandleKillPriv  = (1 << 20)
)
