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
		_ = binary.Read(reader, binary.BigEndian, &aByte)
		hasNext = (aByte & uint8(mask10000000)) == uint8(mask10000000)

		if offset > 0 {
			value += int32(aByte&uint8(mask01111111)) << int32(offset)
		} else {
			value += int32(aByte & uint8(mask01111111))
		}

		offset += uint8(chunkBitSize)

		if !hasNext {
			if value > int32(math.MaxInt16) {
				value -= int32(math.MaxUint16) + int32(1)
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
	var hasNext = false

	for offset < int32Size {
		_ = binary.Read(reader, binary.BigEndian, &aByte)
		hasNext = (aByte & uint8(mask10000000)) == uint8(mask10000000)

		if offset > 0 {
			value += int32(aByte&uint8(mask01111111)) << int32(offset)
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

func WriteVarInt16(buff *bytes.Buffer, value int32) {
	var aByte uint8

	if value > math.MaxInt16 || value < math.MinInt16 {
		panic("forbidden value")
	}

	if value >= 0 && value <= mask01111111 {
		_ = binary.Write(buff, binary.BigEndian, uint8(value))
		return
	}

	var c = value & math.MaxUint16
	for c != 0 {
		aByte = uint8(c & mask01111111)
		c >>= chunkBitSize

		if c > 0 {
			aByte |= uint8(mask10000000)
		}
		_ = binary.Write(buff, binary.BigEndian, aByte)
	}
}

func WriteVarInt64(buff *bytes.Buffer, value float64) {
	var val = int64(value)

	var low = int32(val & 0xffffffff)
	var high = int32(val >> 32)
	if high == 0 {
		writeSpecialInt32(buff, low)
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
		writeSpecialInt32(buff, high>>3)
	}
}

func WriteVarUInt64(buff *bytes.Buffer, value float64) {
	WriteVarInt64(buff, value)
}

func ReadVarInt64(reader *bytes.Reader) int64 {
	var b byte
	var low int32
	var high int32
	var i int32
	for i < 28 {
		_ = binary.Read(reader, binary.BigEndian, &b)

		if (int32(b) & mask10000000) == mask10000000 {
			low |= int32(b) & mask01111111 << i
			i += 7
		} else {
			low |= int32(b) << i

			return int64(low)
		}
	}

	_ = binary.Read(reader, binary.BigEndian, &b)

	if (int32(b) & mask10000000) == mask10000000 {
		b &= byte(mask01111111)
		low |= int32(b) & mask01111111 << i
		high = int32(b) & mask01111111 >> 4
		i = 3

		for {
			_ = binary.Read(reader, binary.BigEndian, &b)
			if (int32(b) & mask10000000) == mask10000000 {
				high |= (int32(b) & mask01111111) << i
			} else {
				break
			}

			i += 7
		}
		high |= int32(b) << i
		return (int64(high) << 32) | (int64(low) & 0xffffffff)
	}
	low |= int32(b) << i
	high = int32(b) >> 4
	return (int64(high) << 32) | (int64(low) & 0xffffffff)
}

func ReadVarUInt64(reader *bytes.Reader) uint64 {
	b := ReadVarInt64(reader)
	return uint64(b)
}

func writeSpecialInt32(buff *bytes.Buffer, value int32) {
	for value >= 128 {
		_ = binary.Write(buff, binary.BigEndian, uint8(value&127|128))
		value >>= 7
	}
	_ = binary.Write(buff, binary.BigEndian, uint8(value))
}
