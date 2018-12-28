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

func doInit(req Req, initOut *kernel.FuseInitOut) int32 {

	se := req.session
	initIn := (*req.Arg).(kernel.FuseInitIn)
	if se.Debug {
		log.Trace.Printf("INIT: %+v \n", initIn)
	}

	bufsize := se.bufsize

	se.connInfo.Major = initIn.Major
	se.connInfo.Minor = initIn.Minor
	se.connInfo.MaxReadahead = initIn.MaxReadahead

	if bufsize < kernel.FuseMinReadBuffer {
		log.Warning.Printf("fuse: warning: buffer size too small: %d\n", bufsize)
		bufsize = kernel.FuseMinReadBuffer
	}

	bufsize -= 4096
	if uint32(bufsize) < se.connInfo.MaxWrite {
		se.connInfo.MaxWrite = uint32(bufsize)
	}

	initOut.Major = kernel.FuseKernelVersion
	initOut.Minor = kernel.FuseKernelMinorVersion
	initOut.MaxReadahead = se.connInfo.MaxReadahead
	initOut.MaxWrite = se.connInfo.MaxWrite
	initOut.TimeGran = se.connInfo.TimeGran
	initOut.CongestionThreshold = se.connInfo.CongestionThreshold

	// To remember what fuse kernel can do
	if initIn.Flags&kernel.FuseAsyncRead > 0 {
		se.connInfo.Capable |= kernel.FuseAsyncRead
	}
	if initIn.Flags&kernel.FusePosixLocks > 0 {
		se.connInfo.Capable |= kernel.FusePosixLocks
	}
	if initIn.Flags&kernel.FuseFileOps > 0 {
		se.connInfo.Capable |= kernel.FuseFileOps
	}
	if initIn.Flags&kernel.FuseAtomicOTrunc > 0 {
		se.connInfo.Capable |= kernel.FuseAtomicOTrunc
	}
	if initIn.Flags&kernel.FuseExportSupport > 0 {
		se.connInfo.Capable |= kernel.FuseExportSupport
	}
	if initIn.Flags&kernel.FuseBigWrites > 0 {
		se.connInfo.Capable |= kernel.FuseBigWrites
	}
	if initIn.Flags&kernel.FuseDontMask > 0 {
		se.connInfo.Capable |= kernel.FuseDontMask
	}
	if initIn.Flags&kernel.FuseSpliceWrite > 0 {
		se.connInfo.Capable |= kernel.FuseSpliceWrite
	}
	if initIn.Flags&kernel.FuseSpliceMove > 0 {
		se.connInfo.Capable |= kernel.FuseSpliceMove
	}
	if initIn.Flags&kernel.FuseSpliceRead > 0 {
		se.connInfo.Capable |= kernel.FuseSpliceRead
	}
	if initIn.Flags&kernel.FuseFlockLocks > 0 {
		se.connInfo.Capable |= kernel.FuseFlockLocks
	}
	if initIn.Flags&kernel.FuseHasIoCtlDir > 0 {
		se.connInfo.Capable |= kernel.FuseHasIoCtlDir
	}
	if initIn.Flags&kernel.FuseAutoInvalData > 0 {
		se.connInfo.Capable |= kernel.FuseAutoInvalData
	}
	if initIn.Flags&kernel.FuseDoReaddirplus > 0 {
		se.connInfo.Capable |= kernel.FuseDoReaddirplus
	}
	if initIn.Flags&kernel.FuseReaddirplusAuto > 0 {
		se.connInfo.Capable |= kernel.FuseReaddirplusAuto
	}
	if initIn.Flags&kernel.FuseAsyncDio > 0 {
		se.connInfo.Capable |= kernel.FuseAsyncDio
	}
	if initIn.Flags&kernel.FuseWritebackCache > 0 {
		se.connInfo.Capable |= kernel.FuseWritebackCache
	}
	if initIn.Flags&kernel.FuseNoOpenSupport > 0 {
		se.connInfo.Capable |= kernel.FuseNoOpenSupport
	}
	if initIn.Flags&kernel.FuseParallelDirops > 0 {
		se.connInfo.Capable |= kernel.FuseParallelDirops
	}
	if initIn.Flags&kernel.FuseHandleKillPriv > 0 {
		se.connInfo.Capable |= kernel.FuseHandleKillPriv
	}
	if initIn.Flags&kernel.FusePosixACL > 0 {
		se.connInfo.Capable |= kernel.FusePosixACL
	}

	// Default settings for modern filesystems.
	// TODO: support write_buf, flock
	if (se.connInfo.Capable & FuseCapAsyncRead) > 0 {
		se.connInfo.Want |= FuseCapAsyncRead
	}
	if (se.connInfo.Capable & FuseCapParallelDirops) > 0 {
		se.connInfo.Want |= FuseCapParallelDirops
	}
	if (se.connInfo.Capable & FuseCapAutoInvalData) > 0 {
		se.connInfo.Want |= FuseCapAutoInvalData
	}
	if (se.connInfo.Capable & FuseCapHandleKillpriv) > 0 {
		se.connInfo.Want |= FuseCapHandleKillpriv
	}
	if (se.connInfo.Capable & FuseCapAsyncDIO) > 0 {
		se.connInfo.Want |= FuseCapAsyncDIO
	}
	if (se.connInfo.Capable & FuseCapIoctlDir) > 0 {
		se.connInfo.Want |= FuseCapIoctlDir
	}
	if (se.connInfo.Capable & FuseCapAtomicOTrunc) > 0 {
		se.connInfo.Want |= FuseCapAtomicOTrunc
	}
	if se.Opts.Getlk != nil && se.Opts.Setlk != nil {
		se.connInfo.Want |= FuseCapPosixLocks
	}
	if se.Opts.Readdirplus != nil {
		se.connInfo.Want |= FuseCapReaddirplus
		se.connInfo.Want |= FuseCapReaddirplusAuto
	}

	// To set what we want fuse kenel to do
	if se.connInfo.Want&kernel.FuseAsyncRead > 0 {
		initOut.Flags |= kernel.FuseAsyncRead
	}
	if se.connInfo.Want&kernel.FusePosixLocks > 0 {
		initOut.Flags |= kernel.FusePosixLocks
	}
	if se.connInfo.Want&kernel.FuseFileOps > 0 {
		initOut.Flags |= kernel.FuseFileOps
	}
	if se.connInfo.Want&kernel.FuseAtomicOTrunc > 0 {
		initOut.Flags |= kernel.FuseAtomicOTrunc
	}
	if se.connInfo.Want&kernel.FuseExportSupport > 0 {
		initOut.Flags |= kernel.FuseExportSupport
	}
	if se.connInfo.Want&kernel.FuseBigWrites > 0 {
		initOut.Flags |= kernel.FuseBigWrites
	}
	if se.connInfo.Want&kernel.FuseDontMask > 0 {
		initOut.Flags |= kernel.FuseDontMask
	}
	if se.connInfo.Want&kernel.FuseSpliceWrite > 0 {
		initOut.Flags |= kernel.FuseSpliceWrite
	}
	if se.connInfo.Want&kernel.FuseSpliceMove > 0 {
		initOut.Flags |= kernel.FuseSpliceMove
	}
	if se.connInfo.Want&kernel.FuseSpliceRead > 0 {
		initOut.Flags |= kernel.FuseSpliceRead
	}
	if se.connInfo.Want&kernel.FuseFlockLocks > 0 {
		initOut.Flags |= kernel.FuseFlockLocks
	}
	if se.connInfo.Want&kernel.FuseHasIoCtlDir > 0 {
		initOut.Flags |= kernel.FuseHasIoCtlDir
	}
	if se.connInfo.Want&kernel.FuseAutoInvalData > 0 {
		initOut.Flags |= kernel.FuseAutoInvalData
	}
	if se.connInfo.Want&kernel.FuseDoReaddirplus > 0 {
		initOut.Flags |= kernel.FuseDoReaddirplus
	}
	if se.connInfo.Want&kernel.FuseReaddirplusAuto > 0 {
		initOut.Flags |= kernel.FuseReaddirplusAuto
	}
	if se.connInfo.Want&kernel.FuseAsyncDio > 0 {
		initOut.Flags |= kernel.FuseAsyncDio
	}
	if se.connInfo.Want&kernel.FuseWritebackCache > 0 {
		initOut.Flags |= kernel.FuseWritebackCache
	}
	if se.connInfo.Want&kernel.FuseNoOpenSupport > 0 {
		initOut.Flags |= kernel.FuseNoOpenSupport
	}
	if se.connInfo.Want&kernel.FuseParallelDirops > 0 {
		initOut.Flags |= kernel.FuseParallelDirops
	}
	if se.connInfo.Want&kernel.FuseHandleKillPriv > 0 {
		initOut.Flags |= kernel.FuseHandleKillPriv
	}
	if se.connInfo.Want&kernel.FusePosixACL > 0 {
		initOut.Flags |= kernel.FusePosixACL
	}

	if se.Opts != nil && se.Opts.Init != nil {
		userdata := (*se.Opts.Init)(se.connInfo)
		se.userdata = userdata

	}

	return errno.SUCCESS
}

