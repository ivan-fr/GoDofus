package messages

import (
	"GoDofus/utils"
	"bytes"
	"encoding/binary"
	"fmt"
)

type helloConnect struct {
	PacketId int32
	Salt     string
	Key      []byte
}

var hConnect = &helloConnect{PacketId: 1030}

func GetHelloConnectNOA() *helloConnect {
	return hConnect
}

func (h *helloConnect) Serialize(buff *bytes.Buffer) {
	utils.WriteUTF(buff, h.Salt)
	utils.WriteVarInt32(buff, int32(len(h.Key)))

	for i := uint(0); i < uint(len(h.Key)); i++ {
		_ = binary.Write(buff, binary.BigEndian, h.Key[i])
	}
}

func (h *helloConnect) Deserialize(reader *bytes.Reader) {
	h.Salt = utils.ReadUTF(reader)
	keyLen := uint(utils.ReadVarInt32(reader))
	h.Key = nil
	for i := uint(0); i < keyLen; i++ {
		var aByte byte
		_ = binary.Read(reader, binary.BigEndian, &aByte)
		h.Key = append(h.Key, aByte)
	}
}

func (h *helloConnect) String() string {
	return fmt.Sprintf("PacketId: %d\nSalt: %s\nKey:%v\n", h.PacketId, h.Salt, h.Key)
}
