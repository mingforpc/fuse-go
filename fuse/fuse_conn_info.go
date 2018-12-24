package fuse

/**************************************************************************
 * Capability bits for 'fuse_conn_info.capable' and 'fuse_conn_info.want' *
 **************************************************************************/

// FuseCapAsyncRead : Indicates that the filesystem supports asynchronous read requests.
// If this capability is not requested/available, the kernel will
// ensure that there is at most one pending read request per
// file-handle at any time, and will attempt to order read requests by
// increasing offset.
//
// This feature is enabled by default when supported by the kernel.
const FuseCapAsyncRead = (1 << 0)

// FuseCapPosixLocks : Indicates that the filesystem supports "remote" locking.
//
// This feature is enabled by default when supported by the kernel,
// and if getlk=() and setlk=() handlers are implemented.
const FuseCapPosixLocks = (1 << 1)

// FuseCapAtomicOTrunc : Indicates that the filesystem supports the O_TRUNC open flag.  If
// disabled, and an application specifies O_TRUNC, fuse first calls
// truncate=() and then open=() with O_TRUNC filtered out.
//
// This feature is enabled by default when supported by the kernel.
const FuseCapAtomicOTrunc = (1 << 3)

// FuseCapExportSupport : Indicates that the filesystem supports lookups of "." and "..".
//
// This feature is disabled by default.
const FuseCapExportSupport = (1 << 4)

// FuseCapDontMask : Indicates that the kernel should not apply the umask to the
// file mode on create operations.
//
// This feature is disabled by default.
const FuseCapDontMask = (1 << 6)

// FuseCapSliceWrite : Indicates that libfuse should try to use splice=() when writing to
// the fuse device. This may improve performance.
//
// This feature is disabled by default.
const FuseCapSliceWrite = (1 << 7)

// FuseCapSpliceMove : Indicates that libfuse should try to move pages instead of copying when
// writing to / reading from the fuse device. This may improve performance.
//
// This feature is disabled by default.
const FuseCapSpliceMove = (1 << 8)

// FuseCapSpliceRead : Indicates that libfuse should try to use splice=() when reading from
// the fuse device. This may improve performance.
//
// This feature is enabled by default when supported by the kernel and
// if the filesystem implements a write_buf=() handler.
const FuseCapSpliceRead = (1 << 9)

// FuseCapFlockLocks : If set, the calls to flock=(2) will be emulated using POSIX locks and must
// then be handled by the filesystem's setlock=() handler.
//
// If not set, flock=(2) calls will be handled by the FUSE kernel module
// internally =(so any access that does not go through the kernel cannot be taken
// into account).
//
// This feature is enabled by default when supported by the kernel and
// if the filesystem implements a flock=() handler.
const FuseCapFlockLocks = (1 << 10)

// FuseCapIoctlDir : Indicates that the filesystem supports ioctl's on directories.
//
// This feature is enabled by default when supported by the kernel.
const FuseCapIoctlDir = (1 << 11)

// FuseCapAutoInvalData : Traditionally, while a file is open the FUSE kernel module only
// asks the filesystem for an update of the file's attributes when a
// client attempts to read beyond EOF. This is unsuitable for
// e.g. network filesystems, where the file contents may change
// without the kernel knowing about it.
//
// If this flag is set, FUSE will check the validity of the attributes
// on every read. If the attributes are no longer valid =(i.e., if the
// *attr_timeout* passed to fuse_reply_attr=() or set in `struct
// fuse_entry_param` has passed), it will first issue a `getattr`
// request. If the new mtime differs from the previous value, any
// cached file *contents* will be invalidated as well.
//
// This flag should always be set when available. If all file changes
// go through the kernel, *attr_timeout* should be set to a very large
// number to avoid unnecessary getattr=() calls.
//
// This feature is enabled by default when supported by the kernel.
const FuseCapAutoInvalData = (1 << 12)

// FuseCapReaddirplus : Indicates that the filesystem supports readdirplus
//
// This feature is enabled by default when supported by the kernel and if the
// filesystem implements a readdirplus=() handler.
const FuseCapReaddirplus = (1 << 13)

// FuseCapReaddirplusAuto : Indicates that the filesystem supports adaptive readdirplus
//
// This feature is enabled by default when supported by the kernel and if the
// filesystem implements a readdirplus=() handler.
const FuseCapReaddirplusAuto = (1 << 14)

// FuseCapAsyncDIO : Indicates that the filesystem supports asynchronous direct I/O submission.
//
// If this capability is not requested/available, the kernel will ensure that
// there is at most one pending read and one pending write request per direct
// I/O file-handle at any time.
//
// This feature is enabled by default when supported by the kernel.
const FuseCapAsyncDIO = (1 << 15)

// FuseCapWritebackCache : Indicates that writeback caching should be enabled. This means that
// individual write request may be buffered and merged in the kernel
// before they are send to the filesystem.
//
// This feature is disabled by default.
const FuseCapWritebackCache = (1 << 16)

//FuseCapNoOpenSupport : Indicates support for zero-message opens. If this flag is set in
// the `capable` field of the `fuse_conn_info` structure, then the
// filesystem may return `ENOSYS` from the open=() handler to indicate
// success. Further attempts to open files will be handled in the
// kernel. =(If this flag is not set, returning ENOSYS will be treated
// as an error and signaled to the caller).
//
// Setting =(or unsetting) this flag in the `want` field has *no
// effect*.
const FuseCapNoOpenSupport = (1 << 17)

// FuseCapParallelDirops : Indicates support for parallel directory operations. If this flag
// is unset, the FUSE kernel module will ensure that lookup=() and
// readdir=() requests are never issued concurrently for the same
// directory.
//
// This feature is enabled by default when supported by the kernel.
const FuseCapParallelDirops = (1 << 18)

// FuseCapPosixACL : Indicates support for POSIX ACLs.
//
// If this feature is enabled, the kernel will cache and have
// responsibility for enforcing ACLs. ACL will be stored as xattrs and
// passed to userspace, which is responsible for updating the ACLs in
// the filesystem, keeping the file mode in sync with the ACL, and
// ensuring inheritance of default ACLs when new filesystem nodes are
// created. Note that this requires that the file system is able to
// parse and interpret the xattr representation of ACLs.
//
// Enabling this feature implicitly turns on the
// ``default_permissions`` mount option =(even if it was not passed to
// mount=(2)).
//
// This feature is disabled by default.
const FuseCapPosixACL = (1 << 19)

// FuseCapHandleKillpriv : Indicates that the filesystem is responsible for unsetting
// setuid and setgid bits when a file is written, truncated, or
// its owner is changed.
//
// This feature is enabled by default when supported by the kernel.
const FuseCapHandleKillpriv = (1 << 20)
