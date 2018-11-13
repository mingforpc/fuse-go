package fuse

import (
	"bytes"
	"fmt"

	"github.com/mingforpc/fuse-go/fuse/common"
	"github.com/mingforpc/fuse-go/fuse/errno"
	"github.com/mingforpc/fuse-go/fuse/kernel"
	"github.com/mingforpc/fuse-go/fuse/log"
)

// The loop to read/write '/dev/fuse'
func FuseLoop(se *FuseSession) {

	respChannel := make(chan []byte, 1024)

	go func() {
		for true {

			res := <-respChannel
			err := writeCmd(se, res)

			if err != nil {
				log.Error.Println(err)
			}
		}

	}()

	for true {

		inheader, buf, err := readCmd(se)

		if err != nil {

			fmt.Println(err)
			break
		}

		req := FuseReq{}
		req.Init(se, *inheader)

		go func() {
			res, err := distribute(&req, inheader, buf)

			if err == kernel.NoNeedReplyErr {
				// This request no need to reply
				return
			}

			if err != nil {
				log.Error.Println(err)
			} else {
				respChannel <- res
			}
		}()

	}

}

// Read event from '/dev/fuse'
func readCmd(se *FuseSession) (*kernel.FuseInHeader, []byte, error) {
	var cmdLenBytes = make([]byte, se.Bufsize)

	n, err := se.Dev.Read(cmdLenBytes)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	cmdLenBytes = cmdLenBytes[0:n]

	var inheader = kernel.FuseInHeader{}

	headerbytes := cmdLenBytes[:40]
	opsbytes := cmdLenBytes[40:]

	err = inheader.ParseBinary(headerbytes)

	if se.Debug {
		log.Trace.Printf("cmdLenBytes[%+v]", cmdLenBytes)
		log.Trace.Printf("inheader: %+v, content[%+v]", inheader, opsbytes)
	}

	return &inheader, opsbytes, err
}

// Write response to '/dev/fuse'
func writeCmd(se *FuseSession, resp []byte) error {
	if se.Debug {
		log.Trace.Printf("resp[%+v] \n", resp)
	}
	_, err := se.Dev.Write(resp)

	return err
}

