package messages

import (
	"GoDofus/utils"
	"bytes"
)

type protocol struct {
	packetId int32
	Version  string
}

var proto = &protocol{packetId: 9546}

func GetProtocolNOA() *protocol {
	return proto
}

func (p *protocol) Serialize(buff *bytes.Buffer) {
	utils.WriteUTF(buff, p.Version)
}

func (p *protocol) Deserialize(reader *bytes.Reader) {
	p.Version = utils.ReadUTF(reader)
}