func doDestory(req Req) {
	se := req.session

	if se.Debug {
		log.Trace.Println("Destory")
	}

	if se.Opts != nil && se.Opts.Destory != nil {
		(*se.Opts.Destory)(se.userdata)

	}

}

func doLookup(req Req, nodeid uint64, entryOut *kernel.FuseEntryOut) int32 {
	lookupIn := (*req.Arg).(kernel.FuseLookupIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Lookup: %+v \n", lookupIn)
	}

	if se.Opts != nil && se.Opts.Lookup != nil {

		var fsStat *FileStat

		fsStat, res = (*se.Opts.Lookup)(req, nodeid, lookupIn.Name)

		if res == errno.SUCCESS {
			entryOut.NodeID = fsStat.Nodeid
			entryOut.Generation = fsStat.Generation
			entryOut.AttrValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			entryOut.AttrValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			entryOut.EntryValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			entryOut.EntryValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			setFuseAttr(&entryOut.Attr, fsStat.Stat)
		}

	}

	return res
}

func doForget(req Req, nodeid uint64) {
	forgetIn := (*req.Arg).(kernel.FuseForgetIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("Forget: %+v \n", forgetIn)
	}

	if se.Opts != nil && se.Opts.Forget != nil {
		(*se.Opts.Forget)(req, nodeid, forgetIn.Nlookup)
	}
}

