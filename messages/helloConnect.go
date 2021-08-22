package messages

import (
	"GoDofus/utils"
	"bytes"
	"encoding/binary"
)

type helloConnect struct {
	packetId int32
	salt     string
	key      []byte
}

var hConnect = &helloConnect{packetId: 1030}

func GetHelloConnect(salt string, key []byte) *helloConnect {
	hConnect.salt = salt
	hConnect.key = key
	return hConnect
}

func GetHelloConnectNOA() *helloConnect {
	return hConnect
}

func (h *helloConnect) Serialize(buff *bytes.Buffer) {
	utils.WriteUTF(buff, h.salt)
	utils.WriteVarInt32(buff, int32(len(h.key)))

	for i := uint(0); i < uint(len(h.key)); i++ {
		_ = binary.Write(buff, binary.BigEndian, uint8(h.key[i]))
	}
}

func (h *helloConnect) Deserialize(reader *bytes.Reader) {
	h.salt = utils.ReadUTF(reader)
	keyLen := uint(utils.ReadVarInt32(reader))
	h.key = nil
	for i := uint(0); i < keyLen; i++ {
		var aByte byte
		_ = binary.Read(reader, binary.BigEndian, aByte)
		h.key = append(h.key, aByte)
	}
}
