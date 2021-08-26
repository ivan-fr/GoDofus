package managers

import (
	"GoDofus/messages"
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
	"gopkg.in/yaml.v2"
	"math/rand"
	"os"
	"time"
)

type myLogin struct {
	Ndc  string `yaml:"nomdecompte"`
	Pass string `yaml:"motdepasse"`
}

func getConf() *myLogin {
	var login = &myLogin{}

	yamlFile, err := os.ReadFile("./login.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, login)
	if err != nil {
		panic(err)
	}

	return login
}

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
	lang      []byte
	publicKey *rsa.PublicKey
	salt      []byte
}

var authenticate_ = &Authentification{AESKey: generateAESKey(), lang: []byte("fr")}
var AESLength = uint(32)

var myLogin_ = getConf()
var publicVerifyPem = utils.ReadRSA("./binaryData/verify_key.bin")
var blockVerify = utils.DecodePem(publicVerifyPem)
var publicKeyVerify = utils.PublicKeyOf(blockVerify)

func GetAuthentification() *Authentification {
	return authenticate_
}

func (a *Authentification) initLoginAction() {
	la := &loginAction{autoSelectServer: true}
	la.username = myLogin_.Ndc
	la.password = myLogin_.Pass
	a.lA = la
	var randomTime = time.Duration(rand.Intn(2) + 2)
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
	identification := messages.GetIdentificationNOA()

	identification.AesKEY_ = make([]byte, len(a.AESKey))
	copy(identification.AesKEY_, a.AESKey)

	identification.Lang = a.lang
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
	hc := messages.GetHelloConnectNOA()

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
	hc := messages.GetHelloConnectNOA()
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