func doGetattr(req Req, nodeid uint64, attrOut *kernel.FuseAttrOut) int32 {

	getattrIn := (*req.Arg).(kernel.FuseGetattrIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Getattr: %+v \n", getattrIn)
	}

	if se.Opts != nil && se.Opts.Getattr != nil {

		var fsStat *FileStat
		fsStat, res = (*se.Opts.Getattr)(req, nodeid)

		if res == errno.SUCCESS {
			attrOut.AttrValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			attrOut.AttrValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			attrOut.Dummp = getattrIn.Dummy
			setFuseAttr(&attrOut.Attr, fsStat.Stat)
		}

	}

	return res
}

func doSetattr(req Req, nodeid uint64, attrOut *kernel.FuseAttrOut) int32 {

	setattrIn := (*req.Arg).(kernel.FuseSetattrIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Setattr: %+v \n", setattrIn)
	}

	if se.Opts != nil && se.Opts.Setattr != nil {
		fsStat := FileStat{}

		setattrInToStat(setattrIn, &fsStat.Stat)

		res = (*se.Opts.Setattr)(req, nodeid, fsStat, setattrIn.Valid)

		if res == errno.SUCCESS {
			attrOut.AttrValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			attrOut.AttrValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			setFuseAttr(&attrOut.Attr, fsStat.Stat)
		}

	}

	return res
}

func doReadlink(req Req, nodeid uint64, attrOut *kernel.FuseReadlinkOut) int32 {

	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Println("Readlink")
	}

	if se.Opts != nil && se.Opts.Readlink != nil {

		var path string

		path, res = (*se.Opts.Readlink)(req, nodeid)

		if res == errno.SUCCESS {
			attrOut.Path = path
		}

	}

	return res
}

func doMknod(req Req, nodeid uint64, entryOut *kernel.FuseEntryOut) int32 {

	mknodIn := (*req.Arg).(kernel.FuseMknodIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Mknod: %+v \n", mknodIn)
	}

	if se.Opts != nil && se.Opts.Mknod != nil {

		var stat *FileStat

		stat, res = (*se.Opts.Mknod)(req, nodeid, mknodIn.Name, mknodIn.Mode, mknodIn.Rdev)

		if res == errno.SUCCESS {
			entryOut.NodeID = stat.Nodeid
			entryOut.Generation = stat.Generation
			entryOut.AttrValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			entryOut.AttrValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			entryOut.EntryValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			entryOut.EntryValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			setFuseAttr(&entryOut.Attr, stat.Stat)
		}

	}

	return res
}

