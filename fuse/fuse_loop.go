package fuse

import (
	"bytes"
	"errors"
	"syscall"

	"github.com/mingforpc/fuse-go/fuse/evloop"

	"github.com/mingforpc/fuse-go/fuse/errno"
	"github.com/mingforpc/fuse-go/fuse/kernel"
	"github.com/mingforpc/fuse-go/fuse/log"
)

// FuseLoop : the loop to read/write '/dev/fuse'
func (se *Session) FuseLoop() {

	if !se.IsInited() {
		panic(kernel.ErrNotInit)
	}

	se.Running = true

	se.readChan = make(chan []byte, se.maxGoro)
	se.writeChan = make(chan []byte, se.maxGoro)
	se.closeCh = make(chan interface{})

	// Write goroutine
	// 用来写"/dev/fuse"的goroutine
	go se.writeGoro()

	// Read goroutine
	// 用来读取"/dev/fuse"的goroutine
	go se.readGoro()

	for se.Running {

		brep, ok := <-se.readChan

		if !ok {
			break
		}

		inheader, buf, err := se.parseHeader(brep)

		if err != nil {

			log.Error.Printf("Session parseHeader[%s] \n", err)
			break
		}

		req := Req{}
		req.Init(se, inheader)

		// 用来处理各个请求的goroutine
		go func() {

			defer func() {
				if err := recover(); err != nil {
					log.Error.Printf("Distribute goroutine error[%s] \n", err)
				}

			}()

			res, err := distribute(&req, inheader, buf)

			if err == kernel.ErrNoNeedReply {
				// This request no need to reply
				return
			}

			if err != nil {
				log.Error.Println(err)
			} else {
				// Avoid srive closed channel
				select {
				case <-se.closeCh:
					return
				case se.writeChan <- res:
				}

			}
		}()

	}

}

func (se *Session) readGoro() {
	defer func() {
		if err := recover(); err != nil {
			log.Error.Printf("Read goroutine error[%s] \n", err)
		}
	}()

	el := se.evloop

	handler := func(el *evloop.EvLoop, fd int, eventmask int, privdata interface{}) {
		breq, err := se.readCmd()
		if err != nil {

			if err == syscall.ENODEV {

				se.Running = false
			} else {
				// 读出错退出
				log.Error.Printf("err: %+v \n", err)
				panic(err)
			}
		}

		// Read 可能block很久，所以再判断一次
		if se.Running {
			se.readChan <- breq
		}
	}

	err := el.Register(se.devFd, evloop.EPOLLIN, handler, nil)

	if err != nil {
		panic(err)
	}

	for se.Running {
		// wait 1 second
		el.Process(1000)
	}

	close(se.readChan)

}

func (se *Session) writeGoro() {
	for se.Running {

		res, ok := <-se.writeChan

		if !ok {
			break
		}

		if se.Running {
			err := se.writeCmd(res)
			if err != nil {
				log.Error.Println(err)
			}
		}

	}

}

// Close : close fuse session
func (se *Session) Close() {
	se.Running = false

	syscall.Close(se.devFd)

	close(se.closeCh)
	close(se.writeChan)

}

// Read event from '/dev/fuse'
func (se *Session) readCmd() ([]byte, error) {
	var cmdLenBytes = make([]byte, se.bufsize)

	n, err := syscall.Read(se.devFd, cmdLenBytes)

	if err != nil {
		return nil, err
	}

	cmdLenBytes = cmdLenBytes[0:n]

	return cmdLenBytes, err
}

func (se *Session) parseHeader(bcontent []byte) (kernel.FuseInHeader, []byte, error) {
	var inheader = kernel.FuseInHeader{}

	headerbytes := bcontent[:40]
	opsbytes := bcontent[40:]

	err := inheader.ParseBinary(headerbytes)

	if se.Debug {
		log.Trace.Printf("cmdLenBytes[%+v]", bcontent)
		log.Trace.Printf("inheader: %+v, content[%+v]", inheader, opsbytes)
	}

	return inheader, opsbytes, err
}

