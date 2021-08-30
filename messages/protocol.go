package messages

import (
	"GoDofus/utils"
	"bytes"
	"fmt"
)

type protocol struct {
	PacketId uint32
	Version  []byte
}

var protocolMap = make(map[uint]*protocol)

func (p *protocol) GetNOA(instance uint) Message {
	protocol_, ok := protocolMap[instance]

	if ok {
		return protocol_
	}

	protocolMap[instance] = &protocol{PacketId: ProtocolID}
	return protocolMap[instance]
}

func (p *protocol) Serialize(buff *bytes.Buffer) {
	utils.WriteUTF(buff, p.Version)
}

func (p *protocol) Deserialize(reader *bytes.Reader) {
	p.Version = utils.ReadUTF(reader)
}

func (p *protocol) GetPacketId() uint32 {
	return p.PacketId
}

func (p *protocol) String() string {
	return fmt.Sprintf("PacketId: %d\nVersion: %s\n", p.PacketId, p.Version)
}
