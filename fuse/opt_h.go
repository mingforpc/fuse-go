package fuse

// Opt : the struct to save the fuse operations
type Opt struct {
	/**
	 * Init
	 * Initialize filesystem
	 *
	 * This function is called when libfuse establishes
	 * communication with the FUSE kernel module. The file system
	 * should use this module to inspect and/or modify the
	 * connection parameters provided in the `conn` structure.
	 *
	 * Note that some parameters may be overwritten by options
	 * passed to fuse_session_new() which take precedence over the
	 * values set in this handler.
	 *
	 * There's no reply to this function
	 *
	 * conn: The fuse connection info
	 * userdata: userdata saved to session
	 *
	 * Fuse的初始化函数
	 * conn: 是与内核fuse协商的信息，如有需要直接修改
	 * userdata: 是应用需要保存的数据，不需要直接返回nil即可
	 */
	Init *func(conn *ConnInfo) (userdata interface{})

	/**
	* Destory
	* Clean up filesystem.
	*
	* Called on filesystem exit. When this method is called, the
	* connection to the kernel may be gone already, so that eg. calls
	* to fuse_lowlevel_notify_* will fail.
	*
	* There's no reply to this function
	*
	* userdata: The userdata saved in Init
	*
	* 该接口并不是FUSE应用退出时调用的程序，而是当挂载的文件系统的superblock关闭或者出错时，才触发的
	* userdata: 在Init中保存的应用数据
	 */
	Destory *func(userdata interface{})

	/**
	 * Look up a directory entry by name and get its attributes.
	 *
	 * req: request handle
	 * parentId: parent inode number of the parent directory
	 * name: the name to look up
	 * fsStat: the file stat to return
	 * res: the errno to fs
	 *
	 * 根据文件名获取文件的属性
	 */
	Lookup *func(req Req, parentId uint64, name string) (fsStat *FileStat, res int32)

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
	 * req: request handle
	 * nodeId: the inode number
	 * nlookup: the number of lookups to forget
	 */
	Forget *func(req Req, nodeid uint64, nlookup uint64)

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
	 * req: request handle
	 * nodeid: the inode number
	 * fsStat: stat the file stat
	 * res: the errno to fs
	 */
	Getattr *func(req Req, nodeid uint64) (fsStat *FileStat, res int32)

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
	 * req: request handle
	 * nodeid: the inode number
	 * attr: the attributes
	 * toSet: bit mask of attributes which should be set
	 * res: the errno to fs
	 */
	Setattr *func(req Req, nodeid uint64, attr FileStat, toSet uint32) (res int32)

	/**
	 * Read symbolic link
	 *
	 *
	 * req: request handle
	 * nodeid: the inode number
	 * path: the contents of the symbolic link
	 * res: the errno to fs. About readlink, please check [http://man7.org/linux/man-pages/man2/readlink.2.html]
	 */
	Readlink *func(req Req, nodeid uint64) (path string, res int32)

	/**
	 * Create file node
	 *
	 * Create a regular file, character device, block device, fifo or
	 * socket node.
	 *
	 *
	 * req: request handle
	 * parentid: inode number of the parent directory
	 * name: to create
	 * mode: file type and mode with which to create the new file
	 * rdev: the device number (only valid if created file is a device)
	 * fsStat: the file stat
	 * res: the errno to fs. About mknod, please check [http://man7.org/linux/man-pages/man2/mknod.2.html]
	 */
	Mknod *func(req Req, parentid uint64, name string, mode uint32, rdev uint32) (fsStat *FileStat, res int32)

	/**
	 * Create a directory
	 *
	 *
	 * req: request handle
	 * parentid: inode number of the parent directory
	 * name: to create
	 * mode: with which to create the new file
	 * fsStat: the file stat
	 * res: the errno to fs. About mkdir, please check [http://man7.org/linux/man-pages/man2/mkdir.2.html]
	 */
	Mkdir *func(req Req, parentid uint64, name string, mode uint32) (fsStat *FileStat, res int32)

	/**
	 * Remove a directory
	 *
	 * If the directory's inode's lookup count is non-zero, the
	 * file system is expected to postpone any removal of the
	 * inode until the lookup count reaches zero (see description
	 * of the forget function).
	 *
	 *
	 * req: request handle
	 * parentid: inode number of the parent directory
	 * name: to remove
	 * res: the errno to fs. About unlink, please check [http://man7.org/linux/man-pages/man2/unlink.2.html]
	 */
	Unlink *func(req Req, parentid uint64, name string) (res int32)

	/**
	 * Remove a directory
	 *
	 * If the directory's inode's lookup count is non-zero, the
	 * file system is expected to postpone any removal of the
	 * inode until the lookup count reaches zero (see description
	 * of the forget function).
	 *
	 * req: request handle
	 * parentid: inode number of the parent directory
	 * name: to remove
	 * res: the errno to fs. About rmdir, please check [http://man7.org/linux/man-pages/man2/rmdir.2.html]
	 */
	Rmdir *func(req Req, parentid uint64, name string) (res int32)

	/**
	 * Create a symbolic link
	 *
	 *
	 * req: request handle
	 * parentid: inode number of the parent directory
	 * link: the contents of the symbolic link
	 * name: to create
	 * fsStat: the file stat
	 * res: the errno to fs. About symlink, please check[http://man7.org/linux/man-pages/man2/symlink.2.html]
	 */
	Symlink *func(req Req, parentid uint64, link string, name string) (fsStat *FileStat, res int32)

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
	 * req: request handle
	 * parentid: inode number of the old parent directory
	 * name: old name
	 * newparentid: inode number of the new parent directory
	 * newname: new name
	 * res: the errno to fs. About rename, please check[http://man7.org/linux/man-pages/man2/rename.2.html]
	 */
	Rename *func(req Req, parentid uint64, name string, newparentid uint64, newname string) (res int32)

	/**
	 * Create a hard link
	 *
	 *
	 * req: request handle
	 * oldnodeid: the old inode number
	 * newparentid: inode number of the new parent directory
	 * newname: new name to create
	 * fsStat: the new file stat
	 * res: the errno to fs. About link, please check[http://man7.org/linux/man-pages/man2/link.2.html]
	 */
	Link *func(req Req, oldnodeid uint64, newparentid uint64, newname string) (fsStat *FileStat, res int32)

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
	* req: request handle
	* nodeid: the inode number
	* fi: file information
	* res: the errno to fs. About open, please check[http://man7.org/linux/man-pages/man2/open.2.html]
	*/
	Open *func(req Req, nodeid uint64, fi *FileInfo) (res int32)

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
	 * req: request handle
	 * nodeid: the inode number
	 * size: number of bytes to read
	 * offset: offset to read from
	 * fi: file information
	 * content: the content to read
	 * res: the errno to fs. About read, please check[http://man7.org/linux/man-pages/man2/read.2.html]
	 */
	Read *func(req Req, nodeid uint64, size uint32, offset uint64, fi FileInfo) (content []byte, res int32)

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
	 * req: request handle
	 * nodeid: the inode number
	 * buf: data to write
	 * offset: offset to write to
	 * fi: file information
	 * size: the size write to file
	 * res: the errno to fs. About write, please check[http://man7.org/linux/man-pages/man2/write.2.html]
	 */
	Write *func(req Req, nodeid uint64, buf []byte, offset uint64, fi FileInfo) (size uint32, res int32)

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
	 * req: request handle
	 * nodeid: the inode number
	 * fi: file information
	 * res: the errno to fs. About flush, please check[http://man7.org/linux/man-pages/man3/fflush.3.html](may be not correct)
	 */
	Flush *func(req Req, nodeid uint64, fi FileInfo) (res int32)

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
	 * req: request handle
	 * nodeid: the inode number
	 * datasync: flag indicating if only data should be flushed
	 * fi: file information
	 * res: the errno to fs. About fsync, please check[http://man7.org/linux/man-pages/man2/fsync.2.html]
	 */
	Fsync *func(req Req, nodeid uint64, datasync uint32, fi FileInfo) (res int32)

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
	 * req: request handle
	 * nodeid: the inode number
	 * fi: file information
	 * res: the errno to fs. About opendir, please check[http://man7.org/linux/man-pages/man3/opendir.3.html]
	 */
	Opendir *func(req Req, nodeid uint64, fi *FileInfo) (res int32)

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
	 * req: request handle
	 * nodeid: the inode number
	 * size: maximum number of bytes to send
	 * offset: offset to continue reading the directory stream
	 * fi: file information
	 * direntList: list of file in this directory, the binary size of direntList should not larget than size.
	 * res: the errno to fs. About readdir, please check[http://man7.org/linux/man-pages/man3/readdir.3.html]
	 *
	 */
	Readdir *func(req Req, nodeid uint64, size uint32, offset uint64, fi FileInfo) (direntList []Dirent, res int32)

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
	 * req: request handle
	 * nodeid: the inode number
	 * fi: file information
	 * res: the errno to fs.
	 */
	Releasedir *func(req Req, nodeid uint64, fi FileInfo) (res int32)

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
	 *
	 * req: request handle
	 * nodeid: the inode number
	 * fi: file information
	 * res: the errno to fs.
	 */
	Release *func(req Req, nodeid uint64, fi FileInfo) (res int32)

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
	 * req: request handle
	 * nodeid: the inode number
	 * datasync: flag indicating if only data should be flushed
	 * fi: file information.
	 * res: the errno to fs. About fsyncdir, please check[http://man7.org/linux/man-pages/man2/fsync.2.html]
	 */
	Fsyncdir *func(req Req, nodeid uint64, datasync uint32, fi FileInfo) (res int32)

	/**
	 * Get file system statistics
	 *
	 *
	 * req: request handle
	 * nodeid: the inode number, zero means "undefined"
	 * statfs: stat of file system
	 * res: the errno to fs. About fsyncdir, please check[http://man7.org/linux/man-pages/man2/fstatfs.2.html]
	 */
	Statfs *func(req Req, nodeid uint64) (statfs *Statfs, res int32)

	/**
	 * Set an extended attribute
	 *
	 * If this request is answered with an error code of ENOSYS, this is
	 * treated as a permanent failure with error code EOPNOTSUPP, i.e. all
	 * future setxattr() requests will fail with EOPNOTSUPP without being
	 * send to the filesystem process.
	 *
	 * req: request handle
	 * nodeid: the inode number, zero means "undefined"
	 * name: name of attribute
	 * value: value of attribute
	 * flags: setxattr flags
	 * res: the errno to fs. About setxattr, pease check[http://man7.org/linux/man-pages/man2/fsetxattr.2.html]
	 */
	Setxattr *func(req Req, nodeid uint64, name string, value string, flags uint32) (res int32)

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
	 * req: request handle
	 * nodeid: the inode number
	 * name: name of the extended attribute
	 * size: maximum size of the value to send
	 * value: value of the extended attribute
	 * res: the errno to fs. About getxattr, pease check[http://man7.org/linux/man-pages/man2/fgetxattr.2.html]
	 */
	Getxattr *func(req Req, nodeid uint64, name string, size uint32) (value string, res int32)

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
	 * req: request handle
	 * nodeid: the inode number
	 * size: maximum size of the list to send
	 * list: the extended attributes names string, eatch name use '\0' to split it
	 * res: the errno to fs. About getxattr, pease check[http://man7.org/linux/man-pages/man2/flistxattr.2.html]
	 */
	Listxattr *func(req Req, nodeid uint64, size uint32) (list string, res int32)

	/**
	 * Remove an extended attribute
	 *
	 * If this request is answered with an error code of ENOSYS, this is
	 * treated as a permanent failure with error code EOPNOTSUPP, i.e. all
	 * future removexattr() requests will fail with EOPNOTSUPP without being
	 * send to the filesystem process.
	 *
	 *
	 * req: request handle
	 * nodeid: the inode number
	 * name: name of the extended attribute
	 * res: the errno to fs. About removexattr, pease check[http://man7.org/linux/man-pages/man2/removexattr.2.html]
	 */
	Removexattr *func(req Req, nodeid uint64, name string) (res int32)

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
	 * req: request handle
	 * nodeid: the inode number
	 * mask: requested access mode
	 * res: the errno to fs. About access, pease check[http://man7.org/linux/man-pages/man2/access.2.html]
	 */
	Access *func(req Req, nodeid uint64, mask uint32) (res int32)

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
	 * req: request handle
	 * parentid: inode number of the parent directory
	 * name: to create
	 * mode: file type and mode with which to create the new file
	 * fi: file information, use to control open process
	 * fsStat: file stat
	 * res: the errno to fs.
	 */
	Create *func(req Req, parentid uint64, name string, mode uint32, fi *FileInfo) (fsStat *FileStat, res int32)

	/**
	 * Test for a POSIX file lock
	 *
	 *
	 * req: request handle
	 * nodeid: the inode number
	 * fi: file information
	 * lock: the region/type to test
	 * res: the errno to fs. About getlk, please check [http://man7.org/linux/man-pages/man3/lockf.3.html]
	 */
	Getlk *func(req Req, nodeid uint64, fi FileInfo, lock *Flock) (res int32)

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
	 * req: request handle
	 * nodeid: the inode number
	 * fi: file information
	 * lock: the region/type to set
	 * sleep: locking operation may sleep
	 * res: the errno to fs. About setlk, please check [http://man7.org/linux/man-pages/man3/lockf.3.html]
	 */
	Setlk *func(req Req, nodeid uint64, fi FileInfo, lock Flock, lksleep int) (res int32)

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
	 * req: request handle
	 * nodeid: the inode number
	 * blocksize: unit of block index
	 * idx: block index within file
	 * res: the errno to fs.
	 */
	Bmap *func(req Req, nodeid uint64, blocksize uint32, idx *uint64) (res int32)

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
	 * req: request handle
	 * nodeid: the inode number
	 * cmd: ioctl command
	 * arg: ioctl argument
	 * fi: file information
	 * inbuf: data fetched from the caller
	 * outbufsz: maximum size of output data
	 * ioctl: result to kernel
	 * res: the errno to fs. About ioctl, please check[http://man7.org/linux/man-pages/man2/ioctl.2.html]
	 */
	Ioctl *func(req Req, nodeid uint64, cmd uint32, arg uint64, fi FileInfo, inbuf []byte, outbufsz uint32) (ioctl *Ioctl, res int32)

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
	 * req: request handle
	 * nodeid: the inode number
	 * fi: file information
	 * ph: poll handle to be used for notification
	 * revents: returned events
	 * res: the errno to fs. About poll, please check[http://man7.org/linux/man-pages/man2/poll.2.html]
	 */
	Poll *func(req Req, nodeid uint64, fi FileInfo, ph *Pollhandle) (revents uint32, res int32)

	/**
	 * Forget about multiple inodes
	 *
	 * See description of the forget function for more
	 * information.
	 *
	 *
	 * req: request handle
	 * nodeList: the node list to forget
	 */
	ForgetMulti *func(req Req, nodeList []ForgetOne)

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
	 * req: request handle
	 * nodeid: the inode number
	 * mode: determines the operation to be performed on the given range,
	 *             see fallocate(2)
	 * offset: starting point for allocated region
	 * length: size of allocated region
	 * fi: file information
	 * res: the errno to fs. About fallocate, please check[http://man7.org/linux/man-pages/man2/fallocate.2.html]
	 */
	Fallocate *func(req Req, nodeid uint64, mode uint32, offset uint64, length uint64, fi FileInfo) (res int32)

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
	 * req: request handle
	 * nodeid: the inode number
	 * size: maximum number of bytes to send
	 * offset: offset to continue reading the directory stream
	 * fi: file information
	 * buf: result to kernel
	 * res: the errno to fs. About readdirplus , please check[http://man7.org/linux/man-pages/man3/readdir.3.html]
	 */
	Readdirplus *func(req Req, nodeid uint64, size uint32, offset uint64, fi FileInfo) (buf []byte, res int32)

	Interrupt *func(req Req)
}
