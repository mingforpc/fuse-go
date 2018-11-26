package fuse

import (
	"bytes"
	"os"
	"syscall"

	"github.com/mingforpc/fuse-go/fuse/common"
	"github.com/mingforpc/fuse-go/fuse/errno"
	"github.com/mingforpc/fuse-go/fuse/kernel"
	"github.com/mingforpc/fuse-go/fuse/log"
)

func doInit(req FuseReq, nodeid uint64, initOut *kernel.FuseInitOut) int32 {

	se := req.session
	initIn := (*req.Arg).(kernel.FuseInitIn)
	if se.Debug {
		log.Trace.Printf("INIT: %+v \n", initIn)
	}

	bufsize := se.Bufsize

	se.ConnInfo.Major = initIn.Major
	se.ConnInfo.Minor = initIn.Minor
	se.ConnInfo.MaxReadahead = initIn.MaxReadahead

	if bufsize < kernel.FUSE_MIN_READ_BUFFER {
		log.Warning.Printf("fuse: warning: buffer size too small: %d\n", bufsize)
		bufsize = kernel.FUSE_MIN_READ_BUFFER
	}

	bufsize -= 4096
	if uint32(bufsize) < se.ConnInfo.MaxWrite {
		se.ConnInfo.MaxWrite = uint32(bufsize)
	}

	initOut.Major = kernel.FUSE_KERNEL_VERSION
	initOut.Minor = kernel.FUSE_KERNEL_MINOR_VERSION
	initOut.MaxReadahead = se.ConnInfo.MaxReadahead
	initOut.MaxWrite = se.ConnInfo.MaxWrite
	initOut.TimeGran = se.ConnInfo.TimeGran
	initOut.CongestionThreshold = se.ConnInfo.CongestionThreshold

	// To remember what fuse kernel can do
	if initIn.Flags&kernel.FUSE_ASYNC_READ > 0 {
		se.ConnInfo.Capable |= kernel.FUSE_ASYNC_READ
	}
	if initIn.Flags&kernel.FUSE_POSIX_LOCKS > 0 {
		se.ConnInfo.Capable |= kernel.FUSE_POSIX_LOCKS
	}
	if initIn.Flags&kernel.FUSE_FILE_OPS > 0 {
		se.ConnInfo.Capable |= kernel.FUSE_FILE_OPS
	}
	if initIn.Flags&kernel.FUSE_ATOMIC_O_TRUNC > 0 {
		se.ConnInfo.Capable |= kernel.FUSE_ATOMIC_O_TRUNC
	}
	if initIn.Flags&kernel.FUSE_EXPORT_SUPPORT > 0 {
		se.ConnInfo.Capable |= kernel.FUSE_EXPORT_SUPPORT
	}
	if initIn.Flags&kernel.FUSE_BIG_WRITES > 0 {
		se.ConnInfo.Capable |= kernel.FUSE_BIG_WRITES
	}
	if initIn.Flags&kernel.FUSE_DONT_MASK > 0 {
		se.ConnInfo.Capable |= kernel.FUSE_DONT_MASK
	}
	if initIn.Flags&kernel.FUSE_SPLICE_WRITE > 0 {
		se.ConnInfo.Capable |= kernel.FUSE_SPLICE_WRITE
	}
	if initIn.Flags&kernel.FUSE_SPLICE_MOVE > 0 {
		se.ConnInfo.Capable |= kernel.FUSE_SPLICE_MOVE
	}
	if initIn.Flags&kernel.FUSE_SPLICE_READ > 0 {
		se.ConnInfo.Capable |= kernel.FUSE_SPLICE_READ
	}
	if initIn.Flags&kernel.FUSE_FLOCK_LOCKS > 0 {
		se.ConnInfo.Capable |= kernel.FUSE_FLOCK_LOCKS
	}
	if initIn.Flags&kernel.FUSE_HAS_IOCTL_DIR > 0 {
		se.ConnInfo.Capable |= kernel.FUSE_HAS_IOCTL_DIR
	}
	if initIn.Flags&kernel.FUSE_AUTO_INVAL_DATA > 0 {
		se.ConnInfo.Capable |= kernel.FUSE_AUTO_INVAL_DATA
	}
	if initIn.Flags&kernel.FUSE_DO_READDIRPLUS > 0 {
		se.ConnInfo.Capable |= kernel.FUSE_DO_READDIRPLUS
	}
	if initIn.Flags&kernel.FUSE_READDIRPLUS_AUTO > 0 {
		se.ConnInfo.Capable |= kernel.FUSE_READDIRPLUS_AUTO
	}
	if initIn.Flags&kernel.FUSE_ASYNC_DIO > 0 {
		se.ConnInfo.Capable |= kernel.FUSE_ASYNC_DIO
	}
	if initIn.Flags&kernel.FUSE_WRITEBACK_CACHE > 0 {
		se.ConnInfo.Capable |= kernel.FUSE_WRITEBACK_CACHE
	}
	if initIn.Flags&kernel.FUSE_NO_OPEN_SUPPORT > 0 {
		se.ConnInfo.Capable |= kernel.FUSE_NO_OPEN_SUPPORT
	}
	if initIn.Flags&kernel.FUSE_PARALLEL_DIROPS > 0 {
		se.ConnInfo.Capable |= kernel.FUSE_PARALLEL_DIROPS
	}
	if initIn.Flags&kernel.FUSE_HANDLE_KILLPRIV > 0 {
		se.ConnInfo.Capable |= kernel.FUSE_HANDLE_KILLPRIV
	}
	if initIn.Flags&kernel.FUSE_POSIX_ACL > 0 {
		se.ConnInfo.Capable |= kernel.FUSE_POSIX_ACL
	}

	// Default settings for modern filesystems.
	// TODO: support write_buf, flock
	if (se.ConnInfo.Capable & FUSE_CAP_ASYNC_READ) > 0 {
		se.ConnInfo.Want |= FUSE_CAP_ASYNC_READ
	}
	if (se.ConnInfo.Capable & FUSE_CAP_PARALLEL_DIROPS) > 0 {
		se.ConnInfo.Want |= FUSE_CAP_PARALLEL_DIROPS
	}
	if (se.ConnInfo.Capable & FUSE_CAP_AUTO_INVAL_DATA) > 0 {
		se.ConnInfo.Want |= FUSE_CAP_AUTO_INVAL_DATA
	}
	if (se.ConnInfo.Capable & FUSE_CAP_HANDLE_KILLPRIV) > 0 {
		se.ConnInfo.Want |= FUSE_CAP_HANDLE_KILLPRIV
	}
	if (se.ConnInfo.Capable & FUSE_CAP_ASYNC_DIO) > 0 {
		se.ConnInfo.Want |= FUSE_CAP_ASYNC_DIO
	}
	if (se.ConnInfo.Capable & FUSE_CAP_IOCTL_DIR) > 0 {
		se.ConnInfo.Want |= FUSE_CAP_IOCTL_DIR
	}
	if (se.ConnInfo.Capable & FUSE_CAP_ATOMIC_O_TRUNC) > 0 {
		se.ConnInfo.Want |= FUSE_CAP_ATOMIC_O_TRUNC
	}
	if se.Opts.Getlk != nil && se.Opts.Setlk != nil {
		se.ConnInfo.Want |= FUSE_CAP_POSIX_LOCKS
	}
	if se.Opts.Readdirplus != nil {
		se.ConnInfo.Want |= FUSE_CAP_READDIRPLUS
		se.ConnInfo.Want |= FUSE_CAP_READDIRPLUS_AUTO
	}

	// To set what we want fuse kenel to do
	if se.ConnInfo.Want&kernel.FUSE_ASYNC_READ > 0 {
		initOut.Flags |= kernel.FUSE_ASYNC_READ
	}
	if se.ConnInfo.Want&kernel.FUSE_POSIX_LOCKS > 0 {
		initOut.Flags |= kernel.FUSE_POSIX_LOCKS
	}
	if se.ConnInfo.Want&kernel.FUSE_FILE_OPS > 0 {
		initOut.Flags |= kernel.FUSE_FILE_OPS
	}
	if se.ConnInfo.Want&kernel.FUSE_ATOMIC_O_TRUNC > 0 {
		initOut.Flags |= kernel.FUSE_ATOMIC_O_TRUNC
	}
	if se.ConnInfo.Want&kernel.FUSE_EXPORT_SUPPORT > 0 {
		initOut.Flags |= kernel.FUSE_EXPORT_SUPPORT
	}
	if se.ConnInfo.Want&kernel.FUSE_BIG_WRITES > 0 {
		initOut.Flags |= kernel.FUSE_BIG_WRITES
	}
	if se.ConnInfo.Want&kernel.FUSE_DONT_MASK > 0 {
		initOut.Flags |= kernel.FUSE_DONT_MASK
	}
	if se.ConnInfo.Want&kernel.FUSE_SPLICE_WRITE > 0 {
		initOut.Flags |= kernel.FUSE_SPLICE_WRITE
	}
	if se.ConnInfo.Want&kernel.FUSE_SPLICE_MOVE > 0 {
		initOut.Flags |= kernel.FUSE_SPLICE_MOVE
	}
	if se.ConnInfo.Want&kernel.FUSE_SPLICE_READ > 0 {
		initOut.Flags |= kernel.FUSE_SPLICE_READ
	}
	if se.ConnInfo.Want&kernel.FUSE_FLOCK_LOCKS > 0 {
		initOut.Flags |= kernel.FUSE_FLOCK_LOCKS
	}
	if se.ConnInfo.Want&kernel.FUSE_HAS_IOCTL_DIR > 0 {
		initOut.Flags |= kernel.FUSE_HAS_IOCTL_DIR
	}
	if se.ConnInfo.Want&kernel.FUSE_AUTO_INVAL_DATA > 0 {
		initOut.Flags |= kernel.FUSE_AUTO_INVAL_DATA
	}
	if se.ConnInfo.Want&kernel.FUSE_DO_READDIRPLUS > 0 {
		initOut.Flags |= kernel.FUSE_DO_READDIRPLUS
	}
	if se.ConnInfo.Want&kernel.FUSE_READDIRPLUS_AUTO > 0 {
		initOut.Flags |= kernel.FUSE_READDIRPLUS_AUTO
	}
	if se.ConnInfo.Want&kernel.FUSE_ASYNC_DIO > 0 {
		initOut.Flags |= kernel.FUSE_ASYNC_DIO
	}
	if se.ConnInfo.Want&kernel.FUSE_WRITEBACK_CACHE > 0 {
		initOut.Flags |= kernel.FUSE_WRITEBACK_CACHE
	}
	if se.ConnInfo.Want&kernel.FUSE_NO_OPEN_SUPPORT > 0 {
		initOut.Flags |= kernel.FUSE_NO_OPEN_SUPPORT
	}
	if se.ConnInfo.Want&kernel.FUSE_PARALLEL_DIROPS > 0 {
		initOut.Flags |= kernel.FUSE_PARALLEL_DIROPS
	}
	if se.ConnInfo.Want&kernel.FUSE_HANDLE_KILLPRIV > 0 {
		initOut.Flags |= kernel.FUSE_HANDLE_KILLPRIV
	}
	if se.ConnInfo.Want&kernel.FUSE_POSIX_ACL > 0 {
		initOut.Flags |= kernel.FUSE_POSIX_ACL
	}

	return errno.SUCCESS
}

