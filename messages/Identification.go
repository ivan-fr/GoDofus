package messages

import (
	"GoDofus/managers"
	"GoDofus/structs"
	"GoDofus/utils"
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"fmt"
)

type Identification struct {
	PacketId            uint32
	Version             *version
	Lang                []byte
	Credentials         []byte
	ServerId            uint16
	AutoSelectServer    bool
	UseCertificate      bool
	UseLoginToken       bool
	SessionOptionalSalt float64
	FailedAttempts      []uint32
	instance            uint
}

var publicVerifyPem = utils.ReadRSA("./binaryData/verify_key.bin")
var blockVerify = utils.DecodePem(publicVerifyPem)
var publicKeyVerify = utils.PublicKeyOf(blockVerify)

var identificationMap = make(map[uint]*Identification)

func (id *Identification) GetNOA(instance uint) Message {
	identification_, ok := identificationMap[instance]

	if ok {
		return identification_
	}

	identificationMap[instance] = &Identification{PacketId: IdentificationID, Version: new(version),
		UseCertificate: false, UseLoginToken: false, ServerId: 0, instance: instance}
	return identificationMap[instance]
}

func (id *Identification) Serialize(buff *bytes.Buffer) {
	var box uint32
	box = utils.SetFlag(box, 0, id.AutoSelectServer)
	box = utils.SetFlag(box, 1, id.UseCertificate)
	box = utils.SetFlag(box, 2, id.UseLoginToken)

	_ = binary.Write(buff, binary.BigEndian, byte(box))
	id.Version.Serialize(buff)
	utils.WriteUTF(buff, id.Lang)
	utils.WriteVarInt32(buff, int32(len(id.Credentials)))
	_ = binary.Write(buff, binary.BigEndian, id.Credentials)
	_ = binary.Write(buff, binary.BigEndian, id.ServerId)

	if id.SessionOptionalSalt < -9007199254740992 || id.SessionOptionalSalt > 9007199254740992 {
		panic("Forbidden value on element SessionOptionalSalt.")
	}

	utils.WriteVarInt64(buff, id.SessionOptionalSalt)

	_ = binary.Write(buff, binary.BigEndian, uint16(len(id.FailedAttempts)))

	for i := 0; i < len(id.FailedAttempts); i++ {
		utils.WriteVarInt16(buff, int32(id.FailedAttempts[i]))
	}
}

func (id *Identification) Deserialize(reader *bytes.Reader) {
	var box byte
	_ = binary.Read(reader, binary.BigEndian, &box)
	id.AutoSelectServer = utils.GetFlag(uint32(box), 0)
	id.UseCertificate = utils.GetFlag(uint32(box), 1)
	id.UseLoginToken = utils.GetFlag(uint32(box), 2)
	id.Version.Deserialize(reader)
	id.Lang = utils.ReadUTF(reader)

	credLen := utils.ReadVarInt32(reader)
	id.Credentials = make([]byte, credLen)
	_ = binary.Read(reader, binary.BigEndian, id.Credentials)

	_ = binary.Read(reader, binary.BigEndian, &id.ServerId)

	id.SessionOptionalSalt = float64(utils.ReadVarInt64(reader))

	var failsLen uint16
	_ = binary.Read(reader, binary.BigEndian, &failsLen)

	id.FailedAttempts = nil
	for i := uint16(0); i < failsLen; i++ {
		id.FailedAttempts = append(id.FailedAttempts, uint32(utils.ReadVarInt16(reader)))
	}
}

func (id *Identification) GetPacketId() uint32 {
	return id.PacketId
}

func (id *Identification) InitIdentificationMessage() {
	idManager := managers.GetAuthentificationManager(id.instance)
	idManager.InitLoginAction()

	id.Lang = idManager.Lang
	id.AutoSelectServer = idManager.LA.AutoSelectServer

	currentVersion := structs.GetVersionNOA()
	id.Version.Major = currentVersion.Major
	id.Version.Minor = currentVersion.Minor
	id.Version.Code = currentVersion.Code
	id.Version.Build = currentVersion.Build
	id.Version.BuildType = currentVersion.BuildType

	id.Credentials = idManager.GetCipher(id.getPublicKey(), id.getSalt())
}

func (id *Identification) getPublicKey() *rsa.PublicKey {
	idManager := managers.GetAuthentificationManager(id.instance)

	hc := Types_[HelloConnectID].GetNOA(id.instance).(*HelloConnect)

	if idManager.PublicKey != nil && bytes.Compare(hc.Salt, idManager.Salt) == 0 {
		return idManager.PublicKey
	}

	if hc.Key == nil {
		panic("helloMessage wasn't call")
	}

	publicKey := utils.RsaPublicDecrypt(publicKeyVerify, hc.Key)

	pki := fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----",
		base64.StdEncoding.EncodeToString(publicKey))

	publicKeyPem, _ := pem.Decode([]byte(pki))

	pkiX, err := x509.ParsePKIXPublicKey(publicKeyPem.Bytes)
	if err != nil {
		panic(err)
	}

	idManager.PublicKey = pkiX.(*rsa.PublicKey)
	return idManager.PublicKey
}

func (id *Identification) getSalt() []byte {
	idManager := managers.GetAuthentificationManager(id.instance)

	hc := Types_[HelloConnectID].GetNOA(id.instance).(*HelloConnect)
	if hc.Salt == nil {
		panic("helloMessage wasn't call")
	}

	mySalt := hc.Salt
	idManager.Salt = hc.Salt

	if len(mySalt) < 32 {
		for len(mySalt) < 32 {
			mySalt = append(mySalt, byte(' '))
		}
	}

	return mySalt
}

func (id *Identification) String() string {
	return fmt.Sprintf("PacketId: %d\nVersion: %v\nidentification: ...\n\t%s\n\t%v\n\t%v\n\t%v\n\t%v\n\t%v\n",
		id.PacketId, id.Version, id.Lang, id.Credentials, id.SessionOptionalSalt, id.FailedAttempts, id.ServerId, len(id.Credentials))
}
