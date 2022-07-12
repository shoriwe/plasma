package common

import (
	"encoding/binary"
	"math"
)

type (
	integer interface {
		int | int8 | int16 | int32 | int64
	}
	float interface {
		float32 | float64
	}
)

const (
	Zero byte = iota
	Bits8
	Bits16
	Bits32
	Bits64
)

func IntToBytes[T integer](i T) []byte {
	ui64 := uint64(i)
	var bytes []byte
	switch {
	case i == 0:
		return []byte{Zero}
	case ui64 < uint64(^uint8(0)):
		bytes = []byte{Bits8, byte(i)}
	case ui64 < uint64(^uint16(0)):
		bytes = []byte{Bits16, 0, 0}
		binary.BigEndian.PutUint16(bytes[1:], uint16(i))
	case ui64 < uint64(^uint32(0)):
		bytes = []byte{Bits32, 0, 0, 0, 0}
		binary.BigEndian.PutUint32(bytes[1:], uint32(i))
	default:
		bytes = []byte{Bits64, 0, 0, 0, 0, 0, 0, 0, 0}
		binary.BigEndian.PutUint64(bytes[1:], uint64(i))
	}
	return bytes
}

func FloatToBytes[T float](f T) []byte {
	var bytes [8]byte
	u := math.Float64bits(float64(f))
	binary.BigEndian.PutUint64(bytes[:], u)
	return bytes[:]
}