func doDestory(req FuseReq, nodeid uint64) int32 {
	se := req.session

	if se.Debug {
		log.Trace.Println("Destory")
	}

	if se.Opts != nil && se.Opts.Destory != nil {
		res := (*se.Opts.Destory)(req, nodeid)

		return res
	} else {
		return errno.SUCCESS
	}

}

func doLookup(req FuseReq, nodeid uint64, entryOut *kernel.FuseEntryOut) int32 {
	lookupIn := (*req.Arg).(kernel.FuseLookupIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Lookup: %+v \n", lookupIn)
	}

	if se.Opts != nil && se.Opts.Lookup != nil {
		stat := syscall.Stat_t{}
		var generation uint64
		res := (*se.Opts.Lookup)(req, nodeid, lookupIn.Name, &stat, &generation)

		if res == errno.SUCCESS {
			entryOut.NodeId = stat.Ino
			entryOut.Generation = generation
			entryOut.AttrValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			entryOut.AttrValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			entryOut.EntryValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			entryOut.EntryValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			setFuseAttr(&entryOut.Attr, stat)
		}

		return res
	} else {
		return errno.ENOSYS
	}
}

func doForget(req FuseReq, nodeid uint64) {
	forgetIn := (*req.Arg).(kernel.FuseForgetIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Forget: %+v \n", forgetIn)
	}

	if se.Opts != nil && se.Opts.Forget != nil {
		(*se.Opts.Forget)(req, nodeid, forgetIn.Nlookup)
	}
}

