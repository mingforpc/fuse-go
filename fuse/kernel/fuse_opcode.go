package kernel

// Fuse operation code
const (
	FuseOpLookup      = 1
	FuseOpForget      = 2 /* no reply */
	FuseOpGetattr     = 3
	FuseOpSetattr     = 4
	FuseOpReadlink    = 5
	FuseOpSymlink     = 6
	FuseOpMknod       = 8
	FuseOpMkdir       = 9
	FuseOpUnlink      = 10
	FuseOpRmdir       = 11
	FuseOpRename      = 12
	FuseOpLink        = 13
	FuseOpOpen        = 14
	FuseOpRead        = 15
	FuseOpWrite       = 16
	FuseOpStatfs      = 17
	FuseOpRelease     = 18
	FuseOpFsync       = 20
	FuseOpSetxattr    = 21
	FuseOpGetxattr    = 22
	FuseOpListxattr   = 23
	FuseOpRemovexattr = 24
	FuseOpFlush       = 25
	FuseOpInit        = 26
	FuseOpOpendir     = 27
	FuseOpReaddir     = 28
	FuseOpReleasedir  = 29
	FuseOpFsyncdir    = 30
	FuseOpGetlk       = 31
	FuseOpSetlk       = 32
	FuseOpSetlkw      = 33
	FuseOpAccess      = 34
	FuseOpCreate      = 35
	FuseOpInterrupt   = 36
	FuseOpBmap        = 37
	FuseOpDestory     = 38
	FuseOpIoctl       = 39
	FuseOpPoll        = 40
	FuseOpNotifyReply = 41
	FuseOpBatckForget = 42
	FuseOpFallocate   = 43
	FuseOpReaddirplus = 44
	FuseOpRename2     = 45
	FuseOpLseek       = 46

	/* CUSE specific operations */
	CuseInit = 4096
)
