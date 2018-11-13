package kernel

const (
	FUSE_LOOKUP       = 1
	FUSE_FORGET       = 2 /* no reply */
	FUSE_GETATTR      = 3
	FUSE_SETATTR      = 4
	FUSE_READLINK     = 5
	FUSE_SYMLINK      = 6
	FUSE_MKNOD        = 8
	FUSE_MKDIR        = 9
	FUSE_UNLINK       = 10
	FUSE_RMDIR        = 11
	FUSE_RENAME       = 12
	FUSE_LINK         = 13
	FUSE_OPEN         = 14
	FUSE_READ         = 15
	FUSE_WRITE        = 16
	FUSE_STATFS       = 17
	FUSE_RELEASE      = 18
	FUSE_FSYNC        = 20
	FUSE_SETXATTR     = 21
	FUSE_GETXATTR     = 22
	FUSE_LISTXATTR    = 23
	FUSE_REMOVEXATTR  = 24
	FUSE_FLUSH        = 25
	FUSE_INIT         = 26
	FUSE_OPENDIR      = 27
	FUSE_READDIR      = 28
	FUSE_RELEASEDIR   = 29
	FUSE_FSYNCDIR     = 30
	FUSE_GETLK        = 31
	FUSE_SETLK        = 32
	FUSE_SETLKW       = 33
	FUSE_ACCESS       = 34
	FUSE_CREATE       = 35
	FUSE_INTERRUPT    = 36
	FUSE_BMAP         = 37
	FUSE_DESTROY      = 38
	FUSE_IOCTL        = 39
	FUSE_POLL         = 40
	FUSE_NOTIFY_REPLY = 41
	FUSE_BATCH_FORGET = 42
	FUSE_FALLOCATE    = 43
	FUSE_READDIRPLUS  = 44
	FUSE_RENAME2      = 45
	FUSE_LSEEK        = 46

	/* CUSE specific operations */
	CUSE_INIT = 4096
)
