package messages

import (
	"GoDofus/utils"
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

type helloConnect struct {
	PacketId uint32
	Salt     []byte
	Key      []byte
}

var privateKey *rsa.PrivateKey
var publicHelloKey []byte

var hConnectMap = make(map[uint]*helloConnect)

func GetHelloConnectNOA(instance uint) *helloConnect {
	hConnect_, ok := hConnectMap[instance]

	if ok {
		return hConnect_
	}

	hConnectMap[instance] = &helloConnect{PacketId: HelloConnectID}
	return hConnectMap[instance]
}

func (h *helloConnect) Serialize(buff *bytes.Buffer) {
	utils.WriteUTF(buff, h.Salt)

	if privateKey == nil {
		privatePam, _ := os.ReadFile("./sign/private_key.pem")
		block, _ := pem.Decode(privatePam)
		if block == nil {
			log.Fatal("failed to decode PEM block containing private key")
		}

		privateKey, _ = x509.ParsePKCS1PrivateKey(block.Bytes)

		pubHelloPam, _ := os.ReadFile("./sign/hello_public_key.pem")
		block, _ = pem.Decode(pubHelloPam)
		if block == nil {
			log.Fatal("failed to decode PEM block containing pub hello key")
		}

		publicHelloKey = block.Bytes
	}

	key, _ := rsa.SignPKCS1v15(nil, privateKey, crypto.Hash(0), publicHelloKey)

	utils.WriteVarInt32(buff, int32(len(key)))
	_ = binary.Write(buff, binary.BigEndian, key)
}

func (h *helloConnect) Deserialize(reader *bytes.Reader) {
	h.Salt = utils.ReadUTF(reader)
	keyLen := uint(utils.ReadVarInt32(reader))
	h.Key = make([]byte, keyLen)
	_ = binary.Read(reader, binary.BigEndian, &h.Key)
}

func (h *helloConnect) GetPacketId() uint32 {
	return h.PacketId
}

func (h *helloConnect) String() string {
	return fmt.Sprintf("PacketId: %d\nSalt: %s\nKey: %v\nlength Key: %d\n", h.PacketId, h.Salt, h.Key, len(h.Key))
}