// Distribute event to earch function
func distribute(req *FuseReq, inHeader *kernel.FuseInHeader, bcontent []byte) ([]byte, error) {

	var arg interface{}
	var errnum int32
	var resp kernel.FuseRespIntf

	var noreply = false

	switch inHeader.Opcode {
	case kernel.FUSE_INIT:
		// Init event
		var initIn = kernel.FuseInitIn{}
		// common.ParseBinary(bcontent, &initIn)
		initIn.ParseBinary(bcontent)
		arg = initIn
		req.Arg = &arg
		var initOut = kernel.FuseInitOut{}

		errnum = doInit(*req, inHeader.Nodeid, &initOut)

		resp = initOut

	case kernel.FUSE_FORGET:
		// Forget event
		var fotgetIn = kernel.FuseForgetIn{}
		// common.ParseBinary(bcontent, &fotgetIn)
		fotgetIn.ParseBinary(bcontent)
		arg = fotgetIn
		req.Arg = &arg

		doForget(*req, inHeader.Nodeid)

		noreply = true

	case kernel.FUSE_LOOKUP:
		// lookup event
		var lookupIn = kernel.FuseLookupIn{}
		// common.ParseBinary(bcontent, &lookupIn)
		lookupIn.ParseBinary(bcontent)
		arg = lookupIn
		req.Arg = &arg

		var entryOut = kernel.FuseEntryOut{}

		errnum = doLookup(*req, inHeader.Nodeid, &entryOut)

		resp = entryOut

	case kernel.FUSE_GETATTR:
		// Getattr event
		var getattrIn = kernel.FuseGetattrIn{}
		// common.ParseBinary(bcontent, &getattrIn)
		getattrIn.ParseBinary(bcontent)
		arg = getattrIn
		req.Arg = &arg
		var getattrOut = kernel.FuseAttrOut{}

		errnum = doGetattr(*req, inHeader.Nodeid, &getattrOut)

		resp = getattrOut

	case kernel.FUSE_SETATTR:
		// Setattr event
		var setattrIn = kernel.FuseSetattrIn{}
		// common.ParseBinary(bcontent, &setattrIn)
		setattrIn.ParseBinary(bcontent)
		arg = setattrIn
		req.Arg = &arg
		var setattrOut = kernel.FuseAttrOut{}

		errnum = doSetattr(*req, inHeader.Nodeid, &setattrOut)

		resp = setattrOut

	case kernel.FUSE_READLINK:
		// Readlink event
		var readlinkOut = kernel.FuseReadlinkOut{}

		errnum = doReadlink(*req, inHeader.Nodeid, &readlinkOut)

		resp = readlinkOut
	case kernel.FUSE_MKNOD:
		// Mknod event
		var mknodIn = kernel.FuseMknodIn{}
		// common.ParseBinary(bcontent, &mknodIn)
		mknodIn.ParseBinary(bcontent)
		arg = mknodIn
		req.Arg = &arg

		var mknodOut = kernel.FuseEntryOut{}

		errnum = doMknod(*req, inHeader.Nodeid, &mknodOut)

		resp = mknodOut

	case kernel.FUSE_MKDIR:
		// Mkdir event
		var mkdirIn = kernel.FuseMkdirIn{}
		// common.ParseBinary(bcontent, &mkdirIn)
		mkdirIn.ParseBinary(bcontent)
		arg = mkdirIn
		req.Arg = &arg

		var mkdirOut = kernel.FuseEntryOut{}

		errnum = doMkdir(*req, inHeader.Nodeid, &mkdirOut)

		resp = mkdirOut

	case kernel.FUSE_UNLINK:
		// Unlink event
		var unlinkIn = kernel.FuseUnlinkIn{}
		// common.ParseBinary(bcontent, &unlinkIn)
		unlinkIn.ParseBinary(bcontent)
		arg = unlinkIn
		req.Arg = &arg

		errnum = doUnlink(*req, inHeader.Nodeid)

	case kernel.FUSE_RMDIR:
		// Rmdir event
		var rmdirIn = kernel.FuseRmdirIn{}
		// common.ParseBinary(bcontent, &rmdirIn)
		rmdirIn.ParseBinary(bcontent)
		arg = rmdirIn
		req.Arg = &arg

		errnum = doRmdir(*req, inHeader.Nodeid)
	case kernel.FUSE_SYMLINK:
		// Symlink event
		var symlinkIn = kernel.FuseSymlinkIn{}
		// common.ParseBinary(bcontent, &symlinkIn)
		symlinkIn.ParseBinary(bcontent)
		arg = symlinkIn
		req.Arg = &arg

		var symlinkOut = kernel.FuseEntryOut{}

		errnum = doSymlink(*req, inHeader.Nodeid, &symlinkOut)

		resp = symlinkOut

	case kernel.FUSE_RENAME:
		// Rename event
		var renameIn = kernel.FuseRenameIn{}
		// common.ParseBinary(bcontent, &renameIn)
		renameIn.ParseBinary(bcontent)
		arg = renameIn
		req.Arg = &arg

		errnum = doRename(*req, inHeader.Nodeid)

	case kernel.FUSE_RENAME2:
		// Rename2 event
		var renameIn = kernel.FuseRename2In{}
		// common.ParseBinary(bcontent, &renameIn)
		renameIn.ParseBinary(bcontent)
		arg = renameIn
		req.Arg = &arg

		errnum = doRename2(*req, inHeader.Nodeid)
	case kernel.FUSE_LINK:
		// Link event
		var linkIn = kernel.FuseLinkIn{}
		// common.ParseBinary(bcontent, &linkIn)
		linkIn.ParseBinary(bcontent)
		arg = linkIn
		req.Arg = &arg

		var linkOut = kernel.FuseEntryOut{}

		errnum = doLink(*req, inHeader.Nodeid, &linkOut)

		resp = linkOut

	case kernel.FUSE_OPEN:
		// Open event
		var openIn = kernel.FuseOpenIn{}
		// common.ParseBinary(bcontent, &openIn)
		openIn.ParseBinary(bcontent)
		arg = openIn
		req.Arg = &arg

		var openOut = kernel.FuseOpenOut{}

		errnum = doOpen(*req, inHeader.Nodeid, &openOut)

		resp = openOut
	case kernel.FUSE_READ:
		// Read event
		var readIn = kernel.FuseReadIn{}
		// common.ParseBinary(bcontent, &readIn)
		readIn.ParseBinary(bcontent)
		arg = readIn
		req.Arg = &arg

		var readOut = kernel.FuseReadOut{}

		errnum = doRead(*req, inHeader.Nodeid, &readOut)

		resp = readOut

	case kernel.FUSE_WRITE:
		// Write event
		var writeIn = kernel.FuseWriteIn{}
		// common.ParseBinary(bcontent, &writeIn)
		writeIn.ParseBinary(bcontent)
		arg = writeIn
		req.Arg = &arg

		var writeOut = kernel.FuseWriteOut{}

		errnum = doWrite(*req, inHeader.Nodeid, &writeOut)

		resp = writeOut

	case kernel.FUSE_FSYNC:
		// Fsync event
		var fsyncIn = kernel.FuseFsyncIn{}
		// common.ParseBinary(bcontent, &fsyncIn)
		fsyncIn.ParseBinary(bcontent)
		arg = fsyncIn
		req.Arg = &arg

		errnum = doFsync(*req, inHeader.Nodeid)

	case kernel.FUSE_OPENDIR:
		// Opendir event
		var openIn = kernel.FuseOpenIn{}
		// common.ParseBinary(bcontent, &openIn)
		openIn.ParseBinary(bcontent)
		arg = openIn
		req.Arg = &arg

		var openOut = kernel.FuseOpenOut{}

		errnum = doOpendir(*req, inHeader.Nodeid, &openOut)

		resp = openOut

	case kernel.FUSE_READDIR:
		// Readdir event
		var readIn = kernel.FuseReadIn{}
		// common.ParseBinary(bcontent, &readIn)
		readIn.ParseBinary(bcontent)
		arg = readIn
		req.Arg = &arg

		var readOut = kernel.FuseReadOut{}

		errnum = doReaddir(*req, inHeader.Nodeid, &readOut)

		resp = readOut

	case kernel.FUSE_RELEASEDIR:
		// Releasedir event
		var releasedirIn = kernel.FuseReleaseIn{}
		// common.ParseBinary(bcontent, &releasedirIn)
		releasedirIn.ParseBinary(bcontent)
		arg = releasedirIn
		req.Arg = &arg

		errnum = doReleasedir(*req, inHeader.Nodeid)

	case kernel.FUSE_FSYNCDIR:
		// Fsyncdir event
		var fsyncdirIn = kernel.FuseFsyncIn{}
		// common.ParseBinary(bcontent, &fsyncdirIn)
		fsyncdirIn.ParseBinary(bcontent)
		arg = fsyncdirIn
		req.Arg = &arg

		errnum = doFsyncdir(*req, inHeader.Nodeid)

	case kernel.FUSE_STATFS:
		// Statfs event

		var statfsOut = kernel.FuseStatfsOut{}

		errnum = doStatfs(*req, inHeader.Nodeid, &statfsOut)

		resp = statfsOut

	case kernel.FUSE_SETXATTR:
		// Setxattr event

		var setxattrIn = kernel.FuseSetxattrIn{}
		// common.ParseBinary(bcontent, &setxattrIn)
		setxattrIn.ParseBinary(bcontent)
		arg = setxattrIn
		req.Arg = &arg

		errnum = doSetxattr(*req, inHeader.Nodeid)

	case kernel.FUSE_GETXATTR:
		// Getxattr event

		var getxattrIn = kernel.FuseGetxattrIn{}
		// common.ParseBinary(bcontent, &getxattrIn)
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

	case kernel.FUSE_LISTXATTR:
		// Listxattr event

		var listxattrIn = kernel.FuseGetxattrIn{}
		// common.ParseBinary(bcontent, &listxattrIn)
		listxattrIn.ParseBinary(bcontent)
		arg = listxattrIn
		req.Arg = &arg

		var listxattrOut = kernel.FuseGetxattrOut{}
		errnum = doGetxattr(*req, inHeader.Nodeid, &listxattrOut)

		if listxattrOut.Value == "" {
			resp = listxattrOut
		} else {
			resp = listxattrOut.Value
		}

	case kernel.FUSE_REMOVEXATTR:
		// Removexattr event

		var removexattrIn = kernel.FuseRemovexattrIn{}
		// common.ParseBinary(bcontent, &removexattrIn)
		removexattrIn.ParseBinary(bcontent)
		arg = removexattrIn
		req.Arg = &arg

		errnum = doRemovexattr(*req, inHeader.Nodeid)

	case kernel.FUSE_ACCESS:
		// Access event

		var accessIn = kernel.FuseAccessIn{}
		// common.ParseBinary(bcontent, &accessIn)
		accessIn.ParseBinary(bcontent)
		arg = accessIn
		req.Arg = &arg

		errnum = doAccess(*req, inHeader.Nodeid)

	case kernel.FUSE_CREATE:
		// Create event

		var createIn = kernel.FuseCreateIn{}
		// common.ParseBinary(bcontent, &createIn)
		createIn.ParseBinary(bcontent)
		arg = createIn
		req.Arg = &arg

		var createOut = kernel.FuseEntryOut{}

		errnum = doCreate(*req, inHeader.Nodeid, &createOut)

		resp = createOut

	case kernel.FUSE_GETLK:
		// Getlk event

		var getlkIn = kernel.FuseLkIn{}
		// common.ParseBinary(bcontent, &getlkIn)
		getlkIn.ParseBinary(bcontent)
		arg = getlkIn
		req.Arg = &arg

		var lkout = kernel.FuseLkOut{}

		errnum = doGetlk(*req, inHeader.Nodeid, &lkout)

		resp = lkout

	case kernel.FUSE_SETLK:
		// Setlk event

		var setlkIn = kernel.FuseLkIn{}
		// common.ParseBinary(bcontent, &setlkIn)
		setlkIn.ParseBinary(bcontent)
		arg = setlkIn
		req.Arg = &arg

		errnum = doSetlk(*req, inHeader.Nodeid)

	case kernel.FUSE_SETLKW:
		// Getlkw event

		var getlkIn = kernel.FuseLkIn{}
		// common.ParseBinary(bcontent, &getlkIn)
		getlkIn.ParseBinary(bcontent)
		arg = getlkIn
		req.Arg = &arg

		errnum = doSetlkw(*req, inHeader.Nodeid)
	case kernel.FUSE_BMAP:
		// Bmap event

		var bmapIn = kernel.FuseBmapIn{}
		// common.ParseBinary(bcontent, &bmapIn)
		bmapIn.ParseBinary(bcontent)
		arg = bmapIn
		req.Arg = &arg

		var bmapOut = kernel.FuseBmapOut{}

		errnum = doBmap(*req, inHeader.Nodeid, &bmapOut)

		resp = bmapOut

	case kernel.FUSE_IOCTL:
		// Ioctl event

		var ioctlIn = kernel.FuseIoctlIn{}
		// common.ParseBinary(bcontent, ioctlIn)
		ioctlIn.ParseBinary(bcontent)
		arg = ioctlIn
		req.Arg = &arg

		var ioctlOut = kernel.FuseIoctlOut{}

		errnum = doIoctl(*req, inHeader.Nodeid, &ioctlOut)

		resp = ioctlOut

	case kernel.FUSE_POLL:
		// Poll event

		var pollIn = kernel.FusePollIn{}
		// common.ParseBinary(bcontent, pollIn)
		pollIn.ParseBinary(bcontent)
		arg = pollIn
		req.Arg = &arg

		var pollOut = kernel.FusePollOut{}

		errnum = doPoll(*req, inHeader.Nodeid, &pollOut)

		resp = pollOut

	case kernel.FUSE_FALLOCATE:
		// Fallocate event

		var fallocateIn = kernel.FuseFallocateIn{}
		// common.ParseBinary(bcontent, fallocateIn)
		fallocateIn.ParseBinary(bcontent)
		arg = fallocateIn
		req.Arg = &arg

		errnum = doFallocate(*req, inHeader.Nodeid)

	case kernel.FUSE_BATCH_FORGET:
		// Batch forget event

		var batchForgetIn = kernel.FuseBatchForgetIn{}
		// common.ParseBinary(bcontent, batchForgetIn)
		batchForgetIn.ParseBinary(bcontent)
		arg = batchForgetIn
		req.Arg = &arg

		doForgetMulti(*req)

		noreply = true

	case kernel.FUSE_READDIRPLUS:
		// readdirplus event

		var readIn = kernel.FuseReadIn{}
		// common.ParseBinary(bcontent, &readIn)
		readIn.ParseBinary(bcontent)
		arg = readIn
		req.Arg = &arg

		var readOut = kernel.FuseReadOut{}

		errnum = doReaddirplus(*req, inHeader.Nodeid, &readOut)

		resp = readOut
	}

	var bresp []byte
	var err error

	if noreply {
		// means this request no need to reply
		bresp = nil
		err = kernel.NoNeedReplyErr
	} else {

		outHeader := kernel.FuseOutHeader{}
		outHeader.Error = errnum
		outHeader.Unique = inHeader.Unique

		if req.session.Debug {
			log.Trace.Printf("outHeader[%+v], resp[%+v]", outHeader, resp)
		}

		if errnum == errno.SUCCESS {
			bresp, err = generateResp(outHeader, resp)
		} else {
			bresp, err = generateResp(outHeader, nil)
		}

	}

	return bresp, err
}

// Parse opcode from bytes
func parseOpcode(bdata []byte) (uint32, error) {
	var opcode uint32

	err := common.ParseBinary(bdata, &opcode)

	return opcode, err
}

// Function to generate bytes response
func generateResp(outHeader kernel.FuseOutHeader, resp kernel.FuseRespIntf) ([]byte, error) {

	buf := bytes.NewBuffer(nil)

	var bresp []byte
	var err error

	if resp != nil {
		bresp, err = resp.ToBinary()
		if err != nil {
			return nil, err
		}
	}

	outHeader.Len = uint32(kernel.OUT_HEADER_LEN + len(bresp))

	bheader, err := outHeader.ToBinary()
	if err != nil {
		return nil, err
	}

	buf.Write(bheader)
	buf.Write(bresp)

	return buf.Bytes(), nil
}
