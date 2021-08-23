package messages

import (
	"GoDofus/utils"
	"bytes"
	"encoding/binary"
	"fmt"
)

type helloConnect struct {
	PacketId uint32
	Salt     string
	Key      []byte
}

var hConnect = &helloConnect{PacketId: HelloConnectID}

func GetHelloConnectNOA() *helloConnect {
	return hConnect
}

func (h *helloConnect) Serialize(buff *bytes.Buffer) {
	utils.WriteUTF(buff, h.Salt)
	utils.WriteVarInt32(buff, int32(len(h.Key)))
	_ = binary.Write(buff, binary.BigEndian, h.Key)
}

func (h *helloConnect) Deserialize(reader *bytes.Reader) {
	h.Salt = utils.ReadUTF(reader)
	keyLen := uint(utils.ReadVarInt32(reader))
	h.Key = make([]byte, keyLen)
	_ = binary.Read(reader, binary.BigEndian, &h.Key)
}

func (h *helloConnect) String() string {
	return fmt.Sprintf("packetId: %d\nSalt: %s\nKey: %v\nlength Key: %d\n", h.PacketId, h.Salt, h.Key, len(h.Key))
}
