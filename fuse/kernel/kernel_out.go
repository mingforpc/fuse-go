package kernel

import (
	"bytes"
	"encoding/binary"

	"github.com/mingforpc/fuse-go/fuse/common"
)

const OUT_HEADER_LEN = 16

type FuseRespIntf interface {
	ToBinary() ([]byte, error)
}

// each answer starts with a FuseOutHeader
// 16 bytes
type FuseOutHeader struct {
	Len    uint32
	Error  int32
	Unique uint64
}

func (outHeader FuseOutHeader) ToBinary() ([]byte, error) {
	return common.ToBinary(outHeader)
}

// init
// 64 bytes
type FuseInitOut struct {
	// Header FuseOutHeader

	Major               uint32
	Minor               uint32
	MaxReadahead        uint32
	Flags               uint32
	MaxBackground       uint16
	CongestionThreshold uint16
	MaxWrite            uint32
	TimeGran            uint32
	Unused              [9]uint32
}

func (init FuseInitOut) ToBinary() ([]byte, error) {
	return common.ToBinary(init)
}

// getattr, setattr
type FuseAttrOut struct {
	AttrValid     uint64 /* Cache timeout for the attributes */
	AttrValidNsec uint32
	Dummp         uint32

	Attr FuseAttr
}

func (attr FuseAttrOut) ToBinary() ([]byte, error) {
	return common.ToBinary(attr)
}

// readlink
type FuseReadlinkOut struct {
	Path string
}

func (link FuseReadlinkOut) ToBinary() ([]byte, error) {

	buf := bytes.NewBuffer(nil)

	buf.WriteString(link.Path)

	buf.WriteByte(0)

	return buf.Bytes(), nil
}

// lookup, symlink, mknod, mkdir, link
type FuseEntryOut struct {
	NodeId         uint64 /* Inode ID */
	Generation     uint64 /* Inode generation: nodeid:gen must be unique for the fs's lifetime */
	EntryValid     uint64 /* Cache timeout for the name */
	AttrValid      uint64 /* Cache timeout for the attributes */
	EntryValidNsec uint32
	AttrValidNsec  uint32

	Attr FuseAttr
}

func (entry FuseEntryOut) ToBinary() ([]byte, error) {

	return common.ToBinary(entry)
}

// create
// create 要返回FuseEntryOut和FuseOpenOut
type FuseCreateOut struct {
	Entry FuseEntryOut
	Open  FuseOpenOut
}

func (create FuseCreateOut) ToBinary() ([]byte, error) {

	return common.ToBinary(create)
}

// open, opendir
type FuseOpenOut struct {
	Fh        uint64
	OpenFlags uint32
	Padding   uint32
}

func (open FuseOpenOut) ToBinary() ([]byte, error) {

	return common.ToBinary(open)
}

// read, readdir
type FuseReadOut struct {
	Content []byte
}

func (read FuseReadOut) ToBinary() ([]byte, error) {

	return read.Content, nil
}

// 目录的结构体，二进制方式写入readdir的Content中
type FuseDirent struct {
	Ino     uint64
	Off     uint64
	NameLen uint32
	DirType uint32
	Name    string
}

const DIRENT_NAME_OFFSET = 24

func fuseDirentAlign(entlent uint64) uint64 {

	return (entlent + 8 - 1) & (^uint64(7))
}

// ToBinary将FuseDirent转为二进制数据, preOff是在list里面，上一个的偏移量
func (dirent *FuseDirent) ToBinary(preOff *uint64) ([]byte, error) {

	entLen := fuseDirentAlign(DIRENT_NAME_OFFSET + uint64(dirent.NameLen))

	buf := bytes.NewBuffer(nil)

	*preOff += entLen

	binary.Write(buf, binary.LittleEndian, dirent.Ino)
	binary.Write(buf, binary.LittleEndian, *preOff)
	binary.Write(buf, binary.LittleEndian, dirent.NameLen)
	binary.Write(buf, binary.LittleEndian, dirent.DirType)

	buf.WriteString(dirent.Name)

	buf.Write(make([]byte, int(entLen)-buf.Len()))

	return buf.Bytes(), nil
}

// write
type FuseWriteOut struct {
	Size    uint32
	Padding uint32
}

func (write FuseWriteOut) ToBinary() ([]byte, error) {

	return common.ToBinary(write)
}

// statfs
type FuseStatfsOut struct {
	St FuseStatfs
}

func (statfs FuseStatfsOut) ToBinary() ([]byte, error) {

	return common.ToBinary(statfs)
}

type XattrVal string

func (val XattrVal) ToBinary() ([]byte, error) {

	buf := bytes.NewBuffer(nil)

	buf.WriteString(string(val))
	buf.WriteByte(0)

	return buf.Bytes(), nil
}

// getxattr, listxattr
type FuseGetxattrOut struct {
	Size    uint32
	Padding uint32

	Value XattrVal
}

func (getxattr FuseGetxattrOut) ToBinary() ([]byte, error) {

	buf := bytes.NewBuffer(nil)

	binary.Write(buf, binary.LittleEndian, getxattr.Size)
	binary.Write(buf, binary.LittleEndian, getxattr.Padding)

	if len(getxattr.Value) > 0 {
		valb, _ := getxattr.Value.ToBinary()
		buf.Write(valb)
		buf.WriteByte(0)
	}

	return buf.Bytes(), nil
}

// getlk, setlk, setlkw
type FuseLkOut struct {
	Lk FuseFileLock
}

func (lk FuseLkOut) ToBinary() ([]byte, error) {

	return common.ToBinary(lk)
}

// ioctl
type FuseIoctlOut struct {
	Result  int32
	Flags   uint32
	InIovs  uint32
	OutIovs uint32
}

func (ioctl FuseIoctlOut) ToBinary() ([]byte, error) {

	return common.ToBinary(ioctl)
}

// poll
type FusePollOut struct {
	Revents uint32
	Padding uint32
}

func (poll FusePollOut) ToBinary() ([]byte, error) {

	return common.ToBinary(poll)
}

// lseek
type FuseLseekOut struct {
	Offset uint64
}

func (lseek FuseLseekOut) ToBinary() ([]byte, error) {

	return common.ToBinary(lseek)
}

// bmap
type FuseBmapOut struct {
	Block uint64
}

func (bmap FuseBmapOut) ToBinary() ([]byte, error) {

	return common.ToBinary(bmap)
}

// cuse_init
type CuseInitOut struct {
	// Header FuseOutHeader

	Major    uint32
	Minor    uint32
	Unused   uint32
	Flags    uint32
	MaxRead  uint32
	MaxWrite uint32
	DevMajor uint32
	DevMinor uint32
	Spare    [10]uint32
}

func (cuseinit CuseInitOut) ToBinary() ([]byte, error) {

	return common.ToBinary(cuseinit)
}