func doGetattr(req FuseReq, nodeid uint64, attrOut *kernel.FuseAttrOut) int32 {

	getattrIn := (*req.Arg).(kernel.FuseGetattrIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Getattr: %+v \n", getattrIn)
	}

	if se.Opts != nil && se.Opts.Getattr != nil {

		stat := syscall.Stat_t{}

		res := (*se.Opts.Getattr)(req, nodeid, &stat)

		if res == errno.SUCCESS {
			attrOut.AttrValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			attrOut.AttrValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			attrOut.Dummp = getattrIn.Dummy
			setFuseAttr(&attrOut.Attr, stat)
		}

		return res
	} else {
		return errno.ENOSYS
	}

}

func doSetattr(req FuseReq, nodeid uint64, attrOut *kernel.FuseAttrOut) int32 {

	setattrIn := (*req.Arg).(kernel.FuseSetattrIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Setattr: %+v \n", setattrIn)
	}

	if se.Opts != nil && se.Opts.Setattr != nil {
		stat := syscall.Stat_t{}
		setattrInToStat(setattrIn, &stat)

		res := (*se.Opts.Setattr)(req, nodeid, &stat, setattrIn.Valid)

		if res == errno.SUCCESS {
			attrOut.AttrValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			attrOut.AttrValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			setFuseAttr(&attrOut.Attr, stat)
		}

		return res
	} else {
		return errno.ENOSYS
	}
}

func doReadlink(req FuseReq, nodeid uint64, attrOut *kernel.FuseReadlinkOut) int32 {
	se := req.session

	if se.Debug {
		log.Trace.Println("Readlink")
	}

	if se.Opts != nil && se.Opts.Readlink != nil {

		res := (*se.Opts.Readlink)(req, nodeid, &attrOut.Path)

		return res
	} else {
		return errno.ENOSYS
	}
}

