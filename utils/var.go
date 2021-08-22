package utils

import (
	"bytes"
	"encoding/binary"
	"math"
)

var mask10000000 int32 = 128
var mask01111111 int32 = 127
var chunkBitSize int32 = 7
var int32Size uint8 = 32
var int16Size uint8 = 16

func ReadVarInt16(reader *bytes.Reader) int32 {
	var aByte uint8
	var value int32
	var offset uint8
	var hasNext bool

	for offset < int16Size {
		_ = binary.Read(reader, binary.BigEndian, aByte)
		hasNext = (aByte & uint8(mask10000000)) == uint8(mask10000000)

		if offset > 0 {
			value += int32((aByte & uint8(mask01111111)) << offset)
		} else {
			value += int32(aByte & uint8(mask01111111))
		}

		offset += uint8(chunkBitSize)

		if !hasNext {
			if value > int32(math.MaxInt16) {
				value -= int32(math.MaxUint16 + 1)
			}

			return value
		}
	}

	panic("too much data")
}

func ReadVarInt32(reader *bytes.Reader) int32 {
	var aByte uint8
	var value int32
	var offset uint8
	var hasNext bool

	for offset < int32Size {
		_ = binary.Read(reader, binary.BigEndian, aByte)
		hasNext = (aByte & uint8(mask10000000)) == uint8(mask10000000)

		if offset > 0 {
			value += int32((aByte & uint8(mask01111111)) << offset)
		} else {
			value += int32(aByte & uint8(mask01111111))
		}

		offset += uint8(chunkBitSize)

		if !hasNext {
			return value
		}
	}

	panic("too much data")
}

func WriteVarInt32(buff *bytes.Buffer, value int32) {
	var aByte uint8

	if value >= 0 && value <= mask01111111 {
		_ = binary.Write(buff, binary.BigEndian, uint8(value))
		return
	}

	var c = value

	for c != 0 {
		aByte = uint8(c & mask01111111)
		c >>= chunkBitSize

		if c > 0 {
			aByte |= uint8(mask10000000)
		}
		_ = binary.Write(buff, binary.BigEndian, aByte)
	}
}

func WriteVarShort(buff *bytes.Buffer, value int32) {
	var aByte uint8

	if value > math.MaxInt16 || value < math.MinInt16 {
		panic("forbidden value")
	}

	if value >= 0 && value <= mask01111111 {
		_ = binary.Write(buff, binary.BigEndian, uint8(value))
		return
	}

	var c = value & 65535
	for c != 0 {
		aByte = uint8(c & mask01111111)
		c >>= chunkBitSize

		if c > 0 {
			aByte |= uint8(mask10000000)
		}
		_ = binary.Write(buff, binary.BigEndian, aByte)
	}
}

func WriteVarLong(buff *bytes.Buffer, value float64) {
	var val = int64(value)

	var low = int32(val & 0xffffffff)
	var high = int32(val >> 32)
	if high == 0 {
		_ = binary.Write(buff, binary.BigEndian, low)
		return
	}

	for i := 0; i < 4; i++ {
		_ = binary.Write(buff, binary.BigEndian, uint8(low&mask01111111|mask10000000))
		low >>= chunkBitSize
	}

	if (high & 268435455 << 3) == 0 {
		_ = binary.Write(buff, binary.BigEndian, uint8(high<<4|low))
	} else {
		_ = binary.Write(buff, binary.BigEndian, uint8((high<<4|low)&mask01111111|mask10000000))
		_ = binary.Write(buff, binary.BigEndian, high>>3)
	}
}
