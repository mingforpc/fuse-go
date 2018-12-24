package kernel

import (
	"bytes"
	"encoding/binary"

	"github.com/mingforpc/fuse-go/fuse/common"
)

// OutHeaderLen : lenght of FuseOutHeader
const OutHeaderLen = 16

// FuseResponsor : the interface of fuse response
type FuseResponsor interface {
	ToBinary() ([]byte, error)
}

// FuseOutHeader : the header of response,
// each answer starts with this.
// 16 bytes
type FuseOutHeader struct {
	Len    uint32
	Error  int32
	Unique uint64
}

// ToBinary : Parse to binary
func (outHeader FuseOutHeader) ToBinary() ([]byte, error) {
	return common.ToBinary(outHeader)
}

// FuseInitOut : init response
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

// ToBinary : Parse to binary
func (init FuseInitOut) ToBinary() ([]byte, error) {
	return common.ToBinary(init)
}

// FuseAttrOut : getattr, setattr response
type FuseAttrOut struct {
	AttrValid     uint64 /* Cache timeout for the attributes */
	AttrValidNsec uint32
	Dummp         uint32

	Attr FuseAttr
}

// ToBinary : Parse to binary
func (attr FuseAttrOut) ToBinary() ([]byte, error) {
	return common.ToBinary(attr)
}

// FuseReadlinkOut : readlink response
type FuseReadlinkOut struct {
	Path string
}

// ToBinary : Parse to binary
func (link FuseReadlinkOut) ToBinary() ([]byte, error) {

	buf := bytes.NewBuffer(nil)

	buf.WriteString(link.Path)

	buf.WriteByte(0)

	return buf.Bytes(), nil
}

// FuseEntryOut : lookup, symlink, mknod, mkdir, link response
type FuseEntryOut struct {
	NodeID         uint64 /* Inode ID */
	Generation     uint64 /* Inode generation: nodeid:gen must be unique for the fs's lifetime */
	EntryValid     uint64 /* Cache timeout for the name */
	AttrValid      uint64 /* Cache timeout for the attributes */
	EntryValidNsec uint32
	AttrValidNsec  uint32

	Attr FuseAttr
}

// ToBinary : Parse to binary
func (entry FuseEntryOut) ToBinary() ([]byte, error) {

	return common.ToBinary(entry)
}

// FuseCreateOut : create response
// create 要返回FuseEntryOut和FuseOpenOut
type FuseCreateOut struct {
	Entry FuseEntryOut
	Open  FuseOpenOut
}

// ToBinary : Parse to binary
func (create FuseCreateOut) ToBinary() ([]byte, error) {

	return common.ToBinary(create)
}

// FuseOpenOut : open, opendir response
type FuseOpenOut struct {
	Fh        uint64
	OpenFlags uint32
	Padding   uint32
}

// ToBinary : Parse to binary
func (open FuseOpenOut) ToBinary() ([]byte, error) {

	return common.ToBinary(open)
}

// FuseReadOut : read, readdir response
type FuseReadOut struct {
	Content []byte
}

// ToBinary : Parse to binary
func (read FuseReadOut) ToBinary() ([]byte, error) {

	return read.Content, nil
}

// FuseDirent : 目录的结构体，二进制方式写入readdir的Content中
type FuseDirent struct {
	Ino     uint64
	Off     uint64
	NameLen uint32
	DirType uint32
	Name    string
}

// direntNameOffset : the lenght of direntNameOffset
const direntNameOffset = 24

// fuseDirentAlign : 用来保证长度是8的n次方
func fuseDirentAlign(entlent uint64) uint64 {

	return (entlent + 8 - 1) & (^uint64(7))
}

// ToBinary : Parse to binary
// ToBinary将FuseDirent转为二进制数据, preOff是在list里面，上一个的偏移量
func (dirent *FuseDirent) ToBinary(preOff *uint64) ([]byte, error) {

	entLen := fuseDirentAlign(direntNameOffset + uint64(dirent.NameLen))

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

// FuseWriteOut : write response
type FuseWriteOut struct {
	Size    uint32
	Padding uint32
}

// ToBinary : Parse to binary
func (write FuseWriteOut) ToBinary() ([]byte, error) {

	return common.ToBinary(write)
}

// FuseStatfsOut : statfs response
type FuseStatfsOut struct {
	St FuseStatfs
}

// ToBinary : Parse to binary
func (statfs FuseStatfsOut) ToBinary() ([]byte, error) {

	return common.ToBinary(statfs)
}

// XattrVal : value of xattr
type XattrVal string

// ToBinary : Parse to binary
func (val XattrVal) ToBinary() ([]byte, error) {

	buf := bytes.NewBuffer(nil)

	buf.WriteString(string(val))
	buf.WriteByte(0)

	return buf.Bytes(), nil
}

// FuseGetxattrOut : getxattr, listxattr response
type FuseGetxattrOut struct {
	Size    uint32
	Padding uint32

	Value XattrVal
}

// ToBinary : Parse to binary
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

// FuseLkOut : getlk, setlk, setlkw response
type FuseLkOut struct {
	Lk FuseFileLock
}

// ToBinary : Parse to binary
func (lk FuseLkOut) ToBinary() ([]byte, error) {

	return common.ToBinary(lk)
}

// FuseIoctlOut : ioctl response
type FuseIoctlOut struct {
	Result  int32
	Flags   uint32
	InIovs  uint32
	OutIovs uint32
}

// ToBinary : Parse to binary
func (ioctl FuseIoctlOut) ToBinary() ([]byte, error) {

	return common.ToBinary(ioctl)
}

// FusePollOut : poll response
type FusePollOut struct {
	Revents uint32
	Padding uint32
}

// ToBinary : Parse to binary
func (poll FusePollOut) ToBinary() ([]byte, error) {

	return common.ToBinary(poll)
}

// FuseLseekOut : lseek response
type FuseLseekOut struct {
	Offset uint64
}

// ToBinary : Parse to binary
func (lseek FuseLseekOut) ToBinary() ([]byte, error) {

	return common.ToBinary(lseek)
}

// FuseBmapOut : bmap response
type FuseBmapOut struct {
	Block uint64
}

// ToBinary : Parse to binary
func (bmap FuseBmapOut) ToBinary() ([]byte, error) {

	return common.ToBinary(bmap)
}

// CuseInitOut : cuse_init response
type CuseInitOut struct {
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

// ToBinary : Parse to binary
func (cuseinit CuseInitOut) ToBinary() ([]byte, error) {

	return common.ToBinary(cuseinit)
}