func doMknod(req FuseReq, nodeid uint64, entryOut *kernel.FuseEntryOut) int32 {

	mknodIn := (*req.Arg).(kernel.FuseMknodIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Mknod: %+v \n", mknodIn)
	}

	if se.Opts != nil && se.Opts.Mknod != nil {

		stat, res := (*se.Opts.Mknod)(req, nodeid, mknodIn.Name, mknodIn.Mode, mknodIn.Rdev)

		if res == errno.SUCCESS {
			entryOut.NodeId = stat.Nodeid
			entryOut.Generation = stat.Generation
			entryOut.AttrValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			entryOut.AttrValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			entryOut.EntryValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			entryOut.EntryValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			setFuseAttr(&entryOut.Attr, stat.Stat)
		}

		return res
	} else {
		return errno.ENOSYS
	}
}

func doMkdir(req FuseReq, nodeid uint64, entryOut *kernel.FuseEntryOut) int32 {

	mkdirIn := (*req.Arg).(kernel.FuseMkdirIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Mkdir: %+v \n", mkdirIn)
	}

	if se.Opts != nil && se.Opts.Mkdir != nil {

		stat, res := (*se.Opts.Mkdir)(req, nodeid, mkdirIn.Name, mkdirIn.Mode)

		if res == errno.SUCCESS {
			entryOut.NodeId = stat.Nodeid
			entryOut.Generation = stat.Generation
			entryOut.AttrValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			entryOut.AttrValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			entryOut.EntryValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			entryOut.EntryValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			setFuseAttr(&entryOut.Attr, stat.Stat)
		}

		return res
	} else {
		return errno.ENOSYS
	}
}

func doUnlink(req FuseReq, nodeid uint64) int32 {

	unlinkIn := (*req.Arg).(kernel.FuseUnlinkIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Unlink: %+v \n", unlinkIn)
	}

	if se.Opts != nil && se.Opts.Unlink != nil {

		res := (*se.Opts.Unlink)(req, nodeid, unlinkIn.Path)

		return res
	} else {
		return errno.ENOSYS
	}
}

func doRmdir(req FuseReq, nodeid uint64) int32 {

	rmdirIn := (*req.Arg).(kernel.FuseRmdirIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Rmdir: %+v \n", rmdirIn)
	}

	if se.Opts != nil && se.Opts.Rmdir != nil {

		res := (*se.Opts.Rmdir)(req, nodeid, rmdirIn.Path)

		return res
	} else {
		return errno.ENOSYS
	}

}

func doSymlink(req FuseReq, nodeid uint64, entryOut *kernel.FuseEntryOut) int32 {

	symlinkIn := (*req.Arg).(kernel.FuseSymlinkIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Symlink: %+v \n", symlinkIn)
	}

	if se.Opts != nil && se.Opts.Symlink != nil {

		stat, res := (*se.Opts.Symlink)(req, nodeid, symlinkIn.LinkName, symlinkIn.Name)

		if res == errno.SUCCESS {
			entryOut.NodeId = stat.Nodeid
			entryOut.Generation = stat.Generation
			entryOut.AttrValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			entryOut.AttrValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			entryOut.EntryValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			entryOut.EntryValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			setFuseAttr(&entryOut.Attr, stat.Stat)
		}

		return res
	} else {
		return errno.ENOSYS
	}
}

func doRename(req FuseReq, nodeid uint64) int32 {

	renameIn := (*req.Arg).(kernel.FuseRenameIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Rename: %+v \n", renameIn)
	}

	if se.Opts != nil && se.Opts.Rename != nil {

		res := (*se.Opts.Rename)(req, nodeid, renameIn.OldName, renameIn.NewDir, renameIn.NewName)

		return res
	} else {
		return errno.ENOSYS
	}

}

func doRename2(req FuseReq, nodeid uint64) int32 {
	renameIn := (*req.Arg).(kernel.FuseRename2In)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Rename: %+v \n", renameIn)
	}

	if se.Opts != nil && se.Opts.Rename != nil {

		res := (*se.Opts.Rename)(req, nodeid, renameIn.OldName, renameIn.NewDir, renameIn.NewName)

		return res
	} else {
		return errno.ENOSYS
	}
}

func doLink(req FuseReq, nodeid uint64, entryOut *kernel.FuseEntryOut) int32 {

	linklIn := (*req.Arg).(kernel.FuseLinkIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Link: %+v \n", linklIn)
	}

	if se.Opts != nil && se.Opts.Link != nil {

		stat, res := (*se.Opts.Link)(req, linklIn.OldNodeid, nodeid, linklIn.NewName)

		if res == errno.SUCCESS {
			entryOut.NodeId = stat.Nodeid
			entryOut.Generation = stat.Generation
			entryOut.AttrValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			entryOut.AttrValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			entryOut.EntryValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			entryOut.EntryValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			setFuseAttr(&entryOut.Attr, stat.Stat)
		}

		return res
	} else {
		return errno.ENOSYS
	}
}

func doOpen(req FuseReq, nodeid uint64, openOut *kernel.FuseOpenOut) int32 {

	openIn := (*req.Arg).(kernel.FuseOpenIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Open: %+v \n", openIn)
	}

	fi := NewFuseFileInfo()
	fi.Flags = openIn.Flags

	if se.Opts != nil && se.Opts.Open != nil {

		res := (*se.Opts.Open)(req, nodeid, &fi)

		setOpenOut(openOut, fi)

		return res
	} else {

		setOpenOut(openOut, fi)
		return errno.SUCCESS
	}
}

