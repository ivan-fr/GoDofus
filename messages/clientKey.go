package messages

import (
	"GoDofus/utils"
	"bytes"
	"fmt"
)

type clientKey struct {
	PacketId uint32
	key      []byte
}

var cK = &clientKey{PacketId: ClientKeyID}

func GetClientKeyNOA() *clientKey {
	return cK
}

func (ck *clientKey) Serialize(buff *bytes.Buffer) {
	utils.WriteUTF(buff, ck.key)
}

func (ck *clientKey) Deserialize(reader *bytes.Reader) {
	ck.key = utils.ReadUTF(reader)
}

func (ck *clientKey) GetPacketId() uint32 {
	return ck.PacketId
}

func (ck *clientKey) String() string {
	return fmt.Sprintf("PacketId: %d\nkey: %s\n", ck.PacketId, ck.key)
}
