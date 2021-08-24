package messages

import (
	"GoDofus/utils"
	"bytes"
	"fmt"
)

type clientKey struct {
	packetId uint32
	key      string
}

var cK = &clientKey{packetId: ClientKeyID}
var uid = "a19fRCh9EAOvmumjSE"

func GetClientKeyNOA() *clientKey {
	if cK.key == "" {
		cK.key = uid + "#01"
	}
	return cK
}

func (ck *clientKey) Serialize(buff *bytes.Buffer) {
	utils.WriteUTF(buff, ck.key)
}

func (ck *clientKey) Deserialize(reader *bytes.Reader) {
	ck.key = utils.ReadUTF(reader)
}

func (ck *clientKey) String() string {
	return fmt.Sprintf("packetId: %d\nkey: %s\n", ck.packetId, ck.key)
}