func doRead(req FuseReq, nodeid uint64, readOut *kernel.FuseReadOut) int32 {

	readIn := (*req.Arg).(kernel.FuseReadIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Read: %v \n", readIn)
	}

	if se.Opts != nil && se.Opts.Read != nil {

		fi := NewFuseFileInfo()

		fi.Fh = readIn.Fh
		if req.session.ConnInfo.Minor >= 9 {
			fi.LockOwner = readIn.LockOwner
			fi.Flags = readIn.Flags
		}

		buf, res := (*se.Opts.Read)(req, nodeid, readIn.Size, readIn.Offset, fi)

		if res == errno.SUCCESS {
			readOut.Content = buf
		}

		return res
	} else {
		return errno.ENOSYS
	}

}

func doWrite(req FuseReq, nodeid uint64, writeOut *kernel.FuseWriteOut) int32 {

	writeIn := (*req.Arg).(kernel.FuseWriteIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Write: %v \n", writeIn)
	}

	if se.Opts != nil && se.Opts.Write != nil {

		fi := NewFuseFileInfo()

		fi.Fh = writeIn.Fh

		if writeIn.WriteFlags&1 != 0 {
			fi.Writepage = 1
		}

		if req.session.ConnInfo.Minor >= 9 {
			fi.LockOwner = writeIn.LockOwner
			fi.Flags = writeIn.Flags
		}

		size, res := (*se.Opts.Write)(req, nodeid, writeIn.Buf, writeIn.Offset, fi)

		writeOut.Size = size

		return res
	} else {
		return errno.ENOSYS
	}
}

func doFlush(req FuseReq, nodeid uint64) int32 {

	flushIn := (*req.Arg).(kernel.FuseFlushIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Flush: %+v \n", flushIn)
	}

	if se.Opts != nil && se.Opts.Flush != nil {

		fi := NewFuseFileInfo()

		fi.Fh = flushIn.Fh
		fi.Flush = 1

		if req.session.ConnInfo.Minor >= 9 {
			fi.LockOwner = flushIn.LockOwner
		}

		res := (*se.Opts.Flush)(req, nodeid, fi)

		return res
	} else {
		return errno.SUCCESS
	}
}

func doFsync(req FuseReq, nodeid uint64) int32 {

	fsyncIn := (*req.Arg).(kernel.FuseFsyncIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Fsync: %v \n", fsyncIn)
	}

	if se.Opts != nil && se.Opts.Fsync != nil {

		fi := NewFuseFileInfo()

		fi.Fh = fsyncIn.Fh

		datasync := fsyncIn.FsyncFlags & 1

		res := (*se.Opts.Fsync)(req, nodeid, datasync, fi)

		return res
	} else {
		return errno.ENOSYS
	}
}

func doOpendir(req FuseReq, nodeid uint64, openOut *kernel.FuseOpenOut) int32 {

	openIn := (*req.Arg).(kernel.FuseOpenIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Opendir: %+v \n", openIn)
	}

	fi := NewFuseFileInfo()
	fi.Flags = openIn.Flags
	if se.Opts != nil && se.Opts.Opendir != nil {

		res := (*se.Opts.Opendir)(req, nodeid, &fi)

		setOpenOut(openOut, fi)

		return res
	} else {
		setOpenOut(openOut, fi)

		return errno.SUCCESS
	}
}

func doReaddir(req FuseReq, nodeid uint64, readOut *kernel.FuseReadOut) int32 {

	readIn := (*req.Arg).(kernel.FuseReadIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Readdir: %+v \n", readIn)
	}

	if se.Opts != nil && se.Opts.Readdir != nil {

		fi := NewFuseFileInfo()

		fi.Fh = readIn.Fh

		dirList, res := (*se.Opts.Readdir)(req, nodeid, readIn.Size, readIn.Offset, fi)

		if res == errno.SUCCESS {

			buf := bytes.NewBuffer(nil)

			if dirList != nil && len(dirList) > 0 {

				var preOff uint64

				for i, _ := range dirList {
					dirb, _ := dirList[i].ToBinary(&preOff)
					// 判断是否超过readIn.Size的大小限制
					if uint32(buf.Len()+len(dirb)) < readIn.Size {
						buf.Write(dirb)
					} else {
						break
					}

				}

			}

			readOut.Content = buf.Bytes()
		}

		return res
	} else {
		return errno.ENOSYS
	}

}

func doRelease(req FuseReq, nodeid uint64) int32 {

	releaseIn := (*req.Arg).(kernel.FuseReleaseIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Release: %+v \n", releaseIn)
	}

	if se.Opts != nil && se.Opts.Release != nil {

		fi := NewFuseFileInfo()
		fi.Flags = releaseIn.Flags
		fi.Fh = releaseIn.Fh

		res := (*se.Opts.Release)(req, nodeid, fi)

		return res
	} else {
		return errno.SUCCESS
	}

}

