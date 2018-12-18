// Copy from https://github.com/hanwen/go-fuse/blob/master/fuse/mount_linux.go

package mount

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"unsafe"

	"github.com/mingforpc/fuse-go/fuse"
)

func unixgramSocketpair() (l, r *os.File, err error) {
	fd, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_SEQPACKET, 0)
	if err != nil {
		return nil, nil, os.NewSyscallError("socketpair",
			err.(syscall.Errno))
	}
	l = os.NewFile(uintptr(fd[0]), "socketpair-half1")
	r = os.NewFile(uintptr(fd[1]), "socketpair-half2")
	return
}

// Create a FUSE FS on the specified mount point.  The returned
// mount point is always absolute.
func Mount(se *fuse.FuseSession) (err error) {
	local, remote, err := unixgramSocketpair()
	if err != nil {
		return
	}

	defer local.Close()
	defer remote.Close()

	bin, err := fusermountBinary()
	if err != nil {
		return err
	}

	cmd := []string{bin, se.Mountpoint}
	opts := []string{"rw", "nosuid", "nodev"}
	cmd = append(cmd, "-o", strings.Join(opts, ","))

	proc, err := os.StartProcess(bin,
		cmd,
		&os.ProcAttr{
			Env:   []string{"_FUSE_COMMFD=3"},
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr, remote}})

	if err != nil {
		return
	}

	w, err := proc.Wait()
	if err != nil {
		return
	}
	if !w.Success() {
		err = fmt.Errorf("fusermount exited with code %v\n", w.Sys())
		return
	}

	var fd int
	fd, err = getConnection(local)
	if err != nil {
		return err
	}

	se.Dev = os.NewFile(uintptr(fd), "/dev/fuse")

	// golang sets CLOEXEC on file descriptors when they are
	// acquired through normal operations (e.g. open).
	// Buf for fd, we have to set CLOEXEC manually
	syscall.CloseOnExec(fd)

	return err
}

func getConnection(local *os.File) (int, error) {
	var data [4]byte
	control := make([]byte, 4*256)

	// n, oobn, recvflags, from, errno  - todo: error checking.
	_, oobn, _, _,
		err := syscall.Recvmsg(
		int(local.Fd()), data[:], control[:], 0)
	if err != nil {
		return 0, err
	}

	message := *(*syscall.Cmsghdr)(unsafe.Pointer(&control[0]))
	fd := *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&control[0])) + syscall.SizeofCmsghdr))

	if message.Type != 1 {
		return 0, fmt.Errorf("getConnection: recvmsg returned wrong control type: %d", message.Type)
	}
	if oobn <= syscall.SizeofCmsghdr {
		return 0, fmt.Errorf("getConnection: too short control message. Length: %d", oobn)
	}
	if fd < 0 {
		return 0, fmt.Errorf("getConnection: fd < 0: %d", fd)
	}
	return int(fd), nil
}

// lookPathFallback - search binary in PATH and, if that fails,
// in fallbackDir. This is useful if PATH is possible empty.
func lookPathFallback(file string, fallbackDir string) (string, error) {
	binPath, err := exec.LookPath(file)
	if err == nil {
		return binPath, nil
	}

	abs := path.Join(fallbackDir, file)
	return exec.LookPath(abs)
}

// 需要先把"/dev/fuse"文件关闭
func Unmount(mountPoint string) (err error) {
	bin, err := fusermountBinary()
	if err != nil {
		return err
	}
	errBuf := bytes.Buffer{}
	cmd := exec.Command(bin, "-uz", mountPoint)
	cmd.Stderr = &errBuf
	err = cmd.Run()
	if errBuf.Len() > 0 {
		return fmt.Errorf("%s (code %v)\n",
			errBuf.String(), err)
	}
	return err
}

func fusermountBinary() (string, error) {
	return lookPathFallback("fusermount", "/bin")
}

func umountBinary() (string, error) {
	return lookPathFallback("umount", "/bin")
}
