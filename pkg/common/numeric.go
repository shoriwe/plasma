package common

import (
	"encoding/binary"
	"math"
)

type (
	integer interface {
		int | int32 | int64
	}
	float interface {
		float32 | float64
	}
)

func IntToBytes[T integer](i T) []byte {
	var bytes [8]byte
	binary.BigEndian.PutUint64(bytes[:], uint64(i))
	return bytes[:]
}

func FloatToBytes[T float](f T) []byte {
	var bytes [8]byte
	u := math.Float64bits(float64(f))
	binary.BigEndian.PutUint64(bytes[:], u)
	return bytes[:]
}