func doReleasedir(req FuseReq, nodeid uint64) int32 {

	releasedirIn := (*req.Arg).(kernel.FuseReleaseIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Releasedir: %+v \n", releasedirIn)
	}

	if se.Opts != nil && se.Opts.Releasedir != nil {

		fi := NewFuseFileInfo()
		fi.Flags = releasedirIn.Flags
		fi.Fh = releasedirIn.Fh

		res := (*se.Opts.Releasedir)(req, nodeid, fi)

		return res
	} else {
		return errno.SUCCESS
	}

}

func doFsyncdir(req FuseReq, nodeid uint64) int32 {

	fsyncdirIn := (*req.Arg).(kernel.FuseFsyncIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Fsyncdir: %+v \n", fsyncdirIn)
	}

	if se.Opts != nil && se.Opts.Fsyncdir != nil {

		fi := NewFuseFileInfo()

		fi.Fh = fsyncdirIn.Fh

		datasync := fsyncdirIn.FsyncFlags & 1

		res := (*se.Opts.Fsyncdir)(req, nodeid, datasync, fi)

		return res

	} else {
		return errno.ENOSYS
	}
}

func doStatfs(req FuseReq, nodeid uint64, statfsOut *kernel.FuseStatfsOut) int32 {

	se := req.session

	if se.Debug {
		log.Trace.Println("Statfs")
	}

	var statfs = kernel.FuseStatfs{}
	if se.Opts != nil && se.Opts.Statfs != nil {

		res := (*se.Opts.Statfs)(req, nodeid, &statfs)

		statfsOut.St = statfs

		return res
	} else {
		statfs.NameLen = 255
		statfs.Bsize = 512

		statfsOut.St = statfs

		return errno.SUCCESS
	}
}

func doSetxattr(req FuseReq, nodeid uint64) int32 {

	setxattrIn := (*req.Arg).(kernel.FuseSetxattrIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Setxattr: %+v \n", setxattrIn)
	}

	if se.Opts != nil && se.Opts.Setxattr != nil {

		res := (*se.Opts.Setxattr)(req, nodeid, setxattrIn.Name, setxattrIn.Value, setxattrIn.Size, setxattrIn.Flags)

		return res
	} else {
		return errno.ENOSYS
	}

}

func doGetxattr(req FuseReq, nodeid uint64, getxattrOut *kernel.FuseGetxattrOut) int32 {

	getxattrIn := (*req.Arg).(kernel.FuseGetxattrIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Getxattr: %+v \n", getxattrIn)
	}

	if se.Opts != nil && se.Opts.Getxattr != nil {

		var value string

		res := (*se.Opts.Getxattr)(req, nodeid, getxattrIn.Name, getxattrIn.Size, &value)

		size := len(value)
		if getxattrIn.Size == 0 {
			getxattrOut.Size = uint32(size)
		} else {
			getxattrOut.Size = uint32(size)
			getxattrOut.Value = kernel.XattrVal(value)
		}

		return res
	} else {
		return errno.ENOSYS
	}
}

func doListxattr(req FuseReq, nodeid uint64, listxattrOut *kernel.FuseGetxattrOut) int32 {

	listxattrIn := (*req.Arg).(kernel.FuseGetxattrIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Listxattr: %+v \n", listxattrIn)
	}

	if se.Opts != nil && se.Opts.Listxattr != nil {

		var attrlist string

		res := (*se.Opts.Listxattr)(req, nodeid, listxattrIn.Size, &attrlist)

		size := len(attrlist)
		if listxattrIn.Size == 0 {
			listxattrOut.Size = uint32(size)
		} else {
			listxattrOut.Size = uint32(size)
			listxattrOut.Value = kernel.XattrVal(attrlist)
		}

		return res
	} else {
		return errno.ENOSYS
	}
}

func doRemovexattr(req FuseReq, nodeid uint64) int32 {

	removexattrIn := (*req.Arg).(kernel.FuseRemovexattrIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Removexattr: %+v \n", removexattrIn)
	}

	if se.Opts != nil && se.Opts.Removexattr != nil {

		res := (*se.Opts.Removexattr)(req, nodeid, removexattrIn.Name)

		return res
	} else {
		return errno.ENOSYS
	}
}

func doAccess(req FuseReq, nodeid uint64) int32 {

	accessIn := (*req.Arg).(kernel.FuseAccessIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Access: %+v \n", accessIn)
	}

	if se.Opts != nil && se.Opts.Access != nil {

		res := (*se.Opts.Access)(req, nodeid, accessIn.Mask)

		return res
	} else {
		return errno.ENOSYS
	}

}

