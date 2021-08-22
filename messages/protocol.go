package messages

import (
	"GoDofus/utils"
	"bytes"
)

type protocol struct {
	packetId int32
	version  string
}

var proto = &protocol{packetId: 9546}

func GetProtocol(version string) *protocol {
	proto.version = version
	return proto
}

func GetProtocolNOA() *protocol {
	return proto
}

func (p *protocol) Serialize(buff *bytes.Buffer) {
	utils.WriteUTF(buff, p.version)
}

func (p *protocol) Deserialize(reader *bytes.Reader) {
	p.version = utils.ReadUTF(reader)
}
