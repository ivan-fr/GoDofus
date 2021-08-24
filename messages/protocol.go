package messages

import (
	"GoDofus/utils"
	"bytes"
	"fmt"
)

type protocol struct {
	packetId uint32
	Version  []byte
}

var proto = &protocol{packetId: ProtocolID}

func GetProtocolNOA() *protocol {
	return proto
}

func (p *protocol) Serialize(buff *bytes.Buffer) {
	utils.WriteUTF(buff, p.Version)
}

func (p *protocol) Deserialize(reader *bytes.Reader) {
	p.Version = utils.ReadUTF(reader)
}

func (p *protocol) String() string {
	return fmt.Sprintf("packetId: %d\nVersion: %s\n", p.packetId, p.Version)
}
