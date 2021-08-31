package pack

import (
	"GoDofus/messages"
	"bytes"
	"encoding/binary"
	"fmt"
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

	return 0
}

func WriteWeft(weft *Weft, toClient bool, instance uint) []byte {

	packetId := weft.PackId

	defer func() {
		if toClient {
			fmt.Printf("Instance n°%d: %d -----> myClient.\n", instance, packetId)
		} else {
			fmt.Printf("Instance n°%d: %d -----> offcial.\n", instance, packetId)
		}
	}()

	buff := new(bytes.Buffer)

	typeLength := weft.LengthType

	twoBytesHeader := packetId<<2 | typeLength
	_ = binary.Write(buff, binary.BigEndian, twoBytesHeader)

	if !toClient {
		instanceId++
		_ = binary.Write(buff, binary.BigEndian, instanceId)
	}

	switch typeLength {
	case 1:
		var lenMessage = uint8(weft.Length)
		_ = binary.Write(buff, binary.BigEndian, lenMessage)
	case 2:
		var lenMessage = uint16(weft.Length)
		_ = binary.Write(buff, binary.BigEndian, lenMessage)
	case 3:
		var high = uint8(weft.Length >> 16 & uint32(math.MaxUint8))
		var low = uint16(weft.Length & uint32(math.MaxUint16))
		_ = binary.Write(buff, binary.BigEndian, high)
		_ = binary.Write(buff, binary.BigEndian, low)
	case 0:
		return buff.Bytes()
	default:
		panic("wrong typeLength")
	}

	_ = binary.Write(buff, binary.BigEndian, weft.Message)

	return buff.Bytes()
}

func Write(message messages.Message, toClient bool, instance uint) []byte {
	buffMsg := new(bytes.Buffer)
	message.Serialize(buffMsg)

	messageContent := buffMsg.Bytes()

	packetId := uint16(message.GetPacketId())

	defer func() {
		if toClient {
			fmt.Printf("Instance n°%d: %d -----> myClient.\n", instance, packetId)
		} else {
			fmt.Printf("Instance n°%d: %d -----> offcial.\n", instance, packetId)
		}
	}()

	buff := new(bytes.Buffer)

	typeLength := computeTypeLength(uint32(len(messageContent)))

	twoBytesHeader := packetId<<2 | typeLength
	_ = binary.Write(buff, binary.BigEndian, twoBytesHeader)

	if !toClient {
		instanceId++
		_ = binary.Write(buff, binary.BigEndian, instanceId)
	}

	switch typeLength {
	case 1:
		var lenMessage = uint8(len(messageContent))
		_ = binary.Write(buff, binary.BigEndian, lenMessage)
	case 2:
		var lenMessage = uint16(len(messageContent))
		_ = binary.Write(buff, binary.BigEndian, lenMessage)
	case 3:
		var high = uint8(uint32(len(messageContent)) >> 16 & uint32(math.MaxUint8))
		var low = uint16(uint32(len(messageContent)) & uint32(math.MaxUint16))
		_ = binary.Write(buff, binary.BigEndian, high)
		_ = binary.Write(buff, binary.BigEndian, low)
	case 0:
		return buff.Bytes()
	default:
		panic("wrong typeLength")
	}

	_ = binary.Write(buff, binary.BigEndian, messageContent)

	return buff.Bytes()
}
