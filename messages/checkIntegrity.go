// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-25 13:49:18.9644378 +0200 CEST m=+0.003136701

package messages

import (
	"GoDofus/utils"
	"bytes"
	"encoding/binary"
	"fmt"
)

type checkIntegrity struct {
	PacketId uint32
	data     []byte
}

var checkIntegrity_ = &checkIntegrity{PacketId: CheckIntegrityID}

func GetCheckIntegrityNOA() *checkIntegrity {
	reader := bytes.NewReader(GetRawDataNOA().content)
	checkIntegrity_.data = utils.DecryptV(reader)
	return checkIntegrity_
}

func (c *checkIntegrity) Serialize(buff *bytes.Buffer) {
	utils.WriteVarInt32(buff, int32(len(c.data)))
	_ = binary.Write(buff, binary.BigEndian, c.data)
}

func (c *checkIntegrity) Deserialize(reader *bytes.Reader) {

}

func (c *checkIntegrity) GetPacketId() uint32 {
	return c.PacketId
}

func (c *checkIntegrity) String() string {
	return fmt.Sprintf("packetId: %d\n", c.PacketId)
}
