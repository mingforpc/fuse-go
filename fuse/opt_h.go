package fuse

import (
	"bytes"
	"syscall"

	"github.com/mingforpc/fuse-go/fuse/kernel"
)

type FuseOpt struct {

	/**
	 * Clean up filesystem.
	 *
	 * Called on filesystem exit. When this method is called, the
	 * connection to the kernel may be gone already, so that eg. calls
	 * to fuse_lowlevel_notify_* will fail.
	 *
	 * There's no reply to this function
	 *
	 */
	Destory *func(req FuseReq, nodeid uint64) int32

	/**
	 * Look up a directory entry by name and get its attributes.
	 *
	 * @param req request handle
	 * @param parent inode number of the parent directory
	 * @param name the name to look up
	 * @param stat the file stat to return
	 */
	Lookup *func(req FuseReq, parentId uint64, name string, stat *syscall.Stat_t, generation *uint64) int32

	/**
	 * Forget about an inode
	 *
	 * This function is called when the kernel removes an inode
	 * from its internal caches.
	 *
	 * The inode's lookup count increases by one for every call to
	 * fuse_reply_entry and fuse_reply_create. The nlookup parameter
	 * indicates by how much the lookup count should be decreased.
	 *
	 * Inodes with a non-zero lookup count may receive request from
	 * the kernel even after calls to unlink, rmdir or (when
	 * overwriting an existing file) rename. Filesystems must handle
	 * such requests properly and it is recommended to defer removal
	 * of the inode until the lookup count reaches zero. Calls to
	 * unlink, rmdir or rename will be followed closely by forget
	 * unless the file or directory is open, in which case the
	 * kernel issues forget only after the release or releasedir
	 * calls.
	 *
	 * Note that if a file system will be exported over NFS the
	 * inodes lifetime must extend even beyond forget. See the
	 * generation field in struct fuse_entry_param above.
	 *
	 * On unmount the lookup count for all inodes implicitly drops
	 * to zero. It is not guaranteed that the file system will
	 * receive corresponding forget messages for the affected
	 * inodes.
	 *
	 *
	 * @param req request handle
	 * @param ino the inode number
	 * @param nlookup the number of lookups to forget
	 */
	Forget *func(req FuseReq, nodeId uint64, nlookup uint64)

	/**
	 * Get file attributes.
	 *
	 * If writeback caching is enabled, the kernel may have a
	 * better idea of a file's length than the FUSE file system
	 * (eg if there has been a write that extended the file size,
	 * but that has not yet been passed to the filesystem.n
	 *
	 * In this case, the st_size value provided by the file system
	 * will be ignored.
	 *
	 *
	 * @param req request handle
	 * @param ino the inode number
	 * @param stat the file stat
	 */
	Getattr *func(req FuseReq, nodeid uint64, stat *syscall.Stat_t) int32

	/**
	 * Set file attributes
	 *
	 * In the 'attr' argument only members indicated by the 'to_set'
	 * bitmask contain valid values.  Other members contain undefined
	 * values.
	 *
	 * Unless FUSE_CAP_HANDLE_KILLPRIV is disabled, this method is
	 * expected to reset the setuid and setgid bits if the file
	 * size or owner is being changed.
	 *
	 * If the setattr was invoked from the ftruncate() system call
	 * under Linux kernel versions 2.6.15 or later, the fi->fh will
	 * contain the value set by the open method or will be undefined
	 * if the open method didn't set any value.  Otherwise (not
	 * ftruncate call, or kernel version earlier than 2.6.15) the fi
	 * parameter will be NULL.
	 *
	 *
	 * @param req request handle
	 * @param ino the inode number
	 * @param attr the attributes
	 * @param to_set bit mask of attributes which should be set
	 */
	Setattr *func(req FuseReq, nodeid uint64, attr *syscall.Stat_t, toSet uint32) int32

	/**
	 * Read symbolic link
	 *
	 *
	 * @param req request handle
	 * @param ino the inode number
	 */
	Readlink *func(req FuseReq, nodeid uint64, path *string) int32

	/**
	 * Create file node
	 *
	 * Create a regular file, character device, block device, fifo or
	 * socket node.
	 *
	 *
	 * @param req request handle
	 * @param parent inode number of the parent directory
	 * @param name to create
	 * @param mode file type and mode with which to create the new file
	 * @param rdev the device number (only valid if created file is a device)
	 * @param nodeid the file id
	 * @param generation the file generation id
	 * @param stat the file stat
	 */
	Mknod *func(req FuseReq, parentid uint64, name string, mode uint32, rdev uint32) (*FuseStat, int32)

	/**
	 * Create a directory
	 *
	 *
	 * @param req request handle
	 * @param parent inode number of the parent directory
	 * @param name to create
	 * @param mode with which to create the new file
	 */
	Mkdir *func(req FuseReq, parentid uint64, name string, mode uint32) (*FuseStat, int32)

	/**
	 * Remove a directory
	 *
	 * If the directory's inode's lookup count is non-zero, the
	 * file system is expected to postpone any removal of the
	 * inode until the lookup count reaches zero (see description
	 * of the forget function).
	 *
	 *
	 * @param req request handle
	 * @param parent inode number of the parent directory
	 * @param name to remove
	 */
	Unlink *func(req FuseReq, parentid uint64, name string) int32

	/**
	 * Remove a directory
	 *
	 * If the directory's inode's lookup count is non-zero, the
	 * file system is expected to postpone any removal of the
	 * inode until the lookup count reaches zero (see description
	 * of the forget function).
	 *
	 * @param req request handle
	 * @param parent inode number of the parent directory
	 * @param name to remove
	 */
	Rmdir *func(req FuseReq, parentid uint64, name string) int32

	/**
	 * Create a symbolic link
	 *
	 *
	 * @param req request handle
	 * @param link the contents of the symbolic link
	 * @param parent inode number of the parent directory
	 * @param name to create
	 */
	Symlink *func(req FuseReq, parentid uint64, link string, name string) (*FuseStat, int32)

	/** Rename a file
	 *
	 * If the target exists it should be atomically replaced. If
	 * the target's inode's lookup count is non-zero, the file
	 * system is expected to postpone any removal of the inode
	 * until the lookup count reaches zero (see description of the
	 * forget function).
	 *
	 * If this request is answered with an error code of ENOSYS, this is
	 * treated as a permanent failure with error code EINVAL, i.e. all
	 * future bmap requests will fail with EINVAL without being
	 * send to the filesystem process.
	 *
	 * *flags* may be `RENAME_EXCHANGE` or `RENAME_NOREPLACE`. If
	 * RENAME_NOREPLACE is specified, the filesystem must not
	 * overwrite *newname* if it exists and return an error
	 * instead. If `RENAME_EXCHANGE` is specified, the filesystem
	 * must atomically exchange the two files, i.e. both must
	 * exist and neither may be deleted.
	 *
	 *
	 * @param req request handle
	 * @param parent inode number of the old parent directory
	 * @param name old name
	 * @param newparent inode number of the new parent directory
	 * @param newname new name
	 */
	Rename *func(req FuseReq, parentid uint64, name string, newparentid uint64, newname string) int32

	/**
	 * Create a hard link
	 *
	 *
	 * @param req request handle
	 * @param oldnodeid the old inode number
	 * @param newparent inode number of the new parent directory
	 * @param newname new name to create
	 */
	Link *func(req FuseReq, oldnodeid uint64, newparentid uint64, newname string) (*FuseStat, int32)

	/**
	* Open a file
	*
	* Open flags are available in fi->flags. The following rules
	* apply.
	*
	*  - Creation (O_CREAT, O_EXCL, O_NOCTTY) flags will be
	*    filtered out / handled by the kernel.
	*
	*  - Access modes (O_RDONLY, O_WRONLY, O_RDWR) should be used
	*    by the filesystem to check if the operation is
	*    permitted.  If the ``-o default_permissions`` mount
	*    option is given, this check is already done by the
	*    kernel before calling open() and may thus be omitted by
	*    the filesystem.
	*
	*  - When writeback caching is enabled, the kernel may send
	*    read requests even for files opened with O_WRONLY. The
	*    filesystem should be prepared to handle this.
	*
	*  - When writeback caching is disabled, the filesystem is
	*    expected to properly handle the O_APPEND flag and ensure
	*    that each write is appending to the end of the file.
	*
		*  - When writeback caching is enabled, the kernel will
	*    handle O_APPEND. However, unless all changes to the file
	*    come through the kernel this will not work reliably. The
	*    filesystem should thus either ignore the O_APPEND flag
	*    (and let the kernel handle it), or return an error
	*    (indicating that reliably O_APPEND is not available).
	*
	* Filesystem may store an arbitrary file handle (pointer,
	* index, etc) in fi->fh, and use this in other all other file
	* operations (read, write, flush, release, fsync).
	*
	* Filesystem may also implement stateless file I/O and not store
	* anything in fi->fh.
	*
	* There are also some flags (direct_io, keep_cache) which the
	* filesystem may set in fi, to change the way the file is opened.
	* See fuse_file_info structure in <fuse_common.h> for more details.
	*
	* If this request is answered with an error code of ENOSYS
	* and FUSE_CAP_NO_OPEN_SUPPORT is set in
	* `fuse_conn_info.capable`, this is treated as success and
	* future calls to open will also succeed without being send
	* to the filesystem process.
	*
	*
	* @param req request handle
	* @param ino the inode number
	* @param fi file information
	*/
	Open *func(req FuseReq, nodeid uint64, fi *FuseFileInfo) int32

	/**
	 * Read data
	 *
	 * Read should send exactly the number of bytes requested except
	 * on EOF or error, otherwise the rest of the data will be
	 * substituted with zeroes.  An exception to this is when the file
	 * has been opened in 'direct_io' mode, in which case the return
	 * value of the read system call will reflect the return value of
	 * this operation.
	 *
	 * fi->fh will contain the value set by the open method, or will
	 * be undefined if the open method didn't set any value.
	 *
	 *
	 * @param req request handle
	 * @param ino the inode number
	 * @param size number of bytes to read
	 * @param off offset to read from
	 * @param fi file information
	 */
	Read *func(req FuseReq, nodeid uint64, size uint32, offset uint64, fi FuseFileInfo) ([]byte, int32)

	/**
	 * Write data
	 *
	 * Write should return exactly the number of bytes requested
	 * except on error.  An exception to this is when the file has
	 * been opened in 'direct_io' mode, in which case the return value
	 * of the write system call will reflect the return value of this
	 * operation.
	 *
	 * Unless FUSE_CAP_HANDLE_KILLPRIV is disabled, this method is
	 * expected to reset the setuid and setgid bits.
	 *
	 * fi->fh will contain the value set by the open method, or will
	 * be undefined if the open method didn't set any value.
	 *
	 *
	 * @param req request handle
	 * @param ino the inode number
	 * @param buf data to write
	 * @param size number of bytes to write
	 * @param off offset to write to
	 * @param fi file information
	 */
	Write *func(req FuseReq, nodeid uint64, buf []byte, offset uint64, fi FuseFileInfo) (size uint32, res int32)

	/**
	 * Flush method
	 *
	 * This is called on each close() of the opened file.
	 *
	 * Since file descriptors can be duplicated (dup, dup2, fork), for
	 * one open call there may be many flush calls.
	 *
	 * Filesystems shouldn't assume that flush will always be called
	 * after some writes, or that if will be called at all.
	 *
	 * fi->fh will contain the value set by the open method, or will
	 * be undefined if the open method didn't set any value.
	 *
	 * NOTE: the name of the method is misleading, since (unlike
	 * fsync) the filesystem is not forced to flush pending writes.
	 * One reason to flush data, is if the filesystem wants to return
	 * write errors.
	 *
	 * If the filesystem supports file locking operations (setlk,
	 * getlk) it should remove all locks belonging to 'fi->owner'.
	 *
	 * If this request is answered with an error code of ENOSYS,
	 * this is treated as success and future calls to flush() will
	 * succeed automatically without being send to the filesystem
	 * process.
	 *
	 *
	 * @param req request handle
	 * @param ino the inode number
	 * @param fi file information
	 */
	Flush *func(req FuseReq, nodeid uint64, fi FuseFileInfo) int32

	/**
	 * Synchronize file contents
	 *
	 * If the datasync parameter is non-zero, then only the user data
	 * should be flushed, not the meta data.
	 *
	 * If this request is answered with an error code of ENOSYS,
	 * this is treated as success and future calls to fsync() will
	 * succeed automatically without being send to the filesystem
	 * process.
	 *
	 *
	 * @param req request handle
	 * @param ino the inode number
	 * @param datasync flag indicating if only data should be flushed
	 * @param fi file information
	 */
	Fsync *func(req FuseReq, nodeid uint64, datasync uint32, fi FuseFileInfo) int32

	/**
	 * Open a directory
	 *
	 * Filesystem may store an arbitrary file handle (pointer, index,
	 * etc) in fi->fh, and use this in other all other directory
	 * stream operations (readdir, releasedir, fsyncdir).
	 *
	 * Filesystem may also implement stateless directory I/O and not
	 * store anything in fi->fh, though that makes it impossible to
	 * implement standard conforming directory stream operations in
	 * case the contents of the directory can change between opendir
	 * and releasedir.
	 *
	 *
	 * @param req request handle
	 * @param ino the inode number
	 * @param fi file information
	 */
	Opendir *func(req FuseReq, nodeid uint64, fi *FuseFileInfo) int32

	/**
	 * Read directory
	 *
	 * Send a buffer filled using fuse_add_direntry(), with size not
	 * exceeding the requested size.  Send an empty buffer on end of
	 * stream.
	 *
	 * fi->fh will contain the value set by the opendir method, or
	 * will be undefined if the opendir method didn't set any value.
	 *
	 * Returning a directory entry from readdir() does not affect
	 * its lookup count.
	 *
	 * The function does not have to report the '.' and '..'
	 * entries, but is allowed to do so.
	 *
	 *
	 * @param req request handle
	 * @param ino the inode number
	 * @param size maximum number of bytes to send
	 * @param off offset to continue reading the directory stream
	 * @param fi file information
	 */
	Readdir *func(req FuseReq, nodeid uint64, size uint32, offset uint64, fi FuseFileInfo) ([]kernel.FuseDirent, int32)

	/**
	 * Release an open directory
	 *
	 * For every opendir call there will be exactly one releasedir
	 * call.
	 *
	 * fi->fh will contain the value set by the opendir method, or
	 * will be undefined if the opendir method didn't set any value.
	 *
	 *
	 * @param req request handle
	 * @param ino the inode number
	 * @param fi file information
	 */
	Releasedir *func(req FuseReq, nodeid uint64, fi FuseFileInfo) int32

	/**
	 * Release an open file
	 *
	 * Release is called when there are no more references to an open
	 * file: all file descriptors are closed and all memory mappings
	 * are unmapped.
	 *
	 * For every open call there will be exactly one release call.
	 *
	 * The filesystem may reply with an error, but error values are
	 * not returned to close() or munmap() which triggered the
	 * release.
	 *
	 * fi->fh will contain the value set by the open method, or will
	 * be undefined if the open method didn't set any value.
	 * fi->flags will contain the same flags as for open.
	 *
	 * Valid replies:
	 *   fuse_reply_err
	 *
	 * @param req request handle
	 * @param ino the inode number
	 * @param fi file information
	 */
	Release *func(req FuseReq, nodeid uint64, fi FuseFileInfo) int32

	/**
	 * Synchronize directory contents
	 *
	 * If the datasync parameter is non-zero, then only the directory
	 * contents should be flushed, not the meta data.
	 *
	 * fi->fh will contain the value set by the opendir method, or
	 * will be undefined if the opendir method didn't set any value.
	 *
	 * If this request is answered with an error code of ENOSYS,
	 * this is treated as success and future calls to fsyncdir() will
	 * succeed automatically without being send to the filesystem
	 * process.
	 *
	 *
	 * @param req request handle
	 * @param ino the inode number
	 * @param datasync flag indicating if only data should be flushed
	 * @param fi file information
	 */
	Fsyncdir *func(req FuseReq, nodeid uint64, datasync uint32, fi FuseFileInfo) int32

	/**
	 * Get file system statistics
	 *
	 *
	 * @param req request handle
	 * @param ino the inode number, zero means "undefined"
	 */
	Statfs *func(req FuseReq, nodeid uint64, statfs *kernel.FuseStatfs) int32

	/**
	 * Set an extended attribute
	 *
	 * If this request is answered with an error code of ENOSYS, this is
	 * treated as a permanent failure with error code EOPNOTSUPP, i.e. all
	 * future setxattr() requests will fail with EOPNOTSUPP without being
	 * send to the filesystem process.
	 *
	 */
	Setxattr *func(req FuseReq, nodeid uint64, name string, value string, flags uint32) int32

	/**
	 * Get an extended attribute
	 *
	 * If size is zero, the size of the value should be sent with
	 * fuse_reply_xattr.
	 *
	 * If the size is non-zero, and the value fits in the buffer, the
	 * value should be sent with fuse_reply_buf.
	 *
	 * If the size is too small for the value, the ERANGE error should
	 * be sent.
	 *
	 * If this request is answered with an error code of ENOSYS, this is
	 * treated as a permanent failure with error code EOPNOTSUPP, i.e. all
	 * future getxattr() requests will fail with EOPNOTSUPP without being
	 * send to the filesystem process.
	 *
	 *
	 * @param req request handle
	 * @param ino the inode number
	 * @param name of the extended attribute
	 * @param size maximum size of the value to send
	 */
	Getxattr *func(req FuseReq, nodeid uint64, name string, size uint32) (value string, err int32)

	/**
	 * List extended attribute names
	 *
	 * If size is zero, the total size of the attribute list should be
	 * sent with fuse_reply_xattr.
	 *
	 * If the size is non-zero, and the null character separated
	 * attribute list fits in the buffer, the list should be sent with
	 * fuse_reply_buf.
	 *
	 * If the size is too small for the list, the ERANGE error should
	 * be sent.
	 *
	 * If this request is answered with an error code of ENOSYS, this is
	 * treated as a permanent failure with error code EOPNOTSUPP, i.e. all
	 * future listxattr() requests will fail with EOPNOTSUPP without being
	 * send to the filesystem process.
	 *
	 *
	 * @param req request handle
	 * @param ino the inode number
	 * @param size maximum size of the list to send
	 */
	Listxattr *func(req FuseReq, nodeid uint64, size uint32) (list string, err int32)

	/**
	 * Remove an extended attribute
	 *
	 * If this request is answered with an error code of ENOSYS, this is
	 * treated as a permanent failure with error code EOPNOTSUPP, i.e. all
	 * future removexattr() requests will fail with EOPNOTSUPP without being
	 * send to the filesystem process.
	 *
	 *
	 * @param req request handle
	 * @param ino the inode number
	 * @param name of the extended attribute
	 */
	Removexattr *func(req FuseReq, nodeid uint64, name string) int32

	/**
	 * Check file access permissions
	 *
	 * This will be called for the access() and chdir() system
	 * calls.  If the 'default_permissions' mount option is given,
	 * this method is not called.
	 *
	 * This method is not called under Linux kernel versions 2.4.x
	 *
	 * If this request is answered with an error code of ENOSYS, this is
	 * treated as a permanent success, i.e. this and all future access()
	 * requests will succeed without being send to the filesystem process.
	 *
	 *
	 * @param req request handle
	 * @param ino the inode number
	 * @param mask requested access mode
	 */
	Access *func(req FuseReq, nodeid uint64, mask uint32) int32

	/**
	 * Create and open a file
	 *
	 * If the file does not exist, first create it with the specified
	 * mode, and then open it.
	 *
	 * See the description of the open handler for more
	 * information.
	 *
	 * If this method is not implemented or under Linux kernel
	 * versions earlier than 2.6.15, the mknod() and open() methods
	 * will be called instead.
	 *
	 * If this request is answered with an error code of ENOSYS, the handler
	 * is treated as not implemented (i.e., for this and future requests the
	 * mknod() and open() handlers will be called instead).
	 *
	 *
	 * @param req request handle
	 * @param parent inode number of the parent directory
	 * @param name to create
	 * @param mode file type and mode with which to create the new file
	 * @param fi file information, 用来设置Open的操作
	 */
	Create *func(req FuseReq, parentid uint64, name string, mode uint32, fi *FuseFileInfo) (*FuseStat, int32)

	/**
	 * Test for a POSIX file lock
	 *
	 *
	 * @param req request handle
	 * @param ino the inode number
	 * @param fi file information
	 * @param lock the region/type to test
	 */
	Getlk *func(req FuseReq, nodeid uint64, fi FuseFileInfo, lock *syscall.Flock_t) int32

	/**
	 * Acquire, modify or release a POSIX file lock
	 *
	 * For POSIX threads (NPTL) there's a 1-1 relation between pid and
	 * owner, but otherwise this is not always the case.  For checking
	 * lock ownership, 'fi->owner' must be used.  The l_pid field in
	 * 'struct flock' should only be used to fill in this field in
	 * getlk().
	 *
	 * Note: if the locking methods are not implemented, the kernel
	 * will still allow file locking to work locally.  Hence these are
	 * only interesting for network filesystems and similar.
	 *
	 *
	 * @param req request handle
	 * @param ino the inode number
	 * @param fi file information
	 * @param lock the region/type to set
	 * @param sleep locking operation may sleep
	 */
	Setlk *func(req FuseReq, nodeid uint64, fi FuseFileInfo, lock syscall.Flock_t, lksleep int) int32

	/**
	 * Map block index within file to block index within device
	 *
	 * Note: This makes sense only for block device backed filesystems
	 * mounted with the 'blkdev' option
	 *
	 * If this request is answered with an error code of ENOSYS, this is
	 * treated as a permanent failure, i.e. all future bmap() requests will
	 * fail with the same error code without being send to the filesystem
	 * process.
	 *
	 *
	 * @param req request handle
	 * @param ino the inode number
	 * @param blocksize unit of block index
	 * @param idx block index within file
	 */
	Bmap *func(req FuseReq, nodeid uint64, blocksize uint32, idx *uint64) int32

	/**
	 * Ioctl
	 *
	 * Note: For unrestricted ioctls (not allowed for FUSE
	 * servers), data in and out areas can be discovered by giving
	 * iovs and setting FUSE_IOCTL_RETRY in *flags*.  For
	 * restricted ioctls, kernel prepares in/out data area
	 * according to the information encoded in cmd.
	 *
	 *
	 * @param req request handle
	 * @param ino the inode number
	 * @param cmd ioctl command
	 * @param arg ioctl argument
	 * @param fi file information
	 * @param flags for FUSE_IOCTL_* flags
	 * @param in_buf data fetched from the caller
	 * @param in_bufsz number of fetched bytes
	 * @param out_bufsz maximum size of output data
	 */
	Ioctl *func(req FuseReq, nodeid uint64, cmd uint32, arg uint64, fi FuseFileInfo, inbuf []byte, outbufsz uint32, ioctlOut *kernel.FuseIoctlOut) int32

	/**
	 * Poll for IO readiness
	 *
	 * Note: If ph is non-NULL, the client should notify
	 * when IO readiness events occur by calling
	 * fuse_lowlevel_notify_poll() with the specified ph.
	 *
	 * Regardless of the number of times poll with a non-NULL ph
	 * is received, single notification is enough to clear all.
	 * Notifying more times incurs overhead but doesn't harm
	 * correctness.
	 *
	 * The callee is responsible for destroying ph with
	 * fuse_pollhandle_destroy() when no longer in use.
	 *
	 * If this request is answered with an error code of ENOSYS, this is
	 * treated as success (with a kernel-defined default poll-mask) and
	 * future calls to pull() will succeed the same way without being send
	 * to the filesystem process.
	 *
	 *
	 * @param req request handle
	 * @param ino the inode number
	 * @param fi file information
	 * @param ph poll handle to be used for notification
	 */
	Poll *func(req FuseReq, nodeid uint64, fi FuseFileInfo, ph *FusePollhandle, revents *uint32) int32

	/**
	 * Forget about multiple inodes
	 *
	 * See description of the forget function for more
	 * information.
	 *
	 *
	 * @param req request handle
	 */
	ForgetMulti *func(req FuseReq, count uint32, nodeList []kernel.FuseForgetOne)

	/**
	 * Allocate requested space. If this function returns success then
	 * subsequent writes to the specified range shall not fail due to the lack
	 * of free space on the file system storage media.
	 *
	 * If this request is answered with an error code of ENOSYS, this is
	 * treated as a permanent failure with error code EOPNOTSUPP, i.e. all
	 * future fallocate() requests will fail with EOPNOTSUPP without being
	 * send to the filesystem process.
	 *
	 *
	 * @param req request handle
	 * @param ino the inode number
	 * @param offset starting point for allocated region
	 * @param length size of allocated region
	 * @param mode determines the operation to be performed on the given range,
	 *             see fallocate(2)
	 */
	Fallocate *func(req FuseReq, nodeid uint64, mode uint32, offset uint64, length uint64, fi FuseFileInfo) int32

	/**
	 * Read directory with attributes
	 *
	 * Send a buffer filled using fuse_add_direntry_plus(), with size not
	 * exceeding the requested size.  Send an empty buffer on end of
	 * stream.
	 *
	 * fi->fh will contain the value set by the opendir method, or
	 * will be undefined if the opendir method didn't set any value.
	 *
	 * In contrast to readdir() (which does not affect the lookup counts),
	 * the lookup count of every entry returned by readdirplus(), except "."
	 * and "..", is incremented by one.
	 *
	 *
	 * @param req request handle
	 * @param ino the inode number
	 * @param size maximum number of bytes to send
	 * @param off offset to continue reading the directory stream
	 * @param fi file information
	 */
	Readdirplus *func(req FuseReq, nodeid uint64, size uint32, offset uint64, fi FuseFileInfo, buf *bytes.Buffer) int32

	Interrupt *func(req FuseReq)
}