func doCreate(req FuseReq, nodeid uint64, createOut *kernel.FuseCreateOut) int32 {

	createIn := (*req.Arg).(kernel.FuseCreateIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Create: %+v \n", createIn)
	}

	if se.Opts != nil && se.Opts.Create != nil {

		fi := NewFuseFileInfo()
		fi.Flags = createIn.Flags

		stat, res := (*se.Opts.Create)(req, nodeid, createIn.Name, createIn.Mode, &fi)

		if res == errno.SUCCESS {

			createOut.Entry.NodeId = stat.Nodeid
			createOut.Entry.Generation = stat.Generation
			createOut.Entry.AttrValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			createOut.Entry.AttrValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			createOut.Entry.EntryValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			createOut.Entry.EntryValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			setFuseAttr(&createOut.Entry.Attr, stat.Stat)

			setOpenOut(&createOut.Open, fi)
		}

		return res
	} else {
		return errno.ENOSYS
	}

}

func doGetlk(req FuseReq, nodeid uint64, getlkOut *kernel.FuseLkOut) int32 {

	getlkIn := (*req.Arg).(kernel.FuseLkIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Getlk: %+v \n", getlkIn)
	}

	if se.Opts != nil && se.Opts.Create != nil {

		fi := NewFuseFileInfo()
		var flock = syscall.Flock_t{}

		fi.Fh = getlkIn.Fh
		fi.LockOwner = getlkIn.Owner
		convertFuseFileLock(getlkIn.Lk, &flock)

		res := (*se.Opts.Getlk)(req, nodeid, fi, &flock)

		setFuseLkOut(flock, getlkOut)

		return res

	} else {
		return errno.ENOSYS
	}
}

func doSetlk(req FuseReq, nodeid uint64) int32 {
	return doSetlkCommon(req, nodeid, 0)
}

func doSetlkw(req FuseReq, nodeid uint64) int32 {
	return doSetlkCommon(req, nodeid, 1)
}

func doSetlkCommon(req FuseReq, nodeid uint64, lksleep int) int32 {

	setlkIn := (*req.Arg).(kernel.FuseLkIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Setlk: %+v \n", setlkIn)
	}

	if se.Opts != nil && se.Opts.Setlk != nil {
		var res int32

		fi := NewFuseFileInfo()
		var flock = syscall.Flock_t{}

		fi.Fh = setlkIn.Fh
		fi.LockOwner = setlkIn.Owner

		if (setlkIn.LkFlags & kernel.FUSE_LK_FLOCK) > 0 {

			op := 0

			switch setlkIn.Lk.Type {
			case syscall.F_RDLCK:
				op = syscall.LOCK_SH
			case syscall.F_WRLCK:
				op = syscall.LOCK_EX
			case syscall.F_UNLCK:
				op = syscall.LOCK_UN
			}

			if lksleep == 0 {
				op |= syscall.LOCK_NB
			}

			// TODO: flock

		} else {
			convertFuseFileLock(setlkIn.Lk, &flock)

			res = (*se.Opts.Setlk)(req, nodeid, fi, flock, lksleep)
		}

		return res
	} else {
		return errno.ENOSYS
	}

}

func doBmap(req FuseReq, nodeid uint64, bmapOut *kernel.FuseBmapOut) int32 {

	bmapIn := (*req.Arg).(kernel.FuseBmapIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Bmap: %+v \n", bmapIn)
	}

	if se.Opts != nil && se.Opts.Bmap != nil {

		idx := bmapIn.Block

		res := (*se.Opts.Bmap)(req, nodeid, bmapIn.BlockSize, &idx)

		bmapOut.Block = idx

		return res
	} else {
		return errno.ENOSYS
	}
}

func doIoctl(req FuseReq, nodeid uint64, ioctlOut *kernel.FuseIoctlOut) int32 {

	ioctlIn := (*req.Arg).(kernel.FuseIoctlIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Ioctl: %+v \n", ioctlIn)
	}

	flags := ioctlIn.Flags

	if (flags&kernel.FUSE_IOCTL_DIR) > 0 && (se.ConnInfo.Want&FUSE_CAP_IOCTL_DIR) == 0 {
		return errno.ENOSYS
	}

	if se.Opts != nil && se.Opts.Ioctl != nil {

		fi := NewFuseFileInfo()

		res := (*se.Opts.Ioctl)(req, nodeid, ioctlIn.Cmd, ioctlIn.Arg, fi, ioctlIn.InBuf, ioctlIn.OutSize, ioctlOut)

		return res
	} else {
		return errno.ENOSYS
	}

}

func doPoll(req FuseReq, nodeid uint64, pollOut *kernel.FusePollOut) int32 {

	pollIn := (*req.Arg).(kernel.FusePollIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Poll: %+v \n", pollIn)
	}

	if se.Opts != nil && se.Opts.Poll != nil {

		fi := NewFuseFileInfo()
		fi.Fh = pollIn.Fh
		fi.PollEvent = pollIn.Events

		var ph *FusePollhandle

		if (pollIn.Flags & kernel.FUSE_POLL_SCHEDULE_NOTIFY) > 0 {
			ph = &FusePollhandle{}

			ph.Kh = pollIn.Kh
			ph.Se = *se
		}

		res := (*se.Opts.Poll)(req, nodeid, fi, ph, &pollOut.Revents)

		return res
	} else {
		return errno.ENOSYS
	}

}