func doMkdir(req Req, nodeid uint64, entryOut *kernel.FuseEntryOut) int32 {

	mkdirIn := (*req.Arg).(kernel.FuseMkdirIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Mkdir: %+v \n", mkdirIn)
	}

	if se.Opts != nil && se.Opts.Mkdir != nil {

		var stat *FileStat

		stat, res = (*se.Opts.Mkdir)(req, nodeid, mkdirIn.Name, mkdirIn.Mode)

		if res == errno.SUCCESS {
			entryOut.NodeID = stat.Nodeid
			entryOut.Generation = stat.Generation
			entryOut.AttrValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			entryOut.AttrValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			entryOut.EntryValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			entryOut.EntryValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			setFuseAttr(&entryOut.Attr, stat.Stat)
		}

	}

	return res
}

func doUnlink(req Req, nodeid uint64) int32 {

	unlinkIn := (*req.Arg).(kernel.FuseUnlinkIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Unlink: %+v \n", unlinkIn)
	}

	if se.Opts != nil && se.Opts.Unlink != nil {

		res = (*se.Opts.Unlink)(req, nodeid, unlinkIn.Path)

	}

	return res
}

func doRmdir(req Req, nodeid uint64) int32 {

	rmdirIn := (*req.Arg).(kernel.FuseRmdirIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Rmdir: %+v \n", rmdirIn)
	}

	if se.Opts != nil && se.Opts.Rmdir != nil {

		res = (*se.Opts.Rmdir)(req, nodeid, rmdirIn.Path)

	}

	return res
}

func doSymlink(req Req, nodeid uint64, entryOut *kernel.FuseEntryOut) int32 {

	symlinkIn := (*req.Arg).(kernel.FuseSymlinkIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Symlink: %+v \n", symlinkIn)
	}

	if se.Opts != nil && se.Opts.Symlink != nil {

		var stat *FileStat

		stat, res = (*se.Opts.Symlink)(req, nodeid, symlinkIn.LinkName, symlinkIn.Name)

		if res == errno.SUCCESS {
			entryOut.NodeID = stat.Nodeid
			entryOut.Generation = stat.Generation
			entryOut.AttrValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			entryOut.AttrValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			entryOut.EntryValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			entryOut.EntryValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			setFuseAttr(&entryOut.Attr, stat.Stat)
		}

	}

	return res
}

func doRename(req Req, nodeid uint64) int32 {

	renameIn := (*req.Arg).(kernel.FuseRenameIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Rename: %+v \n", renameIn)
	}

	if se.Opts != nil && se.Opts.Rename != nil {

		res = (*se.Opts.Rename)(req, nodeid, renameIn.OldName, renameIn.NewDir, renameIn.NewName)

	}

	return res

}

func doRename2(req Req, nodeid uint64) int32 {
	renameIn := (*req.Arg).(kernel.FuseRename2In)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Rename: %+v \n", renameIn)
	}

	if se.Opts != nil && se.Opts.Rename != nil {

		res = (*se.Opts.Rename)(req, nodeid, renameIn.OldName, renameIn.NewDir, renameIn.NewName)

	}

	return res
}

func doLink(req Req, nodeid uint64, entryOut *kernel.FuseEntryOut) int32 {

	linklIn := (*req.Arg).(kernel.FuseLinkIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Link: %+v \n", linklIn)
	}

	if se.Opts != nil && se.Opts.Link != nil {

		var stat *FileStat

		stat, res = (*se.Opts.Link)(req, linklIn.OldNodeid, nodeid, linklIn.NewName)

		if res == errno.SUCCESS {
			entryOut.NodeID = stat.Nodeid
			entryOut.Generation = stat.Generation
			entryOut.AttrValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			entryOut.AttrValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			entryOut.EntryValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			entryOut.EntryValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			setFuseAttr(&entryOut.Attr, stat.Stat)
		}

	}

	return res
}

func doOpen(req Req, nodeid uint64, openOut *kernel.FuseOpenOut) int32 {

	openIn := (*req.Arg).(kernel.FuseOpenIn)
	se := req.session
	var res int32 = errno.SUCCESS

	if se.Debug {
		log.Trace.Printf("Open: %+v \n", openIn)
	}

	fi := NewFuseFileInfo()
	fi.Flags = openIn.Flags

	if se.Opts != nil && se.Opts.Open != nil {

		res = (*se.Opts.Open)(req, nodeid, &fi)
	}

	setOpenOut(openOut, fi)

	return res
}

