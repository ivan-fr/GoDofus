package managers

import (
	"GoDofus/messages"
	"GoDofus/settings"
	"GoDofus/structs"
	"GoDofus/utils"
	"bytes"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"math/rand"
	"time"
)

type loginAction struct {
	username         string
	password         string
	autoSelectServer bool
	serverId         uint16
	host             string
}

type Authentification struct {
	AESKey    []byte
	lA        *loginAction
	Lang      []byte
	publicKey *rsa.PublicKey
	salt      []byte
	instance  uint
}

var AESLength = uint(32)

var publicVerifyPem = utils.ReadRSA("./binaryData/verify_key.bin")
var blockVerify = utils.DecodePem(publicVerifyPem)
var publicKeyVerify = utils.PublicKeyOf(blockVerify)

var authentificationMap = make(map[uint]*Authentification)

func GetAuthentificationManager(instance uint) *Authentification {
	authentificationMap_, ok := authentificationMap[instance]

	if ok {
		return authentificationMap_
	}

	authentificationMap[instance] = &Authentification{AESKey: generateAESKey(), Lang: []byte("fr"), instance: instance}
	return authentificationMap[instance]
}

func (a *Authentification) initLoginAction() {
	la := &loginAction{autoSelectServer: true}
	la.username = settings.Settings.Ndc
	la.password = settings.Settings.Pass
	a.lA = la
	var randomTime = time.Duration(rand.Intn(1) + 1)
	time.Sleep(time.Second * randomTime)
}

func (a *Authentification) getCipher() []byte {
	buff := new(bytes.Buffer)
	mySalt := a.getSalt()
	_ = binary.Write(buff, binary.BigEndian, mySalt)
	_ = binary.Write(buff, binary.BigEndian, a.AESKey)
	_ = binary.Write(buff, binary.BigEndian, byte(len(a.lA.username)))
	_ = binary.Write(buff, binary.BigEndian, []byte(a.lA.username))
	_ = binary.Write(buff, binary.BigEndian, []byte(a.lA.password))

	credentials, err := rsa.EncryptPKCS1v15(cryptoRand.Reader, a.getPublicKey(), buff.Bytes())
	if err != nil {
		panic(err)
	}
	return credentials
}

func (a *Authentification) InitIdentificationMessage() {
	a.initLoginAction()
	identification := messages.Types_[int(messages.IdentificationID)].GetNOA(a.instance).(*messages.Identification)

	identification.Lang = a.Lang
	identification.AutoSelectServer = a.lA.autoSelectServer

	currentVersion := structs.GetVersionNOA()
	identification.Version.Major = currentVersion.Major
	identification.Version.Minor = currentVersion.Minor
	identification.Version.Code = currentVersion.Code
	identification.Version.Build = currentVersion.Build
	identification.Version.BuildType = currentVersion.BuildType

	identification.Credentials = a.getCipher()
}

func (a *Authentification) getPublicKey() *rsa.PublicKey {
	hc := messages.Types_[messages.HelloConnectID].GetNOA(a.instance).(*messages.HelloConnect)

	if a.publicKey != nil && bytes.Compare(hc.Salt, a.salt) == 0 {
		return a.publicKey
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

	a.publicKey = pkiX.(*rsa.PublicKey)
	return a.publicKey
}

func (a *Authentification) getSalt() []byte {
	hc := messages.Types_[messages.HelloConnectID].GetNOA(a.instance).(*messages.HelloConnect)
	if hc.Salt == nil {
		panic("helloMessage wasn't call")
	}

	mySalt := hc.Salt
	a.salt = hc.Salt

	if len(mySalt) < 32 {
		for len(mySalt) < 32 {
			mySalt = append(mySalt, byte(' '))
		}
	}

	return mySalt
}

func generateAESKey() []byte {
	aes := make([]byte, AESLength)

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := uint(0); i < AESLength; i++ {
		aes[i] = byte(random.Intn(255))
	}

	return aes
}