func doFallocate(req FuseReq, nodeid uint64) int32 {

	fallocateIn := (*req.Arg).(kernel.FuseFallocateIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Fallocate: %+v \n", fallocateIn)
	}

	if se.Opts != nil && se.Opts.Fallocate != nil {

		fi := NewFuseFileInfo()
		fi.Fh = fallocateIn.Fh

		res := (*se.Opts.Fallocate)(req, nodeid, fallocateIn.Mode, fallocateIn.Offset, fallocateIn.Length, fi)

		return res
	} else {
		return errno.ENOSYS
	}

}

func doForgetMulti(req FuseReq) {
	batchForgetIn := (*req.Arg).(kernel.FuseBatchForgetIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("BatchForget: %+v \n", batchForgetIn)
	}

	if se.Opts != nil && se.Opts.ForgetMulti != nil {
		(*se.Opts.ForgetMulti)(req, batchForgetIn.Count, batchForgetIn.NodeList)
	} else if se.Opts != nil && se.Opts.Forget != nil {
		for _, node := range batchForgetIn.NodeList {
			(*se.Opts.Forget)(req, node.Nodeid, node.Nlookup)
		}

	}
}

func doReaddirplus(req FuseReq, nodeid uint64, readOut *kernel.FuseReadOut) int32 {

	readIn := (*req.Arg).(kernel.FuseReadIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Readdirplus: %+v \n", readIn)
	}

	if se.Opts != nil && se.Opts.Readdirplus != nil {

		fi := NewFuseFileInfo()

		fi.Fh = readIn.Fh

		buf := bytes.NewBuffer(nil)

		res := (*se.Opts.Readdirplus)(req, nodeid, readIn.Size, readIn.Offset, fi, buf)

		if res == errno.SUCCESS {
			readOut.Content = buf.Bytes()
		}

		return res
	} else {
		return errno.ENOSYS
	}

}

const OFFSET_MAX = 0x7fffffffffffffff

func convertFuseFileLock(lock kernel.FuseFileLock, flock *syscall.Flock_t) {
	flock.Type = int16(lock.Type)
	flock.Whence = int16(os.SEEK_SET)
	flock.Start = int64(lock.Start)
	if lock.End == OFFSET_MAX {
		flock.Len = 0
	} else {
		flock.Len = int64(lock.End - lock.Start + 1)
	}
	flock.Pid = int32(lock.Pid)
}

func setFuseLkOut(flock syscall.Flock_t, getlkOut *kernel.FuseLkOut) {
	getlkOut.Lk.Type = uint32(flock.Type)
	if flock.Type != syscall.F_UNLCK {
		getlkOut.Lk.Start = uint64(flock.Start)
		if flock.Len == 0 {
			getlkOut.Lk.End = OFFSET_MAX
		} else {
			getlkOut.Lk.End = uint64(flock.Start + flock.Len - 1)
		}
	}
	getlkOut.Lk.Pid = uint32(flock.Pid)
}

func setFuseAttr(attr *kernel.FuseAttr, stat syscall.Stat_t) {
	attr.Ino = stat.Ino
	attr.Size = uint64(stat.Size)
	attr.Blocks = uint64(stat.Blocks)
	attr.Atime = uint64(stat.Atim.Sec)
	attr.AtimeNsec = uint32(stat.Atim.Nsec)
	attr.Mtime = uint64(stat.Mtim.Sec)
	attr.MtimeNsec = uint32(stat.Mtim.Nsec)
	attr.Ctime = uint64(stat.Ctim.Sec)
	attr.CtimeNsec = uint32(stat.Ctim.Nsec)
	attr.Mode = stat.Mode
	attr.Nlink = uint32(stat.Nlink)
	attr.Uid = stat.Uid
	attr.Gid = stat.Gid
	attr.Rdev = uint32(stat.Rdev)
	attr.Blksize = uint32(stat.Blksize)
}

func setattrInToStat(setattrIn kernel.FuseSetattrIn, stat *syscall.Stat_t) {

	stat.Size = int64(setattrIn.Size)
	stat.Atim.Sec = int64(setattrIn.Atime)
	stat.Atim.Nsec = int64(setattrIn.AtimeNsec)
	stat.Mtim.Sec = int64(setattrIn.Mtime)
	stat.Mtim.Nsec = int64(setattrIn.MtimeNsec)
	stat.Ctim.Sec = int64(setattrIn.Ctime)
	stat.Ctim.Nsec = int64(setattrIn.CtimeNsec)
	stat.Mode = setattrIn.Mode
	stat.Uid = setattrIn.Uid
	stat.Gid = setattrIn.Gid

}

func setOpenOut(openOut *kernel.FuseOpenOut, fi FuseFileInfo) {
	openOut.Fh = fi.Fh

	if fi.DirectIo > 0 {
		openOut.OpenFlags |= FOPEN_DIRECT_IO
	}

	if fi.KeepCache > 0 {
		openOut.OpenFlags |= FOPEN_KEEP_CACHE
	}

	if fi.Nonseekable > 0 {
		openOut.OpenFlags |= FOPEN_NONSEEKABLE
	}
}