func doRead(req Req, nodeid uint64, readOut *kernel.FuseReadOut) int32 {

	readIn := (*req.Arg).(kernel.FuseReadIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Read: %v \n", readIn)
	}

	if se.Opts != nil && se.Opts.Read != nil {

		fi := NewFuseFileInfo()

		fi.Fh = readIn.Fh
		if req.session.connInfo.Minor >= 9 {
			fi.LockOwner = readIn.LockOwner
			fi.Flags = readIn.Flags
		}

		var buf []byte

		buf, res = (*se.Opts.Read)(req, nodeid, readIn.Size, readIn.Offset, fi)

		if res == errno.SUCCESS {
			readOut.Content = buf
		}

	}

	return res
}

func doWrite(req Req, nodeid uint64, writeOut *kernel.FuseWriteOut) int32 {

	writeIn := (*req.Arg).(kernel.FuseWriteIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Write: %v \n", writeIn)
	}

	if se.Opts != nil && se.Opts.Write != nil {

		fi := NewFuseFileInfo()

		fi.Fh = writeIn.Fh

		if writeIn.WriteFlags&1 != 0 {
			fi.Writepage = 1
		}

		if req.session.connInfo.Minor >= 9 {
			fi.LockOwner = writeIn.LockOwner
			fi.Flags = writeIn.Flags
		}

		var size uint32
		size, res = (*se.Opts.Write)(req, nodeid, writeIn.Buf, writeIn.Offset, fi)
		writeOut.Size = size
	}

	return res
}

func doFlush(req Req, nodeid uint64) int32 {

	flushIn := (*req.Arg).(kernel.FuseFlushIn)
	se := req.session
	var res int32 = errno.SUCCESS

	if se.Debug {
		log.Trace.Printf("Flush: %+v \n", flushIn)
	}

	if se.Opts != nil && se.Opts.Flush != nil {

		fi := NewFuseFileInfo()

		fi.Fh = flushIn.Fh
		fi.Flush = 1

		if req.session.connInfo.Minor >= 9 {
			fi.LockOwner = flushIn.LockOwner
		}

		res = (*se.Opts.Flush)(req, nodeid, fi)

	}

	return res
}

func doFsync(req Req, nodeid uint64) int32 {

	fsyncIn := (*req.Arg).(kernel.FuseFsyncIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Fsync: %v \n", fsyncIn)
	}

	if se.Opts != nil && se.Opts.Fsync != nil {

		fi := NewFuseFileInfo()

		fi.Fh = fsyncIn.Fh

		datasync := fsyncIn.FsyncFlags & 1

		res = (*se.Opts.Fsync)(req, nodeid, datasync, fi)

	}

	return res
}

func doOpendir(req Req, nodeid uint64, openOut *kernel.FuseOpenOut) int32 {

	openIn := (*req.Arg).(kernel.FuseOpenIn)
	se := req.session
	var res int32 = errno.SUCCESS

	if se.Debug {
		log.Trace.Printf("Opendir: %+v \n", openIn)
	}

	fi := NewFuseFileInfo()
	fi.Flags = openIn.Flags
	if se.Opts != nil && se.Opts.Opendir != nil {

		res = (*se.Opts.Opendir)(req, nodeid, &fi)

	}
	setOpenOut(openOut, fi)

	return res
}

