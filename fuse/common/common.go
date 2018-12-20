package common

import (
	"bytes"
	"encoding/binary"
)

const (
	// Uint64Max max value of uint64
	Uint64Max = ^uint64(0)

	// Uint32Max max vaalue of uint32
	Uint32Max = ^uint32(0)
)

// ParseBinary parser binary from source to target
func ParseBinary(source []byte, target interface{}) error {

	err := binary.Read(bytes.NewBuffer(source), binary.LittleEndian, target)
	if err != nil {
		return err
	}

	return nil
}

// ToBinary parser source to binary
func ToBinary(source interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	err := binary.Write(buf, binary.LittleEndian, source)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}

func CalcTimeoutSec(t float64) uint64 {
	if t > float64(Uint64Max) {
		return Uint64Max
	} else if t < 0.0 {
		return 0
	} else {
		return uint64(t)
	}
}

func CalcTimeoutNsec(t float64) uint32 {

	f := t - float64(CalcTimeoutSec(t))

	if f < 0.0 {
		return 0
	} else if f >= 0.999999999 {
		return 999999999
	} else {
		return uint32(f * 1.0e9)
	}
}
