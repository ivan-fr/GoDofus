package managers

import (
	"GoDofus/messages"
	"GoDofus/structs"
	"bytes"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"os"
	"time"
)

type myLogin struct {
	ndc  string `yaml:"nomdecompte"`
	pass string `yaml:"motdepasse"`
}

func getConf() *myLogin {
	var login = &myLogin{}

	yamlFile, err := ioutil.ReadFile("./login.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, login)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
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

type authentification struct {
	AESKey    []byte
	lA        *loginAction
	lang      []byte
	publicKey *rsa.PublicKey
	salt      []byte
}

var authenticate_ = &authentification{AESKey: generateAESKey(), lang: []byte("fr")}
var AESLength = uint(32)

var myLogin_ = getConf()
var publicVerifyPem = readVerify()
var blockVerify = decodeVerifyPem()
var publicKeyVerify = theVerifyPublicKey()

func rsaPublicDecrypt(pubKey *rsa.PublicKey, data []byte) []byte {
	c := new(big.Int)
	m := new(big.Int)
	m.SetBytes(data)
	e := big.NewInt(int64(pubKey.E))
	c.Exp(m, e, pubKey.N)
	out := c.Bytes()
	skip := 0
	for i := 2; i < len(out); i++ {
		if i+1 >= len(out) {
			break
		}
		if out[i] == 0xff && out[i+1] == 0 {
			skip = i + 2
			break
		}
	}
	return out[skip:]
}

func readVerify() []byte {
	publicVerifyPem, err := os.ReadFile("./binaryData/verify_key.bin")
	if err != nil {
		panic(err)
	}
	return publicVerifyPem
}

func decodeVerifyPem() *pem.Block {
	var blockVerify, _ = pem.Decode(publicVerifyPem)
	if blockVerify == nil {
		panic("block empty")
	}
	return blockVerify
}

func theVerifyPublicKey() *rsa.PublicKey {
	publicKeyVerify, err := x509.ParsePKIXPublicKey(blockVerify.Bytes)
	if err != nil {
		panic(err)
	}
	p := publicKeyVerify.(*rsa.PublicKey)
	return p
}

func GetAuthentification() *authentification {
	return authenticate_
}

func (a *authentification) initLoginAction() {
	la := &loginAction{autoSelectServer: true}
	la.username = myLogin_.ndc
	la.password = myLogin_.pass
	a.lA = la

	time.Sleep(time.Second * 2)
}

func (a *authentification) getCipher() []byte {
	buff := new(bytes.Buffer)
	mySalt := a.getSalt()
	_ = binary.Write(buff, binary.BigEndian, []byte(mySalt))
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

func (a *authentification) InitIdentificationMessage() {
	a.initLoginAction()
	identification := messages.GetIdentificationNOA()
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

func (a *authentification) getPublicKey() *rsa.PublicKey {
	hc := messages.GetHelloConnectNOA()

	if a.publicKey != nil && bytes.Compare(hc.Salt, a.salt) == 0 {
		return a.publicKey
	}

	if hc.Key == nil {
		panic("helloMessage wasn't call")
	}

	publicKey := rsaPublicDecrypt(publicKeyVerify, hc.Key)

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

func (a *authentification) getSalt() []byte {
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