func doReaddir(req Req, nodeid uint64, readOut *kernel.FuseReadOut) int32 {

	readIn := (*req.Arg).(kernel.FuseReadIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Readdir: %+v \n", readIn)
	}

	if se.Opts != nil && se.Opts.Readdir != nil {

		fi := NewFuseFileInfo()
		fi.Fh = readIn.Fh

		var dirList []Dirent

		dirList, res = (*se.Opts.Readdir)(req, nodeid, readIn.Size, readIn.Offset, fi)

		if res == errno.SUCCESS {

			buf := bytes.NewBuffer(nil)

			if dirList != nil && len(dirList) > 0 {

				var preOff uint64

				for _, val := range dirList {

					dir := val

					dirent := kernel.FuseDirent(dir)

					dirb, _ := dirent.ToBinary(&preOff)
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

	}

	return res
}

func doRelease(req Req, nodeid uint64) int32 {

	releaseIn := (*req.Arg).(kernel.FuseReleaseIn)
	se := req.session
	var res int32 = errno.SUCCESS

	if se.Debug {
		log.Trace.Printf("Release: %+v \n", releaseIn)
	}

	if se.Opts != nil && se.Opts.Release != nil {

		fi := NewFuseFileInfo()
		fi.Flags = releaseIn.Flags
		fi.Fh = releaseIn.Fh

		res = (*se.Opts.Release)(req, nodeid, fi)

	}

	return res
}

func doReleasedir(req Req, nodeid uint64) int32 {

	releasedirIn := (*req.Arg).(kernel.FuseReleaseIn)
	se := req.session
	var res int32 = errno.SUCCESS

	if se.Debug {
		log.Trace.Printf("Releasedir: %+v \n", releasedirIn)
	}

	if se.Opts != nil && se.Opts.Releasedir != nil {

		fi := NewFuseFileInfo()
		fi.Flags = releasedirIn.Flags
		fi.Fh = releasedirIn.Fh

		res = (*se.Opts.Releasedir)(req, nodeid, fi)

	}

	return res
}

func doFsyncdir(req Req, nodeid uint64) int32 {

	fsyncdirIn := (*req.Arg).(kernel.FuseFsyncIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Fsyncdir: %+v \n", fsyncdirIn)
	}

	if se.Opts != nil && se.Opts.Fsyncdir != nil {

		fi := NewFuseFileInfo()

		fi.Fh = fsyncdirIn.Fh

		datasync := fsyncdirIn.FsyncFlags & 1

		res = (*se.Opts.Fsyncdir)(req, nodeid, datasync, fi)

	}

	return res
}

func doStatfs(req Req, nodeid uint64, statfsOut *kernel.FuseStatfsOut) int32 {

	se := req.session
	var res int32 = errno.SUCCESS

	if se.Debug {
		log.Trace.Println("Statfs")
	}

	if se.Opts != nil && se.Opts.Statfs != nil {

		var statfs *Statfs
		statfs, res = (*se.Opts.Statfs)(req, nodeid)

		if res == errno.SUCCESS {
			statfsOut.St = kernel.FuseStatfs(*statfs)
		}

	} else {
		statfs := kernel.FuseStatfs{}
		statfs.NameLen = 255
		statfs.Bsize = 512

		statfsOut.St = statfs

	}

	return res
}

func doSetxattr(req Req, nodeid uint64) int32 {

	setxattrIn := (*req.Arg).(kernel.FuseSetxattrIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Setxattr: %+v \n", setxattrIn)
	}

	if se.Opts != nil && se.Opts.Setxattr != nil {

		res = (*se.Opts.Setxattr)(req, nodeid, setxattrIn.Name, setxattrIn.Value, setxattrIn.Flags)

	}

	return res
}

func doGetxattr(req Req, nodeid uint64, getxattrOut *kernel.FuseGetxattrOut) int32 {

	getxattrIn := (*req.Arg).(kernel.FuseGetxattrIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Getxattr: %+v \n", getxattrIn)
	}

	if se.Opts != nil && se.Opts.Getxattr != nil {

		var value string
		value, res = (*se.Opts.Getxattr)(req, nodeid, getxattrIn.Name, getxattrIn.Size)

		size := len(value)
		getxattrOut.Size = uint32(size)
		getxattrOut.Value = kernel.XattrVal(value)

	}

	return res
}

func doListxattr(req Req, nodeid uint64, listxattrOut *kernel.FuseGetxattrOut) int32 {

	listxattrIn := (*req.Arg).(kernel.FuseGetxattrIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Listxattr: %+v \n", listxattrIn)
	}

	if se.Opts != nil && se.Opts.Listxattr != nil {

		var attrlist string
		attrlist, res = (*se.Opts.Listxattr)(req, nodeid, listxattrIn.Size)

		size := len(attrlist)
		if listxattrIn.Size == 0 {
			listxattrOut.Size = uint32(size)
		} else {
			listxattrOut.Size = uint32(size)
			listxattrOut.Value = kernel.XattrVal(attrlist)
		}

	}

	return res
}

func doRemovexattr(req Req, nodeid uint64) int32 {

	removexattrIn := (*req.Arg).(kernel.FuseRemovexattrIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Removexattr: %+v \n", removexattrIn)
	}

	if se.Opts != nil && se.Opts.Removexattr != nil {

		res = (*se.Opts.Removexattr)(req, nodeid, removexattrIn.Name)

	}

	return res
}

func doAccess(req Req, nodeid uint64) int32 {

	accessIn := (*req.Arg).(kernel.FuseAccessIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Access: %+v \n", accessIn)
	}

	if se.Opts != nil && se.Opts.Access != nil {

		res = (*se.Opts.Access)(req, nodeid, accessIn.Mask)

	}

	return res
}

func doCreate(req Req, nodeid uint64, createOut *kernel.FuseCreateOut) int32 {

	createIn := (*req.Arg).(kernel.FuseCreateIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Create: %+v \n", createIn)
	}

	if se.Opts != nil && se.Opts.Create != nil {

		fi := NewFuseFileInfo()
		fi.Flags = createIn.Flags

		var stat *FileStat

		stat, res = (*se.Opts.Create)(req, nodeid, createIn.Name, createIn.Mode, &fi)

		if res == errno.SUCCESS {

			createOut.Entry.NodeID = stat.Nodeid
			createOut.Entry.Generation = stat.Generation
			createOut.Entry.AttrValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			createOut.Entry.AttrValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			createOut.Entry.EntryValid = common.CalcTimeoutSec(se.FuseConfig.AttrTimeout)
			createOut.Entry.EntryValidNsec = common.CalcTimeoutNsec(se.FuseConfig.AttrTimeout)
			setFuseAttr(&createOut.Entry.Attr, stat.Stat)

			setOpenOut(&createOut.Open, fi)
		}

	}

	return res
}

func doGetlk(req Req, nodeid uint64, getlkOut *kernel.FuseLkOut) int32 {

	getlkIn := (*req.Arg).(kernel.FuseLkIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Getlk: %+v \n", getlkIn)
	}

	if se.Opts != nil && se.Opts.Create != nil {

		fi := NewFuseFileInfo()
		var flock = Flock{}

		fi.Fh = getlkIn.Fh
		fi.LockOwner = getlkIn.Owner
		convertFuseFileLock(getlkIn.Lk, &flock)

		res = (*se.Opts.Getlk)(req, nodeid, fi, &flock)

		setFuseLkOut(flock, getlkOut)

	}

	return res
}

func doSetlk(req Req, nodeid uint64) int32 {
	return doSetlkCommon(req, nodeid, 0)
}

func doSetlkw(req Req, nodeid uint64) int32 {
	return doSetlkCommon(req, nodeid, 1)
}

func doSetlkCommon(req Req, nodeid uint64, lksleep int) int32 {

	setlkIn := (*req.Arg).(kernel.FuseLkIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Setlk: %+v \n", setlkIn)
	}

	if se.Opts != nil && se.Opts.Setlk != nil {

		fi := NewFuseFileInfo()
		var flock = Flock{}

		fi.Fh = setlkIn.Fh
		fi.LockOwner = setlkIn.Owner

		if (setlkIn.LkFlags & kernel.FuseLkFlock) > 0 {

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

	}

	return res
}

func doBmap(req Req, nodeid uint64, bmapOut *kernel.FuseBmapOut) int32 {

	bmapIn := (*req.Arg).(kernel.FuseBmapIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Bmap: %+v \n", bmapIn)
	}

	if se.Opts != nil && se.Opts.Bmap != nil {

		idx := bmapIn.Block

		res = (*se.Opts.Bmap)(req, nodeid, bmapIn.BlockSize, &idx)

		bmapOut.Block = idx

	}

	return res
}

func doIoctl(req Req, nodeid uint64, ioctlOut *kernel.FuseIoctlOut) int32 {

	ioctlIn := (*req.Arg).(kernel.FuseIoctlIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Ioctl: %+v \n", ioctlIn)
	}

	flags := ioctlIn.Flags

	if (flags&kernel.FuseIoctlDir) > 0 && (se.connInfo.Want&FuseCapIoctlDir) == 0 {
		return errno.ENOSYS
	}

	if se.Opts != nil && se.Opts.Ioctl != nil {

		fi := NewFuseFileInfo()

		var ioctl Ioctl
		ioctl, res = (*se.Opts.Ioctl)(req, nodeid, ioctlIn.Cmd, ioctlIn.Arg, fi, ioctlIn.InBuf, ioctlIn.OutSize)

		ioctlOut.Result = ioctl.Result
		ioctlOut.Flags = ioctl.Flags
		ioctlOut.InIovs = ioctl.InIovs
		ioctlOut.OutIovs = ioctl.OutIovs

	}

	return res

}

func doPoll(req Req, nodeid uint64, pollOut *kernel.FusePollOut) int32 {

	pollIn := (*req.Arg).(kernel.FusePollIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Poll: %+v \n", pollIn)
	}

	if se.Opts != nil && se.Opts.Poll != nil {

		fi := NewFuseFileInfo()
		fi.Fh = pollIn.Fh
		fi.PollEvent = pollIn.Events

		var ph *Pollhandle

		if (pollIn.Flags & kernel.FusePollScheduleNotify) > 0 {
			ph = &Pollhandle{}

			ph.Kh = pollIn.Kh
			ph.Se = *se
		}

		var revents uint32
		revents, res = (*se.Opts.Poll)(req, nodeid, fi, ph)

		pollOut.Revents = revents

	}

	return res
}

func doFallocate(req Req, nodeid uint64) int32 {

	fallocateIn := (*req.Arg).(kernel.FuseFallocateIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Fallocate: %+v \n", fallocateIn)
	}

	if se.Opts != nil && se.Opts.Fallocate != nil {

		fi := NewFuseFileInfo()
		fi.Fh = fallocateIn.Fh

		res = (*se.Opts.Fallocate)(req, nodeid, fallocateIn.Mode, fallocateIn.Offset, fallocateIn.Length, fi)

	}

	return res
}

func doForgetMulti(req Req) {
	batchForgetIn := (*req.Arg).(kernel.FuseBatchForgetIn)
	se := req.session

	if se.Debug {
		log.Trace.Printf("BatchForget: %+v \n", batchForgetIn)
	}

	if se.Opts != nil && se.Opts.ForgetMulti != nil {

		nodelist := make([]ForgetOne, batchForgetIn.Count)

		for i, val := range batchForgetIn.NodeList {

			node := val
			nodelist[i] = ForgetOne(node)
		}

		(*se.Opts.ForgetMulti)(req, nodelist)
	} else if se.Opts != nil && se.Opts.Forget != nil {
		for _, node := range batchForgetIn.NodeList {
			(*se.Opts.Forget)(req, node.Nodeid, node.Nlookup)
		}

	}
}

func doReaddirplus(req Req, nodeid uint64, readOut *kernel.FuseReadOut) int32 {

	readIn := (*req.Arg).(kernel.FuseReadIn)
	se := req.session
	var res int32 = errno.ENOSYS

	if se.Debug {
		log.Trace.Printf("Readdirplus: %+v \n", readIn)
	}

	if se.Opts != nil && se.Opts.Readdirplus != nil {

		fi := NewFuseFileInfo()
		fi.Fh = readIn.Fh

		var content []byte

		content, res = (*se.Opts.Readdirplus)(req, nodeid, readIn.Size, readIn.Offset, fi)

		if res == errno.SUCCESS {
			readOut.Content = content
		}

	}

	return res
}

func doInterrupt(req Req) {
	se := req.session

	if se.Debug {
		log.Trace.Println("Interrupt")
	}

	if se.Opts != nil && se.Opts.Readdirplus != nil {
		(*se.Opts.Interrupt)(req)
	}
}

const offsetMax = 0x7fffffffffffffff

func convertFuseFileLock(lock kernel.FuseFileLock, flock *Flock) {
	flock.Type = int16(lock.Type)
	flock.Whence = int16(os.SEEK_SET)
	flock.Start = int64(lock.Start)
	if lock.End == offsetMax {
		flock.Len = 0
	} else {
		flock.Len = int64(lock.End - lock.Start + 1)
	}
	flock.Pid = int32(lock.Pid)
}

func setFuseLkOut(flock Flock, getlkOut *kernel.FuseLkOut) {
	getlkOut.Lk.Type = uint32(flock.Type)
	if flock.Type != syscall.F_UNLCK {
		getlkOut.Lk.Start = uint64(flock.Start)
		if flock.Len == 0 {
			getlkOut.Lk.End = offsetMax
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
	attr.UID = stat.Uid
	attr.GID = stat.Gid
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
	stat.Uid = setattrIn.UID
	stat.Gid = setattrIn.Gid

}

func setOpenOut(openOut *kernel.FuseOpenOut, fi FileInfo) {
	openOut.Fh = fi.Fh

	if fi.DirectIo > 0 {
		openOut.OpenFlags |= FOpenDirectIO
	}

	if fi.KeepCache > 0 {
		openOut.OpenFlags |= FOpenKeepCache
	}

	if fi.Nonseekable > 0 {
		openOut.OpenFlags |= FOpenNonSeekable
	}
}
