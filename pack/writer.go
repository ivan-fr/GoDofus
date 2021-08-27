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

func Write(message messages.Message, toClient bool) []byte {
	buffMsg := new(bytes.Buffer)
	message.Serialize(buffMsg)

	messageContent := buffMsg.Bytes()

	packetId := uint16(message.GetPacketId())

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
	}

	_ = binary.Write(buff, binary.BigEndian, messageContent)

	if toClient {
		fmt.Printf("write to the client %d...\n", packetId)
	} else {
		fmt.Printf("write to the server %d...\n", packetId)
	}

	return buff.Bytes()
}
