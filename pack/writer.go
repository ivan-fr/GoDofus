package pack

import (
	"bytes"
	"encoding/binary"
	"math"
)

var instanceId uint32 = 0

func computeTypeLength(messageLength uint32) uint16 {
	if messageLength > math.MaxUint16 {
		return 3
	}
	if messageLength > math.MaxUint8 {
		return 2
	}
	if messageLength > 0 {
		return 1
	}

	panic("invalid message length")
}

func Write(packetId uint16, message []byte) []byte {
	instanceId++
	buff := new(bytes.Buffer)

	typeLength := computeTypeLength(uint32(len(message)))

	twoBytesHeader := packetId<<2 | typeLength
	_ = binary.Write(buff, binary.BigEndian, twoBytesHeader)
	_ = binary.Write(buff, binary.BigEndian, instanceId)

	switch typeLength {
	case 1:
		var lenMessage = uint8(len(message))
		_ = binary.Write(buff, binary.BigEndian, lenMessage)
	case 2:
		var lenMessage = uint16(len(message))
		_ = binary.Write(buff, binary.BigEndian, lenMessage)
	case 3:
		var high = uint8(uint32(len(message)) >> 16 & uint32(math.MaxUint8))
		var low = uint16(uint32(len(message)) & uint32(math.MaxUint16))
		_ = binary.Write(buff, binary.BigEndian, high)
		_ = binary.Write(buff, binary.BigEndian, low)
	}

	_ = binary.Write(buff, binary.BigEndian, message)
	return buff.Bytes()
}