// Write response to '/dev/fuse'
func (se *Session) writeCmd(resp []byte) error {
	if se.Debug {
		log.Trace.Printf("resp[%+v] \n", resp)
	}
	_, err := syscall.Write(se.devFd, resp)

	return err
}

// Distribute event to earch function
func distribute(req *Req, inHeader kernel.FuseInHeader, bcontent []byte) ([]byte, error) {

	var arg interface{}
	var errnum int32
	var resp kernel.FuseResponsor

	var noreply = false

	switch inHeader.Opcode {
	case kernel.FuseOpInit:
		// Init event
		var initIn = kernel.FuseInitIn{}
		initIn.ParseBinary(bcontent)
		arg = initIn
		req.Arg = &arg
		var initOut = kernel.FuseInitOut{}

		errnum = doInit(*req, &initOut)

		resp = initOut

	case kernel.FuseOpDestory:
		// Destory event

		doDestory(*req)

		noreply = true

	case kernel.FuseOpForget:
		// Forget event
		var fotgetIn = kernel.FuseForgetIn{}
		fotgetIn.ParseBinary(bcontent)
		arg = fotgetIn
		req.Arg = &arg

		doForget(*req, inHeader.Nodeid)

		noreply = true

	case kernel.FuseOpLookup:
		// lookup event
		var lookupIn = kernel.FuseLookupIn{}
		lookupIn.ParseBinary(bcontent)
		arg = lookupIn
		req.Arg = &arg

		var entryOut = kernel.FuseEntryOut{}

		errnum = doLookup(*req, inHeader.Nodeid, &entryOut)

		resp = entryOut

	case kernel.FuseOpGetattr:
		// Getattr event
		var getattrIn = kernel.FuseGetattrIn{}
		getattrIn.ParseBinary(bcontent)
		arg = getattrIn
		req.Arg = &arg
		var getattrOut = kernel.FuseAttrOut{}

		errnum = doGetattr(*req, inHeader.Nodeid, &getattrOut)

		resp = getattrOut

	case kernel.FuseOpSetattr:
		// Setattr event
		var setattrIn = kernel.FuseSetattrIn{}
		setattrIn.ParseBinary(bcontent)
		arg = setattrIn
		req.Arg = &arg
		var setattrOut = kernel.FuseAttrOut{}

		errnum = doSetattr(*req, inHeader.Nodeid, &setattrOut)

		resp = setattrOut

	case kernel.FuseOpReadlink:
		// Readlink event
		var readlinkOut = kernel.FuseReadlinkOut{}

		errnum = doReadlink(*req, inHeader.Nodeid, &readlinkOut)

		resp = readlinkOut
	case kernel.FuseOpMknod:
		// Mknod event
		var mknodIn = kernel.FuseMknodIn{}
		mknodIn.ParseBinary(bcontent)
		arg = mknodIn
		req.Arg = &arg

		var mknodOut = kernel.FuseEntryOut{}

		errnum = doMknod(*req, inHeader.Nodeid, &mknodOut)

		resp = mknodOut

	case kernel.FuseOpMkdir:
		// Mkdir event
		var mkdirIn = kernel.FuseMkdirIn{}
		mkdirIn.ParseBinary(bcontent)
		arg = mkdirIn
		req.Arg = &arg

		var mkdirOut = kernel.FuseEntryOut{}

		errnum = doMkdir(*req, inHeader.Nodeid, &mkdirOut)

		resp = mkdirOut

	case kernel.FuseOpUnlink:
		// Unlink event
		var unlinkIn = kernel.FuseUnlinkIn{}
		unlinkIn.ParseBinary(bcontent)
		arg = unlinkIn
		req.Arg = &arg

		errnum = doUnlink(*req, inHeader.Nodeid)

	case kernel.FuseOpRmdir:
		// Rmdir event
		var rmdirIn = kernel.FuseRmdirIn{}
		rmdirIn.ParseBinary(bcontent)
		arg = rmdirIn
		req.Arg = &arg

		errnum = doRmdir(*req, inHeader.Nodeid)
	case kernel.FuseOpSymlink:
		// Symlink event
		var symlinkIn = kernel.FuseSymlinkIn{}
		symlinkIn.ParseBinary(bcontent)
		arg = symlinkIn
		req.Arg = &arg

		var symlinkOut = kernel.FuseEntryOut{}

		errnum = doSymlink(*req, inHeader.Nodeid, &symlinkOut)

		resp = symlinkOut

	case kernel.FuseOpRename:
		// Rename event
		var renameIn = kernel.FuseRenameIn{}
		renameIn.ParseBinary(bcontent)
		arg = renameIn
		req.Arg = &arg

		errnum = doRename(*req, inHeader.Nodeid)

	case kernel.FuseOpRename2:
		// Rename2 event
		var renameIn = kernel.FuseRename2In{}
		renameIn.ParseBinary(bcontent)
		arg = renameIn
		req.Arg = &arg

		errnum = doRename2(*req, inHeader.Nodeid)
	case kernel.FuseOpLink:
		// Link event
		var linkIn = kernel.FuseLinkIn{}
		linkIn.ParseBinary(bcontent)
		arg = linkIn
		req.Arg = &arg

		var linkOut = kernel.FuseEntryOut{}

		errnum = doLink(*req, inHeader.Nodeid, &linkOut)

		resp = linkOut

	case kernel.FuseOpOpen:
		// Open event
		var openIn = kernel.FuseOpenIn{}
		openIn.ParseBinary(bcontent)
		arg = openIn
		req.Arg = &arg

		var openOut = kernel.FuseOpenOut{}

		errnum = doOpen(*req, inHeader.Nodeid, &openOut)

		resp = openOut
	case kernel.FuseOpRead:
		// Read event
		var readIn = kernel.FuseReadIn{}
		readIn.ParseBinary(bcontent)
		arg = readIn
		req.Arg = &arg

		var readOut = kernel.FuseReadOut{}

		errnum = doRead(*req, inHeader.Nodeid, &readOut)

		resp = readOut

	case kernel.FuseOpWrite:
		// Write event
		var writeIn = kernel.FuseWriteIn{}
		writeIn.ParseBinary(bcontent)
		arg = writeIn
		req.Arg = &arg

		var writeOut = kernel.FuseWriteOut{}

		errnum = doWrite(*req, inHeader.Nodeid, &writeOut)

		resp = writeOut

	case kernel.FuseOpFsync:
		// Fsync event
		var fsyncIn = kernel.FuseFsyncIn{}
		fsyncIn.ParseBinary(bcontent)
		arg = fsyncIn
		req.Arg = &arg

		errnum = doFsync(*req, inHeader.Nodeid)

	case kernel.FuseOpOpendir:
		// Opendir event
		var openIn = kernel.FuseOpenIn{}
		openIn.ParseBinary(bcontent)
		arg = openIn
		req.Arg = &arg

		var openOut = kernel.FuseOpenOut{}

		errnum = doOpendir(*req, inHeader.Nodeid, &openOut)

		resp = openOut

	case kernel.FuseOpReaddir:
		// Readdir event
		var readIn = kernel.FuseReadIn{}
		readIn.ParseBinary(bcontent)
		arg = readIn
		req.Arg = &arg

		var readOut = kernel.FuseReadOut{}

		errnum = doReaddir(*req, inHeader.Nodeid, &readOut)

		resp = readOut

	case kernel.FuseOpRelease:
		// Release event
		var releaseIn = kernel.FuseReleaseIn{}
		releaseIn.ParseBinary(bcontent)
		arg = releaseIn
		req.Arg = &arg

		errnum = doRelease(*req, inHeader.Nodeid)

	case kernel.FuseOpReleasedir:
		// Releasedir event
		var releasedirIn = kernel.FuseReleaseIn{}
		releasedirIn.ParseBinary(bcontent)
		arg = releasedirIn
		req.Arg = &arg

		errnum = doReleasedir(*req, inHeader.Nodeid)

	case kernel.FuseOpFlush:
		var flushIn = kernel.FuseFlushIn{}
		flushIn.ParseBinary(bcontent)
		arg = flushIn
		req.Arg = &arg

		errnum = doFlush(*req, inHeader.Nodeid)

	case kernel.FuseOpFsyncdir:
		// Fsyncdir event
		var fsyncdirIn = kernel.FuseFsyncIn{}
		fsyncdirIn.ParseBinary(bcontent)
		arg = fsyncdirIn
		req.Arg = &arg

		errnum = doFsyncdir(*req, inHeader.Nodeid)

	case kernel.FuseOpStatfs:
		// Statfs event

		var statfsOut = kernel.FuseStatfsOut{}

		errnum = doStatfs(*req, inHeader.Nodeid, &statfsOut)

		resp = statfsOut

	case kernel.FuseOpSetxattr:
		// Setxattr event

		var setxattrIn = kernel.FuseSetxattrIn{}
		setxattrIn.ParseBinary(bcontent)
		arg = setxattrIn
		req.Arg = &arg

		errnum = doSetxattr(*req, inHeader.Nodeid)

	case kernel.FuseOpGetxattr:
		// Getxattr event

		var getxattrIn = kernel.FuseGetxattrIn{}
		getxattrIn.ParseBinary(bcontent)
		arg = getxattrIn
		req.Arg = &arg

		var getxattrOut = kernel.FuseGetxattrOut{}
		errnum = doGetxattr(*req, inHeader.Nodeid, &getxattrOut)

		if getxattrOut.Value == "" {
			resp = getxattrOut
		} else {
			resp = getxattrOut.Value
		}

	case kernel.FuseOpListxattr:
		// Listxattr event

		var listxattrIn = kernel.FuseGetxattrIn{}
		listxattrIn.ParseBinary(bcontent)
		arg = listxattrIn
		req.Arg = &arg

		var listxattrOut = kernel.FuseGetxattrOut{}
		errnum = doListxattr(*req, inHeader.Nodeid, &listxattrOut)

		if listxattrOut.Value == "" {
			resp = listxattrOut
		} else {
			resp = listxattrOut.Value
		}

	case kernel.FuseOpRemovexattr:
		// Removexattr event

		var removexattrIn = kernel.FuseRemovexattrIn{}
		removexattrIn.ParseBinary(bcontent)
		arg = removexattrIn
		req.Arg = &arg

		errnum = doRemovexattr(*req, inHeader.Nodeid)

	case kernel.FuseOpAccess:
		// Access event

		var accessIn = kernel.FuseAccessIn{}
		accessIn.ParseBinary(bcontent)
		arg = accessIn
		req.Arg = &arg

		errnum = doAccess(*req, inHeader.Nodeid)

	case kernel.FuseOpCreate:
		// Create event

		var createIn = kernel.FuseCreateIn{}
		createIn.ParseBinary(bcontent)
		arg = createIn
		req.Arg = &arg

		var createOut = kernel.FuseCreateOut{}

		errnum = doCreate(*req, inHeader.Nodeid, &createOut)

		if errnum == errno.ENOENT {
			noreply = true
		} else {
			resp = createOut
		}

	case kernel.FuseOpGetlk:
		// Getlk event

		var getlkIn = kernel.FuseLkIn{}
		getlkIn.ParseBinary(bcontent)
		arg = getlkIn
		req.Arg = &arg

		var lkout = kernel.FuseLkOut{}

		errnum = doGetlk(*req, inHeader.Nodeid, &lkout)

		resp = lkout

	case kernel.FuseOpSetlk:
		// Setlk event

		var setlkIn = kernel.FuseLkIn{}
		setlkIn.ParseBinary(bcontent)
		arg = setlkIn
		req.Arg = &arg

		errnum = doSetlk(*req, inHeader.Nodeid)

	case kernel.FuseOpSetlkw:
		// Getlkw event

		var getlkIn = kernel.FuseLkIn{}
		getlkIn.ParseBinary(bcontent)
		arg = getlkIn
		req.Arg = &arg

		errnum = doSetlkw(*req, inHeader.Nodeid)
	case kernel.FuseOpBmap:
		// Bmap event

		var bmapIn = kernel.FuseBmapIn{}
		bmapIn.ParseBinary(bcontent)
		arg = bmapIn
		req.Arg = &arg

		var bmapOut = kernel.FuseBmapOut{}

		errnum = doBmap(*req, inHeader.Nodeid, &bmapOut)

		resp = bmapOut

	case kernel.FuseOpIoctl:
		// Ioctl event

		var ioctlIn = kernel.FuseIoctlIn{}
		ioctlIn.ParseBinary(bcontent)
		arg = ioctlIn
		req.Arg = &arg

		var ioctlOut = kernel.FuseIoctlOut{}

		errnum = doIoctl(*req, inHeader.Nodeid, &ioctlOut)

		resp = ioctlOut

	case kernel.FuseOpPoll:
		// Poll event

		var pollIn = kernel.FusePollIn{}
		pollIn.ParseBinary(bcontent)
		arg = pollIn
		req.Arg = &arg

		var pollOut = kernel.FusePollOut{}

		errnum = doPoll(*req, inHeader.Nodeid, &pollOut)

		resp = pollOut

	case kernel.FuseOpFallocate:
		// Fallocate event

		var fallocateIn = kernel.FuseFallocateIn{}
		fallocateIn.ParseBinary(bcontent)
		arg = fallocateIn
		req.Arg = &arg

		errnum = doFallocate(*req, inHeader.Nodeid)

	case kernel.FuseOpBatckForget:
		// Batch forget event

		var batchForgetIn = kernel.FuseBatchForgetIn{}
		batchForgetIn.ParseBinary(bcontent)
		arg = batchForgetIn
		req.Arg = &arg

		doForgetMulti(*req)

		noreply = true

	case kernel.FuseOpReaddirplus:
		// readdirplus event

		var readIn = kernel.FuseReadIn{}
		readIn.ParseBinary(bcontent)
		arg = readIn
		req.Arg = &arg

		var readOut = kernel.FuseReadOut{}

		errnum = doReaddirplus(*req, inHeader.Nodeid, &readOut)

		resp = readOut

	case kernel.FuseOpInterrupt:
		// interrupt event

		doInterrupt(*req)

	default:
		panic(errors.New("未实现的操作！！！"))
	}

	var bresp []byte
	var err error

	if noreply {
		// means this request no need to reply
		bresp = nil
		err = kernel.ErrNoNeedReply
	} else {

		outHeader := kernel.FuseOutHeader{}
		outHeader.Error = errnum
		outHeader.Unique = inHeader.Unique

		if errnum == errno.SUCCESS {

			if req.session.Debug {
				log.Trace.Printf("errnum[%d], outHeader[%+v], resp[%+v]", errnum, outHeader, resp)
			}

			bresp, err = generateResp(outHeader, resp)
		} else {

			if req.session.Debug {
				log.Trace.Printf("errnum[%d], outHeader[%+v]", errnum, outHeader)
			}

			bresp, err = generateResp(outHeader, nil)
		}

	}

	return bresp, err
}

// Function to generate bytes response
func generateResp(outHeader kernel.FuseOutHeader, resp kernel.FuseResponsor) ([]byte, error) {

	buf := bytes.NewBuffer(nil)

	var bresp []byte
	var err error

	if resp != nil {
		bresp, err = resp.ToBinary()
		if err != nil {
			return nil, err
		}
	}

	outHeader.Len = uint32(kernel.OutHeaderLen + len(bresp))

	bheader, err := outHeader.ToBinary()
	if err != nil {
		return nil, err
	}

	buf.Write(bheader)
	buf.Write(bresp)

	return buf.Bytes(), nil
}
