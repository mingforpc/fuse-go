package common

import (
	"bytes"
	"encoding/binary"
)

const UINT64_MAX = ^uint64(0)
const UINT32_MAX = ^uint32(0)

func ParseBinary(source []byte, target interface{}) error {

	err := binary.Read(bytes.NewBuffer(source), binary.LittleEndian, target)
	if err != nil {
		return err
	}

	return nil
}

func ToBinary(source interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	err := binary.Write(buf, binary.LittleEndian, source)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}

func CalcTimeoutSec(t float64) uint64 {
	if t > float64(UINT64_MAX) {
		return UINT64_MAX
	} else if t < 0.0 {
		return 0
	} else {
		return uint64(t)
	}
}

func CalcTimeoutNsec(t float64) uint32 {

	var f float64 = t - float64(CalcTimeoutSec(t))

	if f < 0.0 {
		return 0
	} else if f >= 0.999999999 {
		return 999999999
	} else {
		return uint32(f * 1.0e9)
	}
}
